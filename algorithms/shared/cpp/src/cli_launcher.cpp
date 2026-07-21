#include "compresskit/cli_launcher.hpp"

#include <iostream>

namespace compresskit {
namespace cli {

namespace {
// program + mode + input + output
constexpr int EXPECTED_ARGC = 4;
}  // namespace

int run(const Algorithm& algo, int argc, char** argv) {
    if (argc != EXPECTED_ARGC) {
        std::cerr << "Usage: " << argv[0] << " encode|decode input output\n";
        return 1;
    }

    std::string mode = argv[1];
    std::string input_path = argv[2];
    std::string output_path = argv[3];

    bool ok = true;

    if (mode == "encode") {
        ok = encode_file_via_buffer(algo.encode, input_path, output_path);
    } else if (mode == "decode") {
        ok = decode_file_via_buffer(algo.decode, input_path, output_path);
    } else {
        std::cerr << "unknown mode, expected encode or decode\n";
        return 1;
    }

    return ok ? 0 : 1;
}

}  // namespace cli
}  // namespace compresskit
