# TDD Execution Rules

The required loop for every new CompressKit capability is:

`spec delta -> failing test -> minimal code -> green -> refactor -> docs -> bench`

## Usage Notes

- Do not start parallel language ports until the reference implementation has passing tests.
- Keep failing tests scoped to one contract change at a time.
- Treat docs and benchmark work as post-green activities unless they are required to define the contract itself.
