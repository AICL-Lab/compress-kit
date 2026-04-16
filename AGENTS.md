# Project Philosophy: Spec-Driven Development (SDD)

This project strictly follows the **Spec-Driven Development (SDD)** paradigm. All code implementations are based on the specification documents in the `/specs` directory as the Single Source of Truth.

## Directory Context

| Directory | Purpose |
|-----------|---------|
| `/specs/product/` | Product feature definitions and acceptance criteria |
| `/specs/rfc/` | Technical design documents (architecture, patterns, decisions) |
| `/specs/api/` | API interface definitions (OpenAPI, schemas) |
| `/specs/db/` | Database model definitions (if applicable) |
| `/specs/testing/` | Test specifications and cross-language verification rules |
| `/docs/` | User-facing documentation (VitePress site) |

## AI Agent Workflow Instructions

When you (AI) are asked to develop a new feature, modify an existing one, or fix a bug, **you must strictly follow this workflow without skipping any steps**:

### Step 1: Review Specs

- First, read the relevant documents in `/specs` (product requirements, RFCs, API definitions).
- If the user's request conflicts with existing specs, **stop immediately** and point out the conflict. Ask the user whether to update the spec first.

### Step 2: Spec-First Update

- If this is a new feature, or if existing interfaces/database structures need to change, **you must first propose modifying or creating the appropriate spec documents** (e.g., `openapi.yaml` or RFC documents).
- Wait for user confirmation of the spec changes before proceeding to the code writing phase.

### Step 3: Implementation

- When writing code, **100% adhere to the definitions in the specs** (including variable naming, API paths, data types, status codes, etc.).
- Do not add features in the code that are not defined in the specs (No Gold-Plating).

### Step 4: Test Against Spec

- Write unit and integration tests based on the acceptance criteria in `/specs`.
- Ensure test cases cover all boundary conditions described in the specs.

## Code Generation Rules

- Any externally exposed API changes must synchronously update the relevant spec files in `/specs/`.
- If uncertain about technical details, consult the architecture conventions in `/specs/rfc/`. Do not invent design patterns on your own.
- All error messages must be in English.
- Follow language-specific conventions (C++17, Go 1.21+, Rust 1.70+).

## Why This Matters

1. **Prevent AI Hallucinations**: AI tends to "freestyle" without context. Forcing it to read `/specs` in the first step anchors its thinking scope.
2. **Constrain Modification Path**: Declaring "modify specs before code" ensures document-code synchronization.
3. **Improve PR Quality**: When AI generates Pull Requests, the implementation will be highly aligned with business logic because it's developed based on the acceptance criteria defined in the specs.
