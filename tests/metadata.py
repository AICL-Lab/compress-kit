"""CompressKit test and benchmark metadata module.

Single source of truth for algorithm registry, corpus registry, Range corpus cap,
benchmark job registry, and report schema.  Conformance, smoke, test-data, and
benchmark scripts are thin adapters over this module.

Design decision: docs/adr/0001-validation-metadata-module-shape.md
Range corpus cap decision: docs/adr/0003-range-coder-corpus-cap-policy.md

Consumers must assert the version they depend on:

    import metadata
    assert metadata.METADATA_VERSION == "1.0"
"""

from __future__ import annotations

from dataclasses import dataclass, field
from pathlib import Path

# ---------------------------------------------------------------------------
# Version
# ---------------------------------------------------------------------------

METADATA_VERSION = "1.0"

# ---------------------------------------------------------------------------
# Paths
# ---------------------------------------------------------------------------

ROOT = Path(__file__).resolve().parent.parent
TESTS_DIR = ROOT / "tests"
DATA_DIR = TESTS_DIR / "data"
SCRIPTS_DIR = ROOT / "scripts"
REPORTS_DIR = ROOT / "reports"
DOCS_DATA_DIR = ROOT / "docs" / ".vitepress" / "data"

# ---------------------------------------------------------------------------
# Algorithm and binary registry
# ---------------------------------------------------------------------------


@dataclass(frozen=True)
class BinaryEntry:
    language: str
    path: Path


@dataclass(frozen=True)
class AlgorithmEntry:
    """One algorithm supported by CompressKit."""

    name: str
    """Canonical lower-case algorithm name (e.g. 'huffman')."""

    extension: str
    """File extension used for encoded output (e.g. 'huf')."""

    binaries: tuple[BinaryEntry, ...]
    """Shipped binaries, in display order (cpp)."""

    benchmark_input: str
    """Corpus file name used as the benchmark input for this algorithm."""

    benchmark_driver: Path
    """Path to the per-algorithm benchmark driver script."""


ALGORITHMS: tuple[AlgorithmEntry, ...] = (
    AlgorithmEntry(
        name="huffman",
        extension="huf",
        binaries=(
            BinaryEntry("cpp", ROOT / "algorithms/huffman/cpp/huffman_cpp"),
            BinaryEntry("go", ROOT / "algorithms/huffman/go/huffman_go"),
            BinaryEntry("rust", ROOT / "algorithms/huffman/rust/huffman_rust"),
        ),
        benchmark_input="textlike_10MiB.bin",
        benchmark_driver=ROOT / "algorithms/huffman/benchmark/bench.py",
    ),
    AlgorithmEntry(
        name="arithmetic",
        extension="aenc",
        binaries=(
            BinaryEntry("cpp", ROOT / "algorithms/arithmetic/cpp/arithmetic_cpp"),
            BinaryEntry("go", ROOT / "algorithms/arithmetic/go/arithmetic_go"),
            BinaryEntry("rust", ROOT / "algorithms/arithmetic/rust/arithmetic_rust"),
        ),
        benchmark_input="textlike_10MiB.bin",
        benchmark_driver=ROOT / "algorithms/arithmetic/benchmark/bench.py",
    ),
    AlgorithmEntry(
        name="range",
        extension="rcnc",
        binaries=(
            BinaryEntry("cpp", ROOT / "algorithms/range/cpp/rangecoder_cpp"),
            BinaryEntry("go", ROOT / "algorithms/range/go/rangecoder_go"),
            BinaryEntry("rust", ROOT / "algorithms/range/rust/target/release/rangecoder"),
        ),
        benchmark_input="small_dictionary_like.bin",
        benchmark_driver=ROOT / "algorithms/range/benchmark/bench.py",
    ),
    AlgorithmEntry(
        name="rle",
        extension="rle",
        binaries=(
            BinaryEntry("cpp", ROOT / "algorithms/rle/cpp/rle_cpp"),
            BinaryEntry("go", ROOT / "algorithms/rle/go/rle_go"),
            BinaryEntry("rust", ROOT / "algorithms/rle/rust/rle_rust"),
        ),
        benchmark_input="repetitive_10MiB.bin",
        benchmark_driver=ROOT / "algorithms/rle/benchmark/bench.py",
    ),
)

ALGORITHM_ORDER: tuple[str, ...] = tuple(a.name for a in ALGORITHMS)
LANGUAGE_ORDER: tuple[str, ...] = ("cpp",)

# ---------------------------------------------------------------------------
# Corpus registry
# ---------------------------------------------------------------------------


@dataclass(frozen=True)
class CorpusEntry:
    """A test corpus file under tests/data/."""

    name: str
    """File name (e.g. 'empty.bin')."""

    is_large: bool = False
    """If True, only included when --include-large is passed to conformance."""


CORPUS: tuple[CorpusEntry, ...] = (
    CorpusEntry("empty.bin"),
    CorpusEntry("single_byte.bin"),
    CorpusEntry("alternating.bin"),
    CorpusEntry("small_dictionary_like.bin"),
    CorpusEntry("random_1MiB.bin", is_large=True),
    CorpusEntry("random_10MiB.bin", is_large=True),
    CorpusEntry("repetitive_10MiB.bin", is_large=True),
    CorpusEntry("textlike_10MiB.bin", is_large=True),
)

DEFAULT_CORPUS: tuple[str, ...] = tuple(c.name for c in CORPUS if not c.is_large)
LARGE_CORPUS: tuple[str, ...] = tuple(c.name for c in CORPUS if c.is_large)

# ---------------------------------------------------------------------------
# Range corpus cap
#
# See docs/adr/0003-range-coder-corpus-cap-policy.md for the policy decision.
# Do NOT change this value without an OpenSpec change.
# ---------------------------------------------------------------------------

RANGE_CORPUS_CAP_BYTES: int = 100 * 1024  # 100 KiB
RANGE_CORPUS_CAP_REASON: str = "range_coder_corpus_cap_100_kib"

# ---------------------------------------------------------------------------
# Benchmark report schema
#
# Field names used in benchmark JSON output.  Changing these breaks the docs
# page that consumes benchmarks.json.
# ---------------------------------------------------------------------------

BENCHMARK_REPORT_FIELDS: tuple[str, ...] = (
    "algorithm",
    "language",
    "dataset",
    "encodeTime",
    "decodeTime",
    "encodeSpeed",
    "decodeSpeed",
    "compressionRatio",
    "throughput",
)

# ---------------------------------------------------------------------------
# Convenience accessors
# ---------------------------------------------------------------------------


def algorithm_by_name(name: str) -> AlgorithmEntry:
    """Return the AlgorithmEntry for *name*, raising KeyError if not found."""
    for algo in ALGORITHMS:
        if algo.name == name:
            return algo
    raise KeyError(f"unknown algorithm: {name!r}")


def all_binaries() -> list[tuple[str, str, Path]]:
    """Return [(algorithm_name, language, binary_path)] for all shipped binaries."""
    return [
        (algo.name, binary.language, binary.path)
        for algo in ALGORITHMS
        for binary in algo.binaries
    ]


def corpus_files(include_large: bool = False) -> tuple[str, ...]:
    """Return corpus file names eligible for conformance testing."""
    if include_large:
        return tuple(c.name for c in CORPUS)
    return DEFAULT_CORPUS
