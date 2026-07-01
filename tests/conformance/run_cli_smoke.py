#!/usr/bin/env python3
"""Run fast CLI smoke checks for every shipped CompressKit binary."""

from __future__ import annotations

import subprocess
import sys
import tempfile
from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]
USAGE_FRAGMENT = "encode|decode input output"
INVALID_MODE = "invalid"
TIMEOUT_SECONDS = 10.0

ALGORITHMS = {
    "huffman": [ROOT / "build/huffman_cpp"],
    "arithmetic": [ROOT / "build/arithmetic_cpp"],
    "range": [ROOT / "build/rangecoder_cpp"],
    "rle": [ROOT / "build/rle_cpp"],
}

CORPUS = (
    ROOT / "tests/data/empty.bin",
    ROOT / "tests/data/single_byte.bin",
    ROOT / "tests/data/alternating.bin",
    ROOT / "tests/data/all_same_byte.bin",
    ROOT / "tests/data/small_dictionary_like.bin",
)


def ensure_files_exist(paths: list[Path] | tuple[Path, ...], label: str, hint: str) -> None:
    missing = [path for path in paths if not path.is_file()]
    if not missing:
        return
    rendered = "\n".join(str(path.relative_to(ROOT)) for path in missing)
    raise SystemExit(f"missing {label}; run `{hint}` first:\n{rendered}")


def run(command: list[str]) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        command,
        cwd=ROOT,
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        timeout=TIMEOUT_SECONDS,
        check=False,
    )


def run_checked(command: list[str]) -> None:
    proc = run(command)
    if proc.returncode == 0:
        return
    rendered = " ".join(command)
    raise RuntimeError(
        f"command failed with exit code {proc.returncode}: {rendered}\n"
        f"stdout:\n{proc.stdout}\n"
        f"stderr:\n{proc.stderr}"
    )


def assert_usage(binary: Path) -> None:
    proc = run([str(binary)])
    assert_usage_error(binary, proc, "without args")


def assert_usage_error(
    binary: Path, proc: subprocess.CompletedProcess[str], scenario: str
) -> None:
    combined = proc.stdout + proc.stderr
    if proc.returncode == 0:
        raise RuntimeError(f"{binary} unexpectedly succeeded {scenario}")
    if "Usage:" not in combined:
        raise RuntimeError(f"{binary} did not print usage {scenario}")
    if USAGE_FRAGMENT not in combined:
        raise RuntimeError(f"{binary} did not advertise the unified CLI contract {scenario}")


def assert_wrong_arity(binary: Path, args: list[str], scenario: str) -> None:
    proc = run([str(binary), *args])
    assert_usage_error(binary, proc, scenario)


def assert_invalid_mode(binary: Path, source: Path, output: Path) -> None:
    proc = run([str(binary), INVALID_MODE, str(source), str(output)])
    combined = proc.stdout + proc.stderr
    lowered = combined.lower()
    if proc.returncode == 0:
        raise RuntimeError(f"{binary} unexpectedly accepted invalid mode {INVALID_MODE!r}")
    if "mode" not in lowered:
        raise RuntimeError(f"{binary} did not explain invalid mode handling")
    if "encode" not in lowered or "decode" not in lowered:
        raise RuntimeError(f"{binary} did not advertise supported modes on invalid mode")


def assert_round_trip(binary: Path, source: Path, encoded: Path, decoded: Path) -> None:
    run_checked([str(binary), "encode", str(source), str(encoded)])
    run_checked([str(binary), "decode", str(encoded), str(decoded)])
    if source.read_bytes() != decoded.read_bytes():
        raise RuntimeError(f"{binary} round-trip mismatch for {source.name}")


def main() -> int:
    binaries = [binary for group in ALGORITHMS.values() for binary in group]
    ensure_files_exist(tuple(binaries), "binary file(s)", "make build")
    ensure_files_exist(CORPUS, "corpus file(s)", "make test-data")

    checks = 0
    for algorithm, algorithm_binaries in ALGORITHMS.items():
        for binary in algorithm_binaries:
            assert_usage(binary)
            checks += 1
            print(f"PASS usage {algorithm} {binary.name}", flush=True)

    with tempfile.TemporaryDirectory(prefix=".cli-smoke-", dir=ROOT / "tests") as tmp:
        tmpdir = Path(tmp)
        for algorithm, algorithm_binaries in ALGORITHMS.items():
            for binary in algorithm_binaries:
                assert_wrong_arity(
                    binary,
                    ["encode", str(CORPUS[0])],
                    "with too few args",
                )
                checks += 1
                print(f"PASS wrong-arity-too-few {algorithm} {binary.name}", flush=True)
                extra_output = tmpdir / f"{algorithm}-{binary.name}.extra"
                assert_wrong_arity(
                    binary,
                    ["encode", str(CORPUS[0]), str(extra_output), "extra"],
                    "with too many args",
                )
                checks += 1
                print(f"PASS wrong-arity-too-many {algorithm} {binary.name}", flush=True)
                invalid_output = tmpdir / f"{algorithm}-{binary.name}.invalid"
                assert_invalid_mode(binary, CORPUS[0], invalid_output)
                checks += 1
                print(f"PASS invalid-mode {algorithm} {binary.name}", flush=True)
                for source in CORPUS:
                    encoded = tmpdir / f"{algorithm}-{binary.name}-{source.name}.encoded"
                    decoded = tmpdir / f"{algorithm}-{binary.name}-{source.name}.decoded"
                    assert_round_trip(binary, source, encoded, decoded)
                    checks += 1
                    print(
                        f"PASS round-trip {algorithm} {binary.name} {source.name}",
                        flush=True,
                    )

    print(f"cli smoke passed: {checks} check(s)")
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
