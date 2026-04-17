#!/usr/bin/env python3
import os
import sys
import time
import subprocess
from pathlib import Path

# Range coder 跨语言 benchmark 脚本
# - 编译 C++, Go, Rust 三种实现
# - 在指定输入文件上执行 encode/decode
# - 校验解码结果与原始输入是否一致
# - 比较三种语言的性能差异
#
# 用法：
#   python3 bench.py [/path/to/input.bin]
# 若未提供参数，则默认使用项目根目录下 tests/data/random_10MiB.bin

ROOT = Path(__file__).resolve().parent.parent
CPP_DIR = ROOT / "cpp"
GO_DIR = ROOT / "go"
RUST_DIR = ROOT / "rust"
BENCH_DIR = ROOT / "benchmark"
TMP_DIR = BENCH_DIR / "tmp"
PROJECT_ROOT = ROOT.parent
TEST_DATA_DIR = PROJECT_ROOT / "tests" / "data"


def run(cmd, cwd):
    start = time.perf_counter()
    proc = subprocess.run(cmd, cwd=cwd, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    end = time.perf_counter()
    if proc.returncode != 0:
        sys.stderr.write(f"Command failed: {' '.join(map(str, cmd))}\n")
        if proc.stderr:
            sys.stderr.write(proc.stderr.decode(errors="ignore") + "\n")
        sys.exit(1)
    return end - start


def ensure_tmp():
    TMP_DIR.mkdir(parents=True, exist_ok=True)


def compile_all():
    times = {}
    times["cpp_build"] = run(["g++", "-std=c++17", "-O2", "main.cpp", "-o", "rangecoder_cpp"], CPP_DIR)
    times["go_build"] = run(["go", "build", "-o", "rangecoder_go", "./cmd"], GO_DIR)
    times["rust_build"] = run(
        ["cargo", "build", "--bin", "rangecoder", "--release"],
        RUST_DIR,
    )
    return times


def bench_lang(name: str, exe: Path, input_path: Path, encoded_path: Path, decoded_path: Path, cwd: Path):
    enc_t = run([str(exe), "encode", str(input_path), str(encoded_path)], cwd=cwd)
    dec_t = run([str(exe), "decode", str(encoded_path), str(decoded_path)], cwd=cwd)
    size = encoded_path.stat().st_size
    return enc_t, dec_t, size


def files_equal(a: Path, b: Path) -> bool:
    with a.open("rb") as fa, b.open("rb") as fb:
        while True:
            ba = fa.read(65536)
            bb = fb.read(65536)
            if ba != bb:
                return False
            if not ba:
                return True


def main():
    ensure_tmp()
    if len(sys.argv) >= 2:
        input_path = Path(sys.argv[1]).resolve()
        if not input_path.is_file():
            sys.stderr.write("Input file does not exist\n")
            sys.exit(1)
    else:
        input_path = TEST_DATA_DIR / "random_10MiB.bin"

    if not input_path.is_file():
        sys.stderr.write(f"Input file not found: {input_path}\n")
        sys.stderr.write("请先运行 tests/gen_testdata.py 生成测试数据\n")
        sys.exit(1)

    build_times = compile_all()
    original_size = input_path.stat().st_size

    cpp_exe = CPP_DIR / "rangecoder_cpp"
    go_exe = GO_DIR / "rangecoder_go"
    rust_exe = RUST_DIR / "target" / "release" / "rangecoder"

    cpp_enc = TMP_DIR / "cpp.rcnc"
    cpp_dec = TMP_DIR / "cpp.out"
    go_enc = TMP_DIR / "go.rcnc"
    go_dec = TMP_DIR / "go.out"
    rust_enc = TMP_DIR / "rust.rcnc"
    rust_dec = TMP_DIR / "rust.out"

    results = []
    enc_t, dec_t, size = bench_lang("cpp", cpp_exe, input_path, cpp_enc, cpp_dec, CPP_DIR)
    results.append(("cpp", enc_t, dec_t, size))
    enc_t, dec_t, size = bench_lang("go", go_exe, input_path, go_enc, go_dec, GO_DIR)
    results.append(("go", enc_t, dec_t, size))
    enc_t, dec_t, size = bench_lang("rust", rust_exe, input_path, rust_enc, rust_dec, RUST_DIR)
    results.append(("rust", enc_t, dec_t, size))

    for name, _, _, _ in results:
        dec = TMP_DIR / f"{name}.out"
        if not files_equal(input_path, dec):
            sys.stderr.write(f"{name} decode output mismatch\n")
            sys.exit(1)

    print(f"Original size: {original_size} bytes")
    print("Build times (s):")
    for k, v in build_times.items():
        print(f"  {k}: {v:.4f}")

    print("\nRuntime (seconds) and compression ratio:")
    print("lang  encode  decode  total  ratio")
    for name, enc_t, dec_t, comp_size in results:
        total = enc_t + dec_t
        ratio = comp_size / original_size if original_size > 0 else 0.0
        print(f"{name:4}  {enc_t:6.4f}  {dec_t:6.4f}  {total:6.4f}  {ratio:6.3f}")


if __name__ == "__main__":
    main()
