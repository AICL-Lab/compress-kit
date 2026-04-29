## Summary

Explain what changed and why.

## Scope

- [ ] Algorithm implementation
- [ ] Shared streaming/buffer layer
- [ ] Cross-language conformance
- [ ] Documentation / Git Pages
- [ ] CI / repository tooling
- [ ] OpenSpec

Affected algorithms/languages:

- Algorithm: Huffman / Arithmetic / Range Coder / RLE / Shared / N/A
- Language: C++17 / Go / Rust / Python scripts / Docs / CI

## Compatibility

- [ ] No binary format change
- [ ] Binary format change, OpenSpec updated
- [ ] Cross-language behavior changed, `make test-conformance` updated

## Verification

Paste the commands run locally:

```bash
make test
npm run docs:build
```

## Review focus

Call out anything reviewers should inspect closely, especially:

- encoder/decoder pair compatibility
- Range Coder large-file behavior
- generated artifacts or ignored files
- OpenSpec requirement alignment
