import importlib.util
import json
import shutil
import sys
import textwrap
import unittest
from pathlib import Path


REPO_ROOT = Path(__file__).resolve().parent.parent
SCRIPT_PATH = REPO_ROOT / "scripts" / "run_all_bench.py"
WORKSPACE = REPO_ROOT / "tests" / "_run_all_bench_workspace"


def load_module():
    spec = importlib.util.spec_from_file_location("run_all_bench_under_test", SCRIPT_PATH)
    module = importlib.util.module_from_spec(spec)
    assert spec.loader is not None
    spec.loader.exec_module(module)
    return module


class RunAllBenchTests(unittest.TestCase):
    def setUp(self):
        if WORKSPACE.exists():
            shutil.rmtree(WORKSPACE)
        WORKSPACE.mkdir(parents=True)
        self.root = WORKSPACE / "repo"
        self.root.mkdir()
        self.module = load_module()
        self._write_gen_testdata()
        self._write_docs_json({"generated": "stale", "version": "1.0.0", "results": []})
        self._patch_module_paths()

    def tearDown(self):
        if WORKSPACE.exists():
            shutil.rmtree(WORKSPACE)

    def test_main_fails_when_any_benchmark_driver_is_missing(self):
        self._write_driver("huffman", sample_output("textlike_10MiB"))
        self._write_driver("arithmetic", sample_output("textlike_10MiB"))
        self._write_driver("range", sample_output("small_dictionary_like"))

        with self.assertRaisesRegex(SystemExit, "missing benchmark driver"):
            self.module.main()

    def test_main_fails_when_a_benchmark_driver_exits_non_zero(self):
        self._write_driver("huffman", sample_output("textlike_10MiB"))
        self._write_driver("arithmetic", sample_output("textlike_10MiB"))
        self._write_driver("range", sample_output("small_dictionary_like"))
        self._write_driver("rle", "import sys\nsys.stderr.write('boom\\n')\nsys.exit(7)\n")

        with self.assertRaisesRegex(SystemExit, "exit code 7"):
            self.module.main()

    def test_main_refreshes_docs_benchmark_json_from_driver_output(self):
        self._write_driver("huffman", sample_output("textlike_10MiB"))
        self._write_driver("arithmetic", sample_output("textlike_10MiB"))
        self._write_driver("range", sample_output("small_dictionary_like"))
        self._write_driver("rle", sample_output("repetitive_10MiB"))

        self.module.main()

        docs_payload = json.loads(
            (self.root / "docs" / ".vitepress" / "data" / "benchmarks.json").read_text(encoding="utf-8")
        )
        self.assertNotEqual(docs_payload["generated"], "stale")
        self.assertEqual(len(docs_payload["results"]), 12)

        by_algorithm = {result["algorithm"]: result for result in docs_payload["results"] if result["language"] == "cpp"}
        self.assertEqual(by_algorithm["huffman"]["dataset"], "textlike_10MiB")
        self.assertEqual(by_algorithm["range"]["dataset"], "small_dictionary_like")
        self.assertEqual(by_algorithm["rle"]["dataset"], "repetitive_10MiB")
        self.assertAlmostEqual(by_algorithm["huffman"]["encodeTime"], 100.0)
        self.assertAlmostEqual(by_algorithm["huffman"]["decodeTime"], 200.0)
        self.assertAlmostEqual(by_algorithm["huffman"]["encodeSpeed"], 100.0)
        self.assertAlmostEqual(by_algorithm["huffman"]["decodeSpeed"], 50.0)
        self.assertAlmostEqual(by_algorithm["huffman"]["compressionRatio"], 2.0)

        report_paths = sorted((self.root / "reports").glob("*_report_*.txt"))
        self.assertEqual(len(report_paths), 4)

    def _patch_module_paths(self):
        self.module.ROOT = self.root
        self.module.TESTS_DIR = self.root / "tests"
        self.module.TEST_DATA_DIR = self.module.TESTS_DIR / "data"
        self.module.REPORTS_DIR = self.root / "reports"
        self.module.PY = sys.executable

    def _write_docs_json(self, payload):
        target = self.root / "docs" / ".vitepress" / "data" / "benchmarks.json"
        target.parent.mkdir(parents=True, exist_ok=True)
        target.write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")

    def _write_gen_testdata(self):
        script = self.root / "tests" / "gen_testdata.py"
        script.parent.mkdir(parents=True, exist_ok=True)
        script.write_text(
            textwrap.dedent(
                """
                from pathlib import Path

                root = Path(__file__).resolve().parent / "data"
                root.mkdir(parents=True, exist_ok=True)
                files = {
                    "random_10MiB.bin": 10 * 1024 * 1024,
                    "random_1MiB.bin": 1024 * 1024,
                    "textlike_10MiB.bin": 10 * 1024 * 1024,
                    "repetitive_10MiB.bin": 10 * 1024 * 1024,
                    "small_dictionary_like.bin": 4096,
                }
                for name, size in files.items():
                    (root / name).write_bytes(bytes([len(name) % 251]) * size)
                print("[gen_testdata] done")
                """
            ).strip()
            + "\n",
            encoding="utf-8",
        )

    def _write_driver(self, algorithm, body):
        bench = self.root / "algorithms" / algorithm / "benchmark" / "bench.py"
        bench.parent.mkdir(parents=True, exist_ok=True)
        bench.write_text(body, encoding="utf-8")


def sample_output(dataset_name):
    return textwrap.dedent(
        f"""
        import sys

        dataset = {dataset_name!r}
        _ = sys.argv[1]
        print("Dataset: " + dataset)
        print("Original size: 10485760 bytes")
        print("Build times (s):")
        print("  cpp_build: 0.0100")
        print()
        print("Runtime (seconds) and compression ratio:")
        print("lang  encode  decode  total  ratio")
        print("cpp   0.1000  0.2000  0.3000  0.500")
        print("go    0.1250  0.2500  0.3750  0.400")
        print("rust  0.0800  0.1600  0.2400  0.250")
        """
    ).strip() + "\n"


if __name__ == "__main__":
    unittest.main()
