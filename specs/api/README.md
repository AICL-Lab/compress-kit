# API Specifications

This directory contains API interface definitions (OpenAPI, GraphQL schemas, etc.).

## Current Status

**Not applicable for this project.**

The Encoding project is a command-line compression tool collection that does not expose HTTP/REST APIs or GraphQL endpoints.

## Future Considerations

If the project evolves to include:
- A web service for online compression
- gRPC/REST API for remote encoding/decoding
- WebAssembly bindings for browser use

Then this directory should contain:
- `openapi.yaml` - REST API definitions
- `graphql/schema.graphql` - GraphQL schema (if applicable)
- `protobuf/` - Protocol buffer definitions (if applicable)
