#ifndef COMPRESSKIT_CLI_LAUNCHER_HPP
#define COMPRESSKIT_CLI_LAUNCHER_HPP

#include <functional>
#include <string>

namespace compresskit {
namespace cli {

using FileTransform = std::function<bool(const std::string&, const std::string&)>;

struct Algorithm {
    FileTransform encode;
    FileTransform decode;
};

int run(const Algorithm& algo, int argc, char** argv);

}  // namespace cli
}  // namespace compresskit

#endif  // COMPRESSKIT_CLI_LAUNCHER_HPP
