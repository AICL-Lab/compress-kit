#!/usr/bin/env python3
import datetime
import json
import re
import subprocess
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parent.parent
TESTS_DIR = ROOT / "tests"
TEST_DATA_DIR = TESTS_DIR / "data"
REPORTS_DIR = ROOT / "reports"

PY = sys.executable or "python3"

ALGORITHM_ORDER = ("huffman", "arithmetic", "range", "rle")
LANGUAGE_ORDER = ("cpp", "go", "rust")
BENCHMARK_ROW_RE = re.compile(
    r"^(cpp|go|rust)\s+([0-9]+(?:\.[0-9]+)?)\s+([0-9]+(?:\.[0-9]+)?)\s+([0-9]+(?:\.[0-9]+)?)\s+([0-9]+(?:\.[0-9]+)?)$",
    re.MULTILINE,
)
ORIGINAL_SIZE_RE = re.compile(r"^Original size:\s+([0-9]+)\s+bytes$", re.MULTILINE)
DATASET_RE = re.compile(r"^Dataset:\s+(.+)$", re.MULTILINE)


def benchmark_jobs():
    return [
        {
            "algorithm": "huffman",
            "title": "Huffman benchmark",
            "driver": ROOT / "algorithms" / "huffman" / "benchmark" / "bench.py",
            "input": TEST_DATA_DIR / "textlike_10MiB.bin",
        },
        {
            "algorithm": "arithmetic",
            "title": "Arithmetic benchmark",
            "driver": ROOT / "algorithms" / "arithmetic" / "benchmark" / "bench.py",
            "input": TEST_DATA_DIR / "textlike_10MiB.bin",
        },
        {
            "algorithm": "range",
            "title": "Range coder benchmark",
            "driver": ROOT / "algorithms" / "range" / "benchmark" / "bench.py",
            "input": TEST_DATA_DIR / "small_dictionary_like.bin",
        },
        {
            "algorithm": "rle",
            "title": "RLE benchmark",
            "driver": ROOT / "algorithms" / "rle" / "benchmark" / "bench.py",
            "input": TEST_DATA_DIR / "repetitive_10MiB.bin",
        },
    ]


def docs_benchmarks_json_path():
    return ROOT / "docs" / ".vitepress" / "data" / "benchmarks.json"


def discover_benchmark_jobs():
    jobs = benchmark_jobs()
    missing = [job["algorithm"] for job in jobs if not job["driver"].is_file()]
    if missing:
        raise SystemExit(f"missing benchmark driver(s): {', '.join(missing)}")
    return jobs


def run_capture(cmd, cwd, report_path, title):
    print(f"[run_all_bench] running {title}: {' '.join(map(str, cmd))} (cwd={cwd})")
    proc = subprocess.run(
        cmd,
        cwd=cwd,
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
        text=True,
    )
    REPORTS_DIR.mkdir(parents=True, exist_ok=True)
    header = (
        f"# {title}\n"
        f"# command: {' '.join(map(str, cmd))}\n"
        f"# cwd: {cwd}\n"
        f"# returncode: {proc.returncode}\n\n"
    )
    report_path.write_text(header + proc.stdout, encoding="utf-8", errors="ignore")
    if proc.returncode != 0:
        raise SystemExit(f"{title} failed with exit code {proc.returncode}; see {report_path}")

    print(f"[run_all_bench] {title} OK, report -> {report_path}")
    return proc.stdout


def ensure_testdata():
    print("[run_all_bench] generating test data...")
    subprocess.check_call([PY, "tests/gen_testdata.py"], cwd=ROOT)
    jobs = benchmark_jobs()
    missing_inputs = sorted({job["input"] for job in jobs if not job["input"].is_file()}, key=str)
    if missing_inputs:
        missing_list = ", ".join(str(path.relative_to(ROOT)) for path in missing_inputs)
        raise SystemExit(f"benchmark input file(s) not found: {missing_list}")
    return jobs


def parse_benchmark_output(algorithm, dataset, output):
    original_size_match = ORIGINAL_SIZE_RE.search(output)
    if original_size_match is None:
        raise SystemExit(f"{algorithm} benchmark output missing original size")

    original_size = int(original_size_match.group(1))
    if original_size <= 0:
        raise SystemExit(f"{algorithm} benchmark output has invalid original size: {original_size}")

    dataset_match = DATASET_RE.search(output)
    dataset_name = dataset_match.group(1).strip() if dataset_match else dataset

    results = []
    seen_languages = set()
    input_mib = original_size / (1024 * 1024)

    for match in BENCHMARK_ROW_RE.finditer(output):
        language = match.group(1)
        encode_seconds = float(match.group(2))
        decode_seconds = float(match.group(3))
        ratio = float(match.group(5))
        seen_languages.add(language)

        encode_speed = compute_speed(input_mib, encode_seconds)
        decode_speed = compute_speed(input_mib, decode_seconds)
        results.append(
            {
                "algorithm": algorithm,
                "language": language,
                "dataset": dataset_name,
                "encodeTime": round(encode_seconds * 1000, 3),
                "decodeTime": round(decode_seconds * 1000, 3),
                "encodeSpeed": round(encode_speed, 1),
                "decodeSpeed": round(decode_speed, 1),
                "compressionRatio": round(ratio, 3),
                "throughput": classify_throughput((encode_speed + decode_speed) / 2.0),
            }
        )

    missing_languages = [language for language in LANGUAGE_ORDER if language not in seen_languages]
    if missing_languages:
        raise SystemExit(f"{algorithm} benchmark output missing language row(s): {', '.join(missing_languages)}")

    return sorted(results, key=lambda item: LANGUAGE_ORDER.index(item["language"]))


def compute_speed(size_mib, seconds):
    if seconds <= 0:
        return 0.0
    return size_mib / seconds
def classify_throughput(speed_mib_per_second):
    if speed_mib_per_second >= 200:
        return "very-high"
    if speed_mib_per_second >= 75:
        return "high"
    if speed_mib_per_second >= 25:
        return "medium"
    return "low"


def emit_docs_benchmark_data(results):
    target = docs_benchmarks_json_path()
    payload = {
        "generated": datetime.date.today().isoformat(),
        "version": "1.0.0",
        "results": sorted(
            results,
            key=lambda item: (
                ALGORITHM_ORDER.index(item["algorithm"]),
                LANGUAGE_ORDER.index(item["language"]),
            ),
        ),
    }
    target.parent.mkdir(parents=True, exist_ok=True)
    target.write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")
    print(f"[run_all_bench] refreshed docs benchmark data -> {target}")


def main():
    REPORTS_DIR.mkdir(parents=True, exist_ok=True)
    ensure_testdata()
    jobs = discover_benchmark_jobs()
    timestamp = datetime.datetime.now().strftime("%Y%m%d-%H%M%S")

    docs_results = []
    report_paths = []
    for job in jobs:
        report_path = REPORTS_DIR / f"{job['algorithm']}_report_{timestamp}.txt"
        stdout = run_capture(
            [PY, str(job["driver"]), str(job["input"])],
            cwd=job["driver"].parent,
            report_path=report_path,
            title=job["title"],
        )
        docs_results.extend(parse_benchmark_output(job["algorithm"], job["input"].stem, stdout))
        report_paths.append(report_path)

    emit_docs_benchmark_data(docs_results)

    print("[run_all_bench] all done. Reports are in:")
    for path in report_paths:
        print("  ", path)


if __name__ == "__main__":
    main()
