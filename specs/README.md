# Specifications

This directory contains all specification documents for the Encoding project. The project follows **Spec-Driven Development (SDD)**, meaning all code implementations must be based on these specs as the Single Source of Truth.

## Directory Structure

| Directory | Purpose |
|-----------|---------|
| `product/` | Product feature definitions and acceptance criteria |
| `rfc/` | Technical design documents (architecture, patterns, decisions) |
| `api/` | API interface definitions (OpenAPI, schemas) |
| `db/` | Database model definitions (if applicable) |
| `testing/` | Test specifications and cross-language verification rules |

## Current Specs

- [Encoding Project](product/encoding-project.md) - Product requirements
- [RFC-0001: Core Architecture](rfc/0001-core-architecture.md) - Architecture design
- [Cross-Language Testing](testing/cross-language.md) - Test specifications

## AI Agent Instructions

For AI assistants working on this project, see [AGENTS.md](../AGENTS.md) for the complete SDD workflow instructions.
