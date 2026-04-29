#!/usr/bin/env python3
"""Run the CompressKit cross-language decode conformance matrix."""

from __future__ import annotations

import argparse
import filecmp
import subprocess
import sys
import tempfile
from dataclasses import dataclass
from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]


@dataclass(frozen=True)
class Binary:
    language: str
    path: Path


@dataclass(frozen=True)
class Algorithm:
    name: str
    extension: str
    binaries: tuple[Binary, ...]


ALGORITHMS = (
    Algorithm(
        "huffman",
        "huf",
        (
            Binary("cpp", ROOT / "algorithms/huffman/cpp/huffman_cpp"),
            Binary("go", ROOT / "algorithms/huffman/go/huffman_go"),
            Binary("rust", ROOT / "algorithms/huffman/rust/huffman_rust"),
        ),
    ),
    Algorithm(
        "arithmetic",
        "aenc",
        (
            Binary("cpp", ROOT / "algorithms/arithmetic/cpp/arithmetic_cpp"),
            Binary("go", ROOT / "algorithms/arithmetic/go/arithmetic_go"),
            Binary("rust", ROOT / "algorithms/arithmetic/rust/arithmetic_rust"),
        ),
    ),
    Algorithm(
        "range",
        "rcnc",
        (
            Binary("cpp", ROOT / "algorithms/range/cpp/rangecoder_cpp"),
            Binary("go", ROOT / "algorithms/range/go/rangecoder_go"),
            Binary("rust", ROOT / "algorithms/range/rust/target/release/rangecoder"),
        ),
    ),
    Algorithm(
        "rle",
        "rle",
        (
            Binary("cpp", ROOT / "algorithms/rle/cpp/rle_cpp"),
            Binary("go", ROOT / "algorithms/rle/go/rle_go"),
            Binary("rust", ROOT / "algorithms/rle/rust/rle_rust"),
        ),
    ),
)


DEFAULT_CORPUS = (
    "empty.bin",
    "single_byte.bin",
    "alternating.bin",
    "small_dictionary_like.bin",
)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Run encode-language × decode-language matrix checks."
    )
    parser.add_argument(
        "--corpus-dir",
        type=Path,
        default=ROOT / "tests/data",
        help="Directory containing generated test corpus files.",
    )
    parser.add_argument(
        "--timeout",
        type=float,
        default=10.0,
        help="Per encode/decode command timeout in seconds.",
    )
    parser.add_argument(
        "--include-large",
        action="store_true",
        help="Also test large generated corpus files. Range still skips files over 100 KiB.",
    )
    return parser.parse_args()


def corpus_files(corpus_dir: Path, include_large: bool) -> list[Path]:
    names = list(DEFAULT_CORPUS)
    if include_large:
        names.extend(
            [
                "random_1MiB.bin",
                "repetitive_10MiB.bin",
                "textlike_10MiB.bin",
            ]
        )

    files = [corpus_dir / name for name in names]
    missing = [path for path in files if not path.is_file()]
    if missing:
        missing_text = "\n".join(str(path.relative_to(ROOT)) for path in missing)
        raise SystemExit(f"missing corpus file(s); run `make test-data` first:\n{missing_text}")
    return files


def ensure_binaries_exist() -> None:
    missing: list[Path] = []
    for algorithm in ALGORITHMS:
        for binary in algorithm.binaries:
            if not binary.path.is_file():
                missing.append(binary.path)

    if missing:
        missing_text = "\n".join(str(path.relative_to(ROOT)) for path in missing)
        raise SystemExit(f"missing binary file(s); run `make build` first:\n{missing_text}")


def run_command(command: list[str], timeout: float) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        command,
        cwd=ROOT,
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        timeout=timeout,
        check=False,
    )


def check_command(
    command: list[str], timeout: float, operation: str
) -> subprocess.CompletedProcess[str]:
    result = run_command(command, timeout)
    if result.returncode != 0:
        rendered = " ".join(command)
        raise RuntimeError(
            f"{operation} failed with exit code {result.returncode}: {rendered}\n"
            f"stdout:\n{result.stdout}\n"
            f"stderr:\n{result.stderr}"
        )
    return result


def should_skip(algorithm: Algorithm, corpus_file: Path) -> str | None:
    if algorithm.name == "range" and corpus_file.stat().st_size > 100 * 1024:
        return "range_coder_corpus_cap_100_kib"
    return None


def run_matrix(corpus: list[Path], timeout: float) -> tuple[int, int]:
    cases = 0
    skips = 0
    with tempfile.TemporaryDirectory(prefix=".conformance-", dir=ROOT / "tests") as tmp:
        tmpdir = Path(tmp)
        for algorithm in ALGORITHMS:
            for corpus_file in corpus:
                skip_reason = should_skip(algorithm, corpus_file)
                if skip_reason:
                    skips += len(algorithm.binaries) * len(algorithm.binaries)
                    print(
                        f"SKIP {algorithm.name} {corpus_file.name}: {skip_reason}",
                        flush=True,
                    )
                    continue

                for encoder in algorithm.binaries:
                    encoded = (
                        tmpdir
                        / f"{algorithm.name}-{encoder.language}-{corpus_file.name}.{algorithm.extension}"
                    )
                    check_command(
                        [
                            str(encoder.path),
                            "encode",
                            str(corpus_file),
                            str(encoded),
                        ],
                        timeout,
                        f"{algorithm.name} encode {encoder.language}",
                    )

                    for decoder in algorithm.binaries:
                        decoded = (
                            tmpdir
                            / f"{algorithm.name}-{encoder.language}-to-{decoder.language}-{corpus_file.name}"
                        )
                        check_command(
                            [
                                str(decoder.path),
                                "decode",
                                str(encoded),
                                str(decoded),
                            ],
                            timeout,
                            f"{algorithm.name} decode {encoder.language}->{decoder.language}",
                        )
                        if not filecmp.cmp(corpus_file, decoded, shallow=False):
                            raise RuntimeError(
                                f"{algorithm.name} {corpus_file.name} "
                                f"{encoder.language}->{decoder.language} produced mismatched output"
                            )
                        cases += 1
                        print(
                            f"PASS {algorithm.name} {corpus_file.name} "
                            f"{encoder.language}->{decoder.language}",
                            flush=True,
                        )
    return cases, skips


def main() -> int:
    args = parse_args()
    ensure_binaries_exist()
    corpus = corpus_files(args.corpus_dir, args.include_large)
    cases, skips = run_matrix(corpus, args.timeout)
    print(f"conformance decode matrix passed: {cases} case(s), {skips} skip(s)")
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except subprocess.TimeoutExpired as exc:
        print(f"command timed out after {exc.timeout}s: {' '.join(exc.cmd)}", file=sys.stderr)
        raise SystemExit(1)
    except RuntimeError as exc:
        print(str(exc), file=sys.stderr)
        raise SystemExit(1)
