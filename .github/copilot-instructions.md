# GitHub Copilot Instructions for CompressKit

CompressKit (`LessUp/compress-kit`) is a multi-language lossless compression
algorithm lab, not a generic application repository.

## Always preserve

- Product name: **CompressKit**
- Default branch: `master`
- Docs URL: <https://lessup.github.io/compress-kit/>
- Unified CLI: `<binary> <encode|decode> <input> <output>`
- Security limits: 4 GiB input, 1 GiB decoded output
- Per-algorithm binary compatibility across C++17, Go, and Rust

## Check specs first

Relevant requirements live in:

- `openspec/specs/encoding-project/spec.md`
- `openspec/specs/core-architecture/spec.md`
- `openspec/specs/cross-language-testing/spec.md`

Open a new OpenSpec change before altering binary formats, public API contracts,
CLI semantics, test-gate semantics, or adding algorithms.

## Validation commands

Use existing project commands only:

```bash
make build
make test
make test-conformance
make lint
npm run docs:build
openspec validate --all
```

For docs-only changes, `npm run docs:build` is the key check.

## Repository boundaries

- `algorithms/<algorithm>/<language>/`: implementation and local tests
- `algorithms/shared/`: streaming and buffer API foundation
- `tests/conformance/`: executable cross-language matrix
- `tests/streaming_api_contract/`: contract documentation and provenance
- `docs/`: VitePress user documentation
- `openspec/`: requirements and archived/deferred proposals

Do not commit generated binaries, `tests/data/*.bin`, Rust `target/`, reports,
or `docs/.vitepress/dist/`.

## Known limitation

Range Coder decode performance on files larger than 500 KiB is a documented
known issue. Do not attempt to fix it unless the task explicitly scopes that
performance work.

## Preferred implementation style

- Add or update tests before behavior changes.
- Use English error messages in code.
- Keep Python tooling deterministic and standard-library-first.
- Keep documentation concise and audience-specific; avoid boilerplate.
