#include <cstdint>
#include <fstream>
#include <iostream>
#include <string>
#include <vector>

#include "compresskit/buffer_api.hpp"
#include "compresskit/frequency_table.hpp"

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

class ArithmeticEncoder {
   public:
    explicit ArithmeticEncoder(BitWriter& w)
        : writer(w), low(0), high(FULL_RANGE - 1), pending_bits(0) {}

    void encode_symbol(uint32_t symbol, const std::vector<uint32_t>& cumulative) {
        uint64_t range = high - low + 1;
        uint64_t total = cumulative.back();
        uint64_t sym_low = cumulative[symbol];
        uint64_t sym_high = cumulative[symbol + 1];

        high = low + (range * sym_high) / total - 1;
        low = low + (range * sym_low) / total;

        for (;;) {
            if (high < HALF_RANGE) {
                output_bit(0);
            } else if (low >= HALF_RANGE) {
                output_bit(1);
                low -= HALF_RANGE;
                high -= HALF_RANGE;
            } else if (low >= FIRST_QUARTER && high < THIRD_QUARTER) {
                pending_bits++;
                low -= FIRST_QUARTER;
                high -= FIRST_QUARTER;
            } else {
                break;
            }
            low <<= 1;
            high = (high << 1) | 1;
        }
    }

    void finish() {
        pending_bits++;
        if (low < FIRST_QUARTER) {
            output_bit(0);
        } else {
            output_bit(1);
        }
        writer.flush();
    }

   private:
    static constexpr uint64_t STATE_BITS = 32;
    static constexpr uint64_t FULL_RANGE = (static_cast<uint64_t>(1) << STATE_BITS);
    static constexpr uint64_t HALF_RANGE = FULL_RANGE >> 1;
    static constexpr uint64_t FIRST_QUARTER = HALF_RANGE >> 1;
    static constexpr uint64_t THIRD_QUARTER = FIRST_QUARTER * 3;

    BitWriter& writer;
    uint64_t low;
    uint64_t high;
    uint64_t pending_bits;

    void output_bit(int bit) {
        writer.write_bit(bit);
        int complement = bit ^ 1;
        while (pending_bits > 0) {
            writer.write_bit(complement);
            pending_bits--;
        }
    }
};

class ArithmeticDecoder {
   public:
    explicit ArithmeticDecoder(BitReader& r) : reader(r), low(0), high(FULL_RANGE - 1), code(0) {
        for (uint64_t i = 0; i < STATE_BITS; i++) {
            code = (code << 1) | static_cast<uint64_t>(reader.read_bit());
        }
    }

    uint32_t decode_symbol(const std::vector<uint32_t>& cumulative) {
        uint64_t range = high - low + 1;
        uint64_t total = cumulative.back();
        uint64_t offset = code - low;
        uint64_t value = ((offset + 1) * total - 1) / range;

        uint32_t lo = 0;
        uint32_t hi = static_cast<uint32_t>(cumulative.size() - 1);
        while (lo + 1 < hi) {
            uint32_t mid = lo + (hi - lo) / 2;
            if (cumulative[mid] > value) {
                hi = mid;
            } else {
                lo = mid;
            }
        }
        uint32_t symbol = lo;

        uint64_t sym_low = cumulative[symbol];
        uint64_t sym_high = cumulative[symbol + 1];

        high = low + (range * sym_high) / total - 1;
        low = low + (range * sym_low) / total;

        for (;;) {
            if (high < HALF_RANGE) {
            } else if (low >= HALF_RANGE) {
                low -= HALF_RANGE;
                high -= HALF_RANGE;
                code -= HALF_RANGE;
            } else if (low >= FIRST_QUARTER && high < THIRD_QUARTER) {
                low -= FIRST_QUARTER;
                high -= FIRST_QUARTER;
                code -= FIRST_QUARTER;
            } else {
                break;
            }
            low <<= 1;
            high = (high << 1) | 1;
            code = (code << 1) | static_cast<uint64_t>(reader.read_bit());
        }

        return symbol;
    }

   private:
    static constexpr uint64_t STATE_BITS = 32;
    static constexpr uint64_t FULL_RANGE = (static_cast<uint64_t>(1) << STATE_BITS);
    static constexpr uint64_t HALF_RANGE = FULL_RANGE >> 1;
    static constexpr uint64_t FIRST_QUARTER = HALF_RANGE >> 1;
    static constexpr uint64_t THIRD_QUARTER = FIRST_QUARTER * 3;

    BitReader& reader;
    uint64_t low;
    uint64_t high;
    uint64_t code;
};

static const uint32_t SYMBOL_LIMIT = 257;
static const uint32_t EOF_SYMBOL = SYMBOL_LIMIT - 1;
static const uint32_t MAX_TOTAL = 1u << 24;
static const uint64_t MAX_INPUT_SIZE = 4ULL * 1024 * 1024 * 1024;  // 4 GiB max

static void scale_frequencies(std::vector<uint32_t>& freq) {
    uint64_t total = 0;
    for (uint32_t f : freq) {
        total += f;
    }
    if (total == 0) {
        for (size_t i = 0; i < freq.size(); i++) {
            freq[i] = 1;
        }
        return;
    }
    if (total <= MAX_TOTAL) {
        return;
    }
    uint64_t new_total = 0;
    for (size_t i = 0; i < freq.size(); i++) {
        if (freq[i] == 0) {
            continue;
        }
        uint64_t scaled = static_cast<uint64_t>(freq[i]) * MAX_TOTAL / total;
        if (scaled == 0) {
            scaled = 1;
        }
        freq[i] = static_cast<uint32_t>(scaled);
        new_total += scaled;
    }
    if (new_total == 0) {
        uint32_t base = MAX_TOTAL / static_cast<uint32_t>(freq.size());
        if (base == 0) {
            base = 1;
        }
        for (size_t i = 0; i < freq.size(); i++) {
            freq[i] = base;
        }
    }
}

static std::vector<uint32_t> build_frequencies_from_file(const std::string& input_path) {
    std::vector<uint32_t> freq(SYMBOL_LIMIT, 0);
    std::ifstream in(input_path, std::ios::binary);
    if (!in) {
        return freq;
    }
    uint32_t overflow_symbol = 0;
    const auto status = compresskit::accumulate_frequencies(in, freq, &overflow_symbol);
    if (status == compresskit::FrequencyCountStatus::IO_ERROR) {
        std::cerr << "Failed to read input file\n";
        freq.clear();
        return freq;
    }
    if (status == compresskit::FrequencyCountStatus::OVERFLOW) {
        std::cerr << "Frequency overflow for symbol " << overflow_symbol << "\n";
        freq.clear();
        return freq;
    }
    freq[EOF_SYMBOL] = 1;
    scale_frequencies(freq);
    return freq;
}

static std::vector<uint32_t> build_cumulative(const std::vector<uint32_t>& freq) {
    std::vector<uint32_t> cumulative(freq.size() + 1, 0);
    for (size_t i = 0; i < freq.size(); i++) {
        cumulative[i + 1] = cumulative[i] + freq[i];
    }
    if (cumulative.back() == 0) {
        for (size_t i = 0; i < freq.size(); i++) {
            cumulative[i + 1] = static_cast<uint32_t>(i + 1);
        }
    }
    return cumulative;
}

static void write_frequencies(std::ostream& out, const std::vector<uint32_t>& freq) {
    compresskit::write_frequency_table(out, freq);
}

static bool read_frequencies(std::istream& in, std::vector<uint32_t>& freq) {
    uint32_t count = 0;
    const auto status = compresskit::read_frequency_table(in, freq, SYMBOL_LIMIT, &count);
    if (status == compresskit::FrequencyTableReadStatus::TRUNCATED) {
        std::cerr << "Failed to read frequency table\n";
        return false;
    }
    if (status == compresskit::FrequencyTableReadStatus::BAD_COUNT) {
        std::cerr << "Bad frequency table size: " << count << "\n";
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
    if (freq.empty()) {
        return false;
    }
    std::vector<uint32_t> cumulative = build_cumulative(freq);

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

    const char magic[4] = {'A', 'E', 'N', 'C'};
    out.write(magic, sizeof(magic));
    write_frequencies(out, freq);

    BitWriter bit_writer(out);
    ArithmeticEncoder encoder(bit_writer);

    char c;
    while (in.get(c)) {
        uint32_t sym = static_cast<uint8_t>(c);
        encoder.encode_symbol(sym, cumulative);
    }
    encoder.encode_symbol(EOF_SYMBOL, cumulative);
    encoder.finish();

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
    if (!in || magic[0] != 'A' || magic[1] != 'E' || magic[2] != 'N' || magic[3] != 'C') {
        std::cerr << "Invalid input file format\n";
        return false;
    }

    std::vector<uint32_t> freq;
    if (!read_frequencies(in, freq)) {
        return false;
    }
    std::vector<uint32_t> cumulative = build_cumulative(freq);

    std::ofstream out(output_path, std::ios::binary);
    if (!out) {
        std::cerr << "Cannot open output file for writing\n";
        return false;
    }

    BitReader bit_reader(in);
    ArithmeticDecoder decoder(bit_reader);

    for (;;) {
        uint32_t sym = decoder.decode_symbol(cumulative);
        if (sym == EOF_SYMBOL) {
            break;
        }
        unsigned char b = static_cast<unsigned char>(sym);
        out.put(static_cast<char>(b));
        if (!out) {
            std::cerr << "Failed to write output file\n";
            return false;
        }
    }

    return true;
}

bool arithmetic_encode_file(const std::string& input_path, const std::string& output_path) {
    return compress_file(input_path, output_path);
}

bool arithmetic_decode_file(const std::string& input_path, const std::string& output_path) {
    return decompress_file(input_path, output_path);
}

#ifndef COMPRESSKIT_NO_MAIN
#include "compresskit/cli_launcher.hpp"

int main(int argc, char** argv) {
    compresskit::cli::Algorithm algo{
        [](const std::string& in, const std::string& out) {
            return compresskit::encode_file_via_buffer(arithmetic_encode_file, in, out);
        },
        [](const std::string& in, const std::string& out) {
            return compresskit::decode_file_via_buffer(arithmetic_decode_file, in, out);
        }};
    return compresskit::cli::run("arithmetic", algo, argc, argv);
}
#endif
