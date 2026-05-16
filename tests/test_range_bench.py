import contextlib
import importlib.util
import io
import shutil
import sys
import unittest
from pathlib import Path


REPO_ROOT = Path(__file__).resolve().parent.parent
SCRIPT_PATH = REPO_ROOT / "algorithms" / "range" / "benchmark" / "bench.py"
WORKSPACE = REPO_ROOT / "tests" / "_range_bench_workspace"


def load_module():
    spec = importlib.util.spec_from_file_location("range_bench_under_test", SCRIPT_PATH)
    module = importlib.util.module_from_spec(spec)
    assert spec.loader is not None
    spec.loader.exec_module(module)
    return module


class RangeBenchTests(unittest.TestCase):
    def setUp(self):
        if WORKSPACE.exists():
            shutil.rmtree(WORKSPACE)
        self.module = load_module()
        self.expected_default = REPO_ROOT / "tests" / "data" / "small_dictionary_like.bin"
        self.module.TMP_DIR = WORKSPACE / "tmp"
        self.module.ensure_tmp = lambda: self.module.TMP_DIR.mkdir(parents=True, exist_ok=True)
        self.module.compile_all = lambda: {"build_all": 0.0}
        self.module.files_equal = lambda _a, _b: True

    def tearDown(self):
        if WORKSPACE.exists():
            shutil.rmtree(WORKSPACE)

    def test_default_input_points_to_shipped_small_fixture(self):
        self.assertEqual(self.module.DEFAULT_INPUT, self.expected_default)
        self.assertTrue(self.module.DEFAULT_INPUT.is_file())

    def test_main_without_args_uses_small_range_fixture(self):
        seen_inputs = []

        def fake_bench_lang(_name, _exe, input_path, _encoded_path, _decoded_path, _cwd):
            seen_inputs.append(input_path)
            return 0.01, 0.02, 32

        self.module.bench_lang = fake_bench_lang
        stdout = io.StringIO()
        original_argv = sys.argv
        try:
            sys.argv = [str(SCRIPT_PATH)]
            with contextlib.redirect_stdout(stdout):
                self.module.main()
        finally:
            sys.argv = original_argv

        self.assertEqual(seen_inputs, [self.module.DEFAULT_INPUT] * 3)
        self.assertIn("Dataset: small_dictionary_like", stdout.getvalue())


if __name__ == "__main__":
    unittest.main()
