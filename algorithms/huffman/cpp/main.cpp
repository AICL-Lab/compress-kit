#include <cstdint>
#include <fstream>
#include <iostream>
#include <memory>
#include <queue>
#include <string>
#include <vector>

#include "compresskit/buffer_api.hpp"

class BitWriter {
   public:
    explicit BitWriter(std::ostream& s) : stream(s), buffer(0), bits_in_buffer(0) {}

    void write_bit(int bit) {
        buffer = static_cast<uint8_t>((buffer << 1) | (bit & 1));
        bits_in_buffer++;
        if (bits_in_buffer == 8) {
            stream.put(static_cast<char>(buffer));
            bits_in_buffer = 0;
            buffer = 0;
        }
    }

    void flush() {
        if (bits_in_buffer > 0) {
            buffer <<= (8 - bits_in_buffer);
            stream.put(static_cast<char>(buffer));
            bits_in_buffer = 0;
            buffer = 0;
        }
    }

   private:
    std::ostream& stream;
    uint8_t buffer;
    int bits_in_buffer;
};

class BitReader {
   public:
    explicit BitReader(std::istream& s)
        : stream(s), current_byte(0), bits_remaining(0), reached_eof(false) {}

    int read_bit() {
        if (bits_remaining == 0) {
            int c = stream.get();
            if (c == EOF) {
                reached_eof = true;
                return 0;
            }
            current_byte = static_cast<uint8_t>(c);
            bits_remaining = 8;
        }
        bits_remaining--;
        return (current_byte >> bits_remaining) & 1;
    }

    bool eof() const { return reached_eof; }

   private:
    std::istream& stream;
    uint8_t current_byte;
    int bits_remaining;
    bool reached_eof;
};

static const uint32_t SYMBOL_LIMIT = 257;
static const uint32_t EOF_SYMBOL = SYMBOL_LIMIT - 1;
static const uint64_t MAX_INPUT_SIZE = 4ULL * 1024 * 1024 * 1024;  // 4 GiB max

struct Node {
    uint32_t symbol;
    uint64_t freq;
    Node* left;
    Node* right;
};

// Custom deleter for Node pointers to enable RAII with smart pointers
struct NodeDeleter {
    void operator()(Node* node) const {
        if (!node)
            return;
        // Recursively delete children
        NodeDeleter deleter;
        if (node->left)
            deleter(node->left);
        if (node->right)
            deleter(node->right);
        delete node;
    }
};

// RAII wrapper for Node pointers
using UniqueNode = std::unique_ptr<Node, NodeDeleter>;

static bool is_leaf(const Node* node) {
    return node->left == nullptr && node->right == nullptr;
}

struct NodeCmp {
    bool operator()(const Node* a, const Node* b) const {
        if (a->freq != b->freq) {
            return a->freq > b->freq;
        }
        return a->symbol > b->symbol;
    }
};

static Node* build_tree(const std::vector<uint32_t>& freq) {
    std::priority_queue<Node*, std::vector<Node*>, NodeCmp> pq;
    for (uint32_t s = 0; s < SYMBOL_LIMIT; s++) {
        if (freq[s] == 0) {
            continue;
        }
        Node* n = new Node;
        n->symbol = s;
        n->freq = freq[s];
        n->left = nullptr;
        n->right = nullptr;
        pq.push(n);
    }
    if (pq.empty()) {
        Node* n = new Node;
        n->symbol = EOF_SYMBOL;
        n->freq = 1;
        n->left = nullptr;
        n->right = nullptr;
        return n;
    }
    if (pq.size() == 1) {
        Node* only = pq.top();
        pq.pop();
        Node* parent = new Node;
        parent->symbol = only->symbol;
        parent->freq = only->freq;
        parent->left = only;
        parent->right = nullptr;
        pq.push(parent);
    }
    while (pq.size() > 1) {
        Node* a = pq.top();
        pq.pop();
        Node* b = pq.top();
        pq.pop();
        Node* parent = new Node;
        parent->symbol = a->symbol < b->symbol ? a->symbol : b->symbol;
        parent->freq = a->freq + b->freq;
        parent->left = a;
        parent->right = b;
        pq.push(parent);
    }
    return pq.top();
}

static void build_codes(Node* node, std::vector<std::string>& codes, std::string& prefix) {
    if (!node) {
        return;
    }
    if (is_leaf(node)) {
        if (prefix.empty()) {
            codes[node->symbol] = "0";
        } else {
            codes[node->symbol] = prefix;
        }
        return;
    }
    prefix.push_back('0');
    build_codes(node->left, codes, prefix);
    prefix.pop_back();
    prefix.push_back('1');
    build_codes(node->right, codes, prefix);
    prefix.pop_back();
}

static std::vector<uint32_t> build_frequencies_from_file(const std::string& input_path) {
    std::vector<uint32_t> freq(SYMBOL_LIMIT, 0);
    std::ifstream in(input_path, std::ios::binary);
    if (!in) {
        return freq;
    }
    char c;
    while (in.get(c)) {
        unsigned char uc = static_cast<unsigned char>(c);
        freq[static_cast<uint32_t>(uc)]++;
    }
    freq[EOF_SYMBOL] = 1;
    return freq;
}

static void write_frequencies(std::ostream& out, const std::vector<uint32_t>& freq) {
    uint32_t count = static_cast<uint32_t>(freq.size());
    out.write(reinterpret_cast<const char*>(&count), sizeof(count));
    for (uint32_t v : freq) {
        out.write(reinterpret_cast<const char*>(&v), sizeof(v));
    }
}

static bool read_frequencies(std::istream& in, std::vector<uint32_t>& freq) {
    uint32_t count = 0;
    in.read(reinterpret_cast<char*>(&count), sizeof(count));
    if (!in) {
        std::cerr << "Failed to read frequency table\n";
        return false;
    }
    if (count != SYMBOL_LIMIT) {
        std::cerr << "Bad frequency table size: " << count << "\n";
        return false;
    }
    freq.assign(count, 0);
    in.read(reinterpret_cast<char*>(freq.data()), freq.size() * sizeof(uint32_t));
    if (!in) {
        std::cerr << "Failed to read frequency table\n";
        return false;
    }
    return true;
}

static bool compress_file(const std::string& input_path, const std::string& output_path) {
    // Check input file size to prevent frequency overflow
    {
        std::ifstream check(input_path, std::ios::binary | std::ios::ate);
        if (check) {
            auto size = check.tellg();
            if (size > 0 && static_cast<uint64_t>(size) > MAX_INPUT_SIZE) {
                std::cerr << "Input file too large (max " << MAX_INPUT_SIZE << " bytes)\n";
                return false;
            }
        }
    }

    std::vector<uint32_t> freq = build_frequencies_from_file(input_path);
    UniqueNode root(build_tree(freq));  // RAII: automatic cleanup
    std::vector<std::string> codes(SYMBOL_LIMIT);
    std::string prefix;
    build_codes(root.get(), codes, prefix);

    std::ifstream in(input_path, std::ios::binary);
    if (!in) {
        std::cerr << "Cannot open input file for reading\n";
        return false;
    }
    std::ofstream out(output_path, std::ios::binary);
    if (!out) {
        std::cerr << "Cannot open output file for writing\n";
        return false;
    }

    const char magic[4] = {'H', 'F', 'M', 'N'};
    out.write(magic, sizeof(magic));
    write_frequencies(out, freq);

    BitWriter bit_writer(out);
    char c;
    while (in.get(c)) {
        uint32_t sym = static_cast<uint8_t>(c);
        const std::string& code = codes[sym];
        for (char b : code) {
            bit_writer.write_bit(b == '1' ? 1 : 0);
        }
    }
    const std::string& eof_code = codes[EOF_SYMBOL];
    for (char b : eof_code) {
        bit_writer.write_bit(b == '1' ? 1 : 0);
    }
    bit_writer.flush();

    if (in.bad()) {
        std::cerr << "Failed to read input file\n";
        return false;
    }
    if (!out) {
        std::cerr << "Failed to write output file\n";
        return false;
    }

    return true;
}

static bool decompress_file(const std::string& input_path, const std::string& output_path) {
    std::ifstream in(input_path, std::ios::binary);
    if (!in) {
        std::cerr << "Cannot open input file for reading\n";
        return false;
    }
    char magic[4] = {};
    in.read(magic, sizeof(magic));
    if (!in || magic[0] != 'H' || magic[1] != 'F' || magic[2] != 'M' || magic[3] != 'N') {
        std::cerr << "Invalid input file format\n";
        return false;
    }

    std::vector<uint32_t> freq;
    if (!read_frequencies(in, freq)) {
        return false;
    }
    UniqueNode root(build_tree(freq));  // RAII: automatic cleanup
    if (!root) {
        return false;
    }

    std::ofstream out(output_path, std::ios::binary);
    if (!out) {
        std::cerr << "Cannot open output file for writing\n";
        return false;
    }

    BitReader bit_reader(in);
    Node* node = root.get();  // Raw pointer for traversal
    bool saw_eof = false;
    bool ok = true;
    while (true) {
        int bit = bit_reader.read_bit();
        if (bit == 0) {
            node = node->left;
        } else {
            node = node->right;
        }
        if (!node) {
            std::cerr << "Input data corrupted or truncated\n";
            ok = false;
            break;
        }
        if (is_leaf(node)) {
            if (node->symbol == EOF_SYMBOL) {
                saw_eof = true;
                break;
            }
            unsigned char b = static_cast<unsigned char>(node->symbol);
            out.put(static_cast<char>(b));
            if (!out) {
                std::cerr << "Failed to write output file\n";
                ok = false;
                break;
            }
            node = root.get();
        }
        if (bit_reader.eof() && node == root.get()) {
            break;
        }
    }

    if (!saw_eof) {
        std::cerr << "Input data corrupted or truncated\n";
        ok = false;
    }
    return ok;
}

bool huffman_encode_file(const std::string& input_path, const std::string& output_path) {
    return compress_file(input_path, output_path);
}

bool huffman_decode_file(const std::string& input_path, const std::string& output_path) {
    return decompress_file(input_path, output_path);
}

#ifndef COMPRESSKIT_NO_MAIN
int main(int argc, char** argv) {
    if (argc != 4) {
        std::cerr << "Usage: " << argv[0] << " encode|decode input output\n";
        return 1;
    }
    std::string mode = argv[1];
    std::string input_path = argv[2];
    std::string output_path = argv[3];

    bool ok = true;

    if (mode == "encode") {
        ok = compresskit::encode_file_via_buffer(huffman_encode_file, input_path, output_path);
    } else if (mode == "decode") {
        ok = compresskit::decode_file_via_buffer(huffman_decode_file, input_path, output_path);
    } else {
        std::cerr << "Unknown mode\n";
        return 1;
    }

    return ok ? 0 : 1;
}
#endif
