# CompressKit contract inventory

Inventory date: 2026-05-31. This document records current observable contracts before architecture deepening. It is not a proposal and does not authorize behavior, binary-format, CLI, or public API changes.

## Source-of-truth map

| Area | Current source(s) | Inventory note |
|---|---|---|
| Product identity and required algorithm/language matrix | `openspec/specs/encoding-project/spec.md:16-27`, `AGENTS.md:3-21` | CompressKit is four algorithms (Huffman, Arithmetic, Range, RLE) in C++17, Go, and Rust with byte-compatible formats. |
| CLI, error, streaming, buffer, frequency, and security requirements | `openspec/specs/core-architecture/spec.md:26-77`, `openspec/specs/core-architecture/spec.md:141-239` | Requirements define unified CLI, 5-state lifecycle, transactional buffer errors, and 4 GiB / 1 GiB limits. |
| Cross-language and benchmark requirements | `openspec/specs/cross-language-testing/spec.md:12-67`, `openspec/specs/cross-language-testing/spec.md:78-106` | Cross-language decode matrix, generated test data, benchmark metrics, Range known issue, and streaming conformance tests are specified here. |
| Domain vocabulary | `CONTEXT.md:7-118`, `CONTEXT.md:147-209` | Domain terms include magic numbers, Frequency Table, streaming/buffer layers, errors, security limits, architecture layers, and docs hierarchy. |
| User docs | `docs/en/api/streaming.md:1-44`, `docs/en/testing/cross-language.md:1-43`, `docs/en/benchmarks/results.md:16-56` | User-facing docs restate selected contracts and the benchmark JSON schema. Treat OpenSpec + code as authoritative if docs drift. |
| Known doc drift to preserve for later cleanup | `docs/en/architecture/index.md:66-73` vs. `algorithms/shared/go/cli/launcher.go:18-37` | The architecture page shows a product-level `compress-kit encode --algo ...` example; shipped binaries currently use `&lt;binary&gt; &lt;encode|decode&gt; &lt;input&gt; &lt;output&gt;`. |

## Codec lifecycle contract

| Fact | Citations |
|---|---|
| The lifecycle states are READY, STREAMING, FLUSHING, FINISHED, and ERROR. Valid transitions and reset semantics are specified; errors enter ERROR except `BUF_TOO_SMALL`, which is transactional. | `openspec/specs/core-architecture/spec.md:145-180`, `docs/en/api/streaming.md:8-20` |
| Go exposes `codec.Encoder`/`codec.Decoder` with `Process`, `Flush`, `Finish`, `Reset`, and `State`; `BufferedEncoder`/`BufferedDecoder` collect input and run single-shot functions at `Finish`. | `algorithms/shared/go/codec/encoder.go:9-23`, `algorithms/shared/go/codec/encoder.go:33-124`, `algorithms/shared/go/codec/encoder.go:252-321` |
| Rust exposes equivalent `Encoder`/`Decoder` traits and `BufferedEncoder`/`BufferedDecoder`; state and transactional `BufTooSmall` behavior are documented and tested in the module. | `algorithms/shared/rust/src/codec/encoder.rs:3-79`, `algorithms/shared/rust/src/codec/buffered.rs:1-18`, `algorithms/shared/rust/src/codec/buffered.rs:55-118` |
| C++ exposes virtual `Encoder`/`Decoder`, `State`, `ByteView`, `MutableByteView`, and `BufferEncoder`/`BufferDecoder`; the buffer adapters wrap file-based transforms. | `algorithms/shared/cpp/include/compresskit/encoder.hpp:10-48`, `algorithms/shared/cpp/include/compresskit/buffer_api.hpp:11-27`, `algorithms/shared/cpp/src/buffer_api.cpp:149-285` |
| Huffman, Arithmetic, Range, and RLE Go adapters currently use buffered streaming because they need complete input or are implemented as one-shot wrappers. | `algorithms/huffman/go/huffman.go:272-296`, `algorithms/arithmetic/go/arithmetic.go:306-330`, `algorithms/range/go/rangecoder.go:207-219`, `algorithms/rle/go/rle.go:159-183` |
| Rust algorithm CLIs call shared buffer helpers over streaming adapters. | `algorithms/huffman/rust/main.rs:16-36`, `algorithms/arithmetic/rust/main.rs:16-36`, `algorithms/range/rust/src/bin/rangecoder.rs:16-36`, `algorithms/rle/rust/main.rs:16-36` |

## Error contract

| Fact | Citations |
|---|---|
| Standard semantic errors are `BUF_TOO_SMALL`, `ERR_TRUNCATED`, `ERR_CORRUPT`, `ERR_INVALID_STATE`, `ERR_SIZE_LIMIT`, `ERR_VERSION_UNSUPPORTED`, and `ERR_UNKNOWN_ALGO`; Go also has `KindIO` for wrapped I/O. | `CONTEXT.md:101-118`, `algorithms/shared/go/codec/errors.go:8-29`, `algorithms/shared/cpp/include/compresskit/result.hpp:7-16`, `algorithms/shared/rust/src/codec/error.rs:20-39` |
| Go uses `ErrorKind`, `CodecError`, sentinel errors, and `errors.Is` mapping. | `algorithms/shared/go/codec/errors.go:31-86`, `algorithms/shared/go/codec/errors.go:88-112` |
| Rust uses a `CodecError` enum, string display, and helper mapping from `io::Error` to codec errors. | `algorithms/shared/rust/src/codec/error.rs:20-60`, `algorithms/shared/rust/src/codec/error.rs:62-99` |
| C++ uses `StatusCode` inside `Result&lt;T&gt;`; CLI file transforms return boolean success/failure. | `algorithms/shared/cpp/include/compresskit/result.hpp:7-26`, `algorithms/shared/cpp/include/compresskit/cli_launcher.hpp:10-17` |
| CLI programs report error by nonzero exit; OpenSpec requires English actionable messages and exit code 1 for errors. | `openspec/specs/core-architecture/spec.md:60-77`, `algorithms/shared/cpp/src/cli_launcher.cpp:8-29`, `algorithms/shared/go/cli/launcher.go:18-42`, `algorithms/shared/rust/src/cli.rs:9-33` |
| Some algorithm-level direct APIs still map errors through language-local messages before the shared codec adapter normalizes them; do not assume structural API identity across languages. | `algorithms/range/rust/src/lib.rs:216-236`, `algorithms/rle/rust/src/lib.rs:118-136`, `algorithms/range/go/rangecoder.go:40-55` |

## Buffer and streaming contract

| Fact | Citations |
|---|---|
| Buffer APIs are equivalent to `new encoder -&gt; process(input) -&gt; finish()` and similarly for decode. | `openspec/specs/core-architecture/spec.md:204-219`, `docs/en/api/streaming.md:35-44` |
| Go and Rust resizing buffers retry on `BufTooSmall` and cap decoded output at 1 GiB; encode buffer growth is bounded by an input-derived encode limit. | `algorithms/shared/go/codec/buffer_policy.go:23-69`, `algorithms/shared/rust/src/codec/buffer.rs:5-29`, `algorithms/shared/rust/src/codec/buffer.rs:32-53` |
| C++ buffer adapters copy through temporary files to adapt existing file-based algorithms; this is an implementation detail but affects locality for future refactors. | `algorithms/shared/cpp/include/compresskit/buffer_api.hpp:11-26`, `algorithms/shared/cpp/src/buffer_api.cpp:97-118` |
| Flush is a no-op for buffered static-model algorithms until `Finish`; RLE is currently also buffered in Go/Rust adapters despite being naturally streamable. | `CONTEXT.md:55-59`, `algorithms/rle/go/rle.go:159-183`, `algorithms/rle/rust/src/lib.rs:95-115` |

## Security limits

| Limit | Current facts and enforcement points |
|---|---|
| Input limit: 4 GiB | Required by OpenSpec and CONTEXT (`openspec/specs/encoding-project/spec.md:39-51`, `CONTEXT.md:89-99`). Shared Go/Rust/C++ codec layers define/enforce it (`algorithms/shared/go/codec/errors.go:114-121`, `algorithms/shared/rust/src/codec/error.rs:58-60`, `algorithms/shared/cpp/src/buffer_api.cpp:17-18`, `algorithms/shared/cpp/src/buffer_api.cpp:287-325`). Algorithm file paths also check it in places such as Huffman/Arithmetic Go and C++ (`algorithms/huffman/go/huffman.go:22-25`, `algorithms/arithmetic/go/arithmetic.go:23-24`, `algorithms/huffman/cpp/main.cpp:218-227`, `algorithms/arithmetic/cpp/main.cpp:295-305`). |
| Decoded output limit: 1 GiB | Required by OpenSpec and CONTEXT (`openspec/specs/core-architecture/spec.md:230-239`, `CONTEXT.md:97-99`). Shared Go/Rust/C++ decode buffers enforce it (`algorithms/shared/go/codec/encoder.go:224-228`, `algorithms/shared/rust/src/codec/buffered.rs:212-225`, `algorithms/shared/cpp/src/buffer_api.cpp:322-347`). Algorithm decoders also enforce it in Go/Rust Huffman/Arithmetic/Range/RLE and C++ Range/RLE (`algorithms/huffman/go/huffman.go:252-255`, `algorithms/arithmetic/go/arithmetic.go:294-297`, `algorithms/range/go/rangecoder.go:199-202`, `algorithms/rle/go/rle.go:114-116`, `algorithms/huffman/rust/src/lib.rs:199-204`, `algorithms/arithmetic/rust/src/lib.rs:167-172`, `algorithms/range/rust/src/lib.rs:207-209`, `algorithms/rle/rust/src/lib.rs:79-85`, `algorithms/range/cpp/main.cpp:273-276`, `algorithms/rle/cpp/main.cpp:169-173`). |

## Binary format and Frequency Table facts

| Algorithm | Current encoded prefix and payload facts | Citations |
|---|---|---|
| Huffman | Magic `HFMN`; writes shared Frequency Table; payload is bit-packed Huffman code stream ending with EOF symbol 256. | `CONTEXT.md:7-28`, `algorithms/huffman/go/huffman.go:176-205`, `algorithms/huffman/go/huffman.go:212-218`, `algorithms/huffman/rust/src/lib.rs:139-156`, `algorithms/huffman/rust/src/lib.rs:166-179`, `algorithms/huffman/cpp/main.cpp:250-267`, `algorithms/huffman/cpp/main.cpp:287-297` |
| Arithmetic | Magic `AENC`; writes scaled Frequency Table (`MaxTotal = 1 << 24`); payload is bit-packed arithmetic-coded data ending with EOF symbol. | `algorithms/arithmetic/go/arithmetic.go:14-24`, `algorithms/arithmetic/go/arithmetic.go:239-263`, `algorithms/arithmetic/go/arithmetic.go:270-279`, `algorithms/arithmetic/rust/src/lib.rs:6-24`, `algorithms/arithmetic/rust/src/lib.rs:105-128`, `algorithms/arithmetic/cpp/main.cpp:323-336`, `algorithms/arithmetic/cpp/main.cpp:356-367` |
| Range Coder | Magic `RCNC`; writes scaled Frequency Table (`MAX_TOTAL = 1 << 24`); payload is byte-oriented range-coded data ending with EOF symbol. Range C++ keeps local `write_header`/`read_header` instead of shared Frequency Table helpers. | `algorithms/range/go/rangecoder.go:7-11`, `algorithms/range/go/rangecoder.go:35-55`, `algorithms/range/rust/src/lib.rs:10-28`, `algorithms/range/rust/src/lib.rs:30-49`, `algorithms/range/cpp/main.cpp:13-19`, `algorithms/range/cpp/main.cpp:80-130`, `algorithms/range/cpp/main.cpp:228-247` |
| RLE | Magic `RLE\x00`; no Frequency Table; payload is repeated `(count: uint32 little-endian, value: byte)` pairs; count must be nonzero. Empty input encodes to magic only. | `CONTEXT.md:141-145`, `algorithms/rle/go/rle.go:14-19`, `algorithms/rle/go/rle.go:21-44`, `algorithms/rle/go/rle.go:99-124`, `algorithms/rle/cpp/main.cpp:10-24`, `algorithms/rle/cpp/main.cpp:80-119`, `algorithms/rle/rust/src/lib.rs:6-37`, `algorithms/rle/rust/src/lib.rs:40-85` |
| Shared Frequency Table | Static-model algorithms use 257 symbols: byte values 0-255 plus EOF 256. Wire format is `uint32_le count` followed by `count` `uint32_le` frequencies; expected exact count is 257 for Huffman/Arithmetic and Range decoders after header validation. | `CONTEXT.md:20-35`, `openspec/specs/core-architecture/spec.md:45-59`, `algorithms/shared/go/codec/frequency.go:10-18`, `algorithms/shared/go/codec/frequency.go:271-326`, `algorithms/shared/rust/src/codec/frequency.rs:4-6`, `algorithms/shared/rust/src/codec/frequency.rs:190-242`, `algorithms/shared/cpp/src/frequency_table.cpp:33-68` |
| Frequency validation | Go/Rust general readers reject count 0 or count > 1024; exact readers reject any count other than the algorithm's expected symbol count. C++ shared reader rejects mismatched expected count but does not impose the 1024 cap in the shared exact path. | `algorithms/shared/go/codec/frequency.go:36-69`, `algorithms/shared/go/codec/frequency.go:306-326`, `algorithms/shared/rust/src/codec/frequency.rs:197-242`, `algorithms/shared/cpp/src/frequency_table.cpp:45-68` |

## CLI contract and exceptions

| Fact | Citations |
|---|---|
| Public CLI contract is `&lt;binary&gt; &lt;encode|decode&gt; &lt;input&gt; &lt;output&gt;` with arity exactly 4 including argv[0] in launchers. | `openspec/specs/encoding-project/spec.md:53-66`, `openspec/specs/core-architecture/spec.md:26-43`, `algorithms/shared/cpp/src/cli_launcher.cpp:8-29`, `algorithms/shared/go/cli/launcher.go:18-37`, `algorithms/shared/rust/src/cli.rs:9-28` |
| C++ and Go algorithm entrypoints are thin adapters to shared CLI launchers. Rust entrypoints use shared `cli::run`. | `algorithms/huffman/cpp/main.cpp:359-371`, `algorithms/arithmetic/cpp/main.cpp:402-414`, `algorithms/rle/cpp/main.cpp:203-215`, `algorithms/huffman/go/cmd/main.go:18-20`, `algorithms/arithmetic/go/cmd/main.go:18-20`, `algorithms/range/go/cmd/main.go:19-21`, `algorithms/rle/go/cmd/main.go:18-20`, `algorithms/shared/rust/src/cli.rs:9-33` |
| Range C++ is the known exception: the production binary advertises and accepts `bench [size_bytes] [iterations]` in the same entrypoint. This is current behavior, not a change request. | `algorithms/range/cpp/main.cpp:378-415` |
| CLI smoke tests encode the current contract by checking usage, wrong arity, invalid mode, and round-trips for all binary paths. | `tests/conformance/run_cli_smoke.py:12-45`, `tests/conformance/run_cli_smoke.py:80-118`, `tests/conformance/run_cli_smoke.py:121-167` |
| Makefile names the shipped binary paths consumed by conformance and smoke tests. | `Makefile:11-32`, `tests/conformance/run_decode_matrix.py:31-68`, `tests/conformance/run_cli_smoke.py:17-38` |

## Validation and test metadata consumers

| Consumer | Current metadata shape | Citations |
|---|---|---|
| Build/test gates | `make test` runs generated data, shared C++/Go/Rust tests, algorithm Go/Rust tests, decode matrix, and CLI smoke. `test-conformance` and `test-cli-smoke` both depend on `build test-data`. | `Makefile:34-80` |
| Test data generator | Generates deterministic or literal corpus names under `tests/data`: random 1/10 MiB, repetitive 10 MiB, textlike 10 MiB, empty, single byte, alternating, small dictionary-like. Random binary uses `os.urandom`, while other nonliteral distributions use seeded `random.Random`. | `tests/gen_testdata.py:6-18`, `tests/gen_testdata.py:38-95` |
| Decode matrix | Hardcodes algorithm name, extension, binary path, language labels, default corpus names, optional large corpus names, per-command timeout, temp dir under `tests`, and encoded/decoded filename schema. | `tests/conformance/run_decode_matrix.py:18-68`, `tests/conformance/run_decode_matrix.py:71-100`, `tests/conformance/run_decode_matrix.py:160-224` |
| Range conformance cap | Range Coder skips corpus files over 100 KiB in decode matrix with reason `range_coder_corpus_cap_100_kib`; README documents this as workaround for known large-file decode performance. | `tests/conformance/run_decode_matrix.py:95-99`, `tests/conformance/run_decode_matrix.py:160-162`, `tests/conformance/README.md:32-34`, `docs/en/algorithms/range.md:160-168` |
| CLI smoke | Separately duplicates algorithm-to-binary metadata and small corpus list, then tests usage, arity, invalid mode, and intra-binary round trips. | `tests/conformance/run_cli_smoke.py:12-45`, `tests/conformance/run_cli_smoke.py:121-167` |
| Conformance docs | README says default matrix is 4 algorithms x 4 corpus files x 3 encoders x 3 decoders = 144 checks. | `tests/conformance/README.md:17-34` |

## Benchmark metadata consumers

| Consumer | Current metadata shape | Citations |
|---|---|---|
| `make bench` | Runs `python3 scripts/run_all_bench.py` after `test-data`. | `Makefile:82-89` |
| Benchmark orchestrator | Hardcodes algorithm order, language order, row regex, original-size regex, dataset regex, job list, driver paths, and selected input file per algorithm. | `scripts/run_all_bench.py:10-25`, `scripts/run_all_bench.py:27-53`, `scripts/run_all_bench.py:103-147` |
| Benchmark job inputs | Huffman and Arithmetic use `textlike_10MiB.bin`, Range uses `small_dictionary_like.bin`, and RLE uses `repetitive_10MiB.bin`. | `scripts/run_all_bench.py:27-53`, `docs/en/benchmarks/results.md:26-43` |
| Benchmark report and docs JSON | Orchestrator writes timestamped `reports/{algorithm}_report_*.txt` files and rewrites `docs/.vitepress/data/benchmarks.json` with `generated`, `version`, and sorted `results`. | `scripts/run_all_bench.py:68-89`, `scripts/run_all_bench.py:163-200`, `docs/en/benchmarks/results.md:16-56` |
| Docs benchmark schema | Each result row contains `algorithm`, `language`, `dataset`, `encodeTime`, `decodeTime`, `encodeSpeed`, `decodeSpeed`, `compressionRatio`, and `throughput`; current JSON stores one row per algorithm/language pair. | `scripts/run_all_bench.py:128-139`, `docs/.vitepress/data/benchmarks.json:1-14`, `docs/en/benchmarks/results.md:16-24` |
| Per-algorithm benchmark drivers | Drivers build one algorithm target, run C++/Go/Rust encode and decode, compare decoded output, write `Dataset`, `Original size`, build times, and rows matching `lang encode decode total ratio`. | `algorithms/huffman/benchmark/bench.py:43-53`, `algorithms/huffman/benchmark/bench.py:82-121`, `algorithms/range/benchmark/bench.py:45-55`, `algorithms/range/benchmark/bench.py:84-123` |
| Range benchmark default | Range benchmark defaults to `small_dictionary_like.bin` and documents that it is below the 100 KiB small-input constraint. | `algorithms/range/benchmark/bench.py:13-27` |

## OpenSpec-triggering decisions

Open a new OpenSpec change before making any of these intentional changes:

- Binary format: magic bytes, Frequency Table layout/count/endian rules, EOF symbol semantics, RLE `(count,value)` layout, Range/Arithmetic precision scaling, or any header/payload bytes (`openspec/specs/encoding-project/spec.md:25-38`, `openspec/specs/core-architecture/spec.md:45-59`, `CONTEXT.md:7-28`).
- Public CLI semantics: arity, mode names, exit-code semantics, usage contract, or the Range C++ `bench` exception's removal/deprecation path (`openspec/specs/encoding-project/spec.md:53-66`, `openspec/specs/core-architecture/spec.md:26-43`, `algorithms/range/cpp/main.cpp:378-415`).
- Public streaming/buffer APIs: lifecycle states, allowed transitions, transactional `BUF_TOO_SMALL`, reset behavior, buffer-layer equivalence, and language-visible error kind semantics (`openspec/specs/core-architecture/spec.md:141-239`, `openspec/specs/encoding-project/spec.md:123-141`).
- Security limits: 4 GiB input limit, 1 GiB decoded-output limit, or whether enforcement occurs at CLI, buffer, streaming, or algorithm boundaries (`openspec/specs/encoding-project/spec.md:39-51`, `openspec/specs/core-architecture/spec.md:230-239`).
- Test-gate semantics: algorithms/languages in the conformance matrix, corpus eligibility, Range 100 KiB skip policy, CLI smoke assertions, or required validation gates (`openspec/specs/cross-language-testing/spec.md:12-39`, `openspec/specs/cross-language-testing/spec.md:78-106`, `tests/conformance/run_decode_matrix.py:160-224`, `tests/conformance/run_cli_smoke.py:121-167`).
- Adding/removing algorithms or supported languages (`openspec/specs/encoding-project/spec.md:16-24`, `openspec/specs/cross-language-testing/spec.md:21-25`).

Contract-preserving refactors that should still get ADRs or design notes: extracting shared metadata without changing consumers, moving header/Frequency helpers behind language-local modules while preserving byte output, documenting RLE buffered-vs-incremental streaming stance, and moving Range C++ benchmark behavior only after an approved migration decision.
