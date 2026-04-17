# Specifications

This directory contains all specification documents for the Encoding project. The project follows **Spec-Driven Development (SDD)**, meaning all code implementations must be based on these specs as the Single Source of Truth.

## Directory Structure

```
specs/
├── product/            # Product requirements (PRD)
│   └── encoding-project.md
├── rfc/                # Technical design documents (RFCs)
│   └── 0001-core-architecture.md
├── api/                # API interface definitions
│   └── README.md       # (Not applicable for this CLI project)
├── db/                 # Database schemas
│   └── README.md       # (Not applicable for this CLI project)
└── testing/            # Test specifications
    └── cross-language.md
```

| Directory | Purpose | Status |
|-----------|---------|--------|
| `product/` | Product feature definitions and acceptance criteria | ✅ Active |
| `rfc/` | Technical design documents (architecture, patterns, decisions) | ✅ Active |
| `api/` | API interface definitions (OpenAPI, GraphQL schemas) | ⚪ N/A |
| `db/` | Database model definitions | ⚪ N/A |
| `testing/` | Test specifications and cross-language verification rules | ✅ Active |

## Current Specs

### Product Requirements
- [Encoding Project](product/encoding-project.md) - Project goals, algorithm implementations, file format compatibility, security and quality requirements

### Technical Design (RFCs)
- [RFC-0001: Core Architecture](rfc/0001-core-architecture.md) - Directory structure, CLI patterns, frequency table format, CI/CD workflow design, error handling strategy

### Testing Specifications
- [Cross-Language Testing](testing/cross-language.md) - Correctness tests, benchmark tests, known issues, future improvements

## SDD Workflow

This project follows Spec-Driven Development. The workflow is:

1. **Review Specs First** - Read relevant specs before coding
2. **Update Specs First** - Propose spec changes for new features before implementation
3. **Implement to Spec** - Code must 100% adhere to spec definitions
4. **Test Against Spec** - Write tests based on acceptance criteria

For complete AI workflow instructions, see [AGENTS.md](../AGENTS.md).

## Contributing to Specs

When adding new features or making architectural changes:

1. Create or update the relevant spec document first
2. Follow the existing naming conventions:
   - Product specs: `feature-name.md` in `product/`
   - RFCs: `NNNN-short-title.md` in `rfc/` (e.g., `0002-oauth2-implementation.md`)
   - Test specs: `feature-or-area.md` in `testing/`
3. Include clear acceptance criteria
4. Reference related specs where applicable
