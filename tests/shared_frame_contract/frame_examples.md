# Shared Frame Contract Examples

## Minimum Frame Layout

`magic | version | algo_id | flags | content_size | checksum | [dictionary_id] | [metadata] | payload`

## Negative Cases

- bad magic
- unsupported version
- unknown algo_id
- checksum mismatch
- truncated header
- truncated payload
