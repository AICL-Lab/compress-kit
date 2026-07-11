## Summary

Explain what changed and why.

## Scope

- [ ] Algorithm implementation
- [ ] Shared buffer layer
- [ ] CI / repository tooling

Affected algorithms:

- Algorithm: Huffman / Arithmetic / Range Coder / RLE / Shared / N/A

## Compatibility

- [ ] No binary format change
- [ ] Binary format change (requires design review)

## Verification

Paste the commands run locally:

```bash
make test
make lint
```

## Review focus

Call out anything reviewers should inspect closely, especially:

- encoder/decoder pair compatibility
- Range Coder large-file behavior
- generated artifacts or ignored files
