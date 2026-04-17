# Project Philosophy: Spec-Driven Development (SDD)

This project strictly follows the **Spec-Driven Development (SDD)** paradigm. All code implementations must be based on the specification documents in the `/specs` directory as the Single Source of Truth.

## Directory Context

| Directory | Purpose |
|-----------|---------|
| `/specs/product/` | Product feature definitions and acceptance criteria |
| `/specs/rfc/` | Technical design documents (architecture, patterns, decisions) |
| `/specs/api/` | API interface definitions (OpenAPI, GraphQL schemas) |
| `/specs/db/` | Database model definitions |
| `/specs/testing/` | Test specifications and BDD test cases |
| `/docs/` | User-facing documentation (VitePress documentation site) |

## AI Agent Workflow Instructions

When you (AI) are asked to develop a new feature, modify an existing one, or fix a bug, **you must strictly follow this workflow without skipping any steps**:

### Step 1: Review Specs (审查规范)

- Before writing any code, first read the relevant documents in `/specs` (product requirements, RFCs, API definitions).
- If the user's request conflicts with existing specs, **stop immediately** and point out the conflict. Ask the user whether to update the spec first.
- 如果用户指令与现有 Spec 冲突，请立即停止编码，并指出冲突点，询问用户是否需要先更新 Spec。

### Step 2: Spec-First Update (规范优先)

- If this is a new feature, or if existing interfaces/database structures need to change, **you must first propose modifying or creating the appropriate spec documents** (e.g., `openapi.yaml` or RFC documents).
- Wait for user confirmation of the spec changes before proceeding to the code writing phase.
- 如果这是一个新功能，或者需要改变现有的接口/数据库结构，**必须首先提议修改或创建相应的 Spec 文档**。等待用户确认 Spec 的修改后，才能进入代码编写阶段。

### Step 3: Implementation (代码实现)

- When writing code, **100% adhere to the definitions in the specs** (including variable naming, API paths, data types, status codes, etc.).
- Do not add features in the code that are not defined in the specs (No Gold-Plating).
- 编写代码时，必须 100% 遵守 Spec 中的定义（包括变量命名、API 路径、数据类型、状态码等）。不要在代码中擅自添加 Spec 中未定义的功能。

### Step 4: Test Against Spec (测试验证)

- Write unit and integration tests based on the acceptance criteria in `/specs`.
- Ensure test cases cover all boundary conditions described in the specs.
- 根据 `/specs` 中的验收标准（Acceptance Criteria）编写单元测试和集成测试。确保测试用例覆盖了 Spec 中描述的所有边界情况。

## Code Generation Rules

- Any externally exposed API changes must synchronously update the relevant spec files in `/specs/api/`.
- If uncertain about technical details, consult the architecture conventions in `/specs/rfc/`. Do not invent design patterns on your own.
- All error messages must be in English.
- Follow language-specific conventions:
  - **C++17**: Google C++ Style Guide, 4-space indentation, snake_case for functions/variables, PascalCase for classes
  - **Go 1.21+**: gofmt formatting, go vet analysis, Effective Go conventions
  - **Rust 1.70+**: rustfmt formatting, clippy linting, Rust API Guidelines
  - **Python 3.8+**: PEP 8 style, 4-space indentation

## Why This Matters

1. **Prevent AI Hallucinations (防范 AI 幻觉)**: AI tends to "freestyle" without context. Forcing it to read `/specs` in the first step anchors its thinking scope.
2. **Constrain Modification Path (约束修改路径)**: Declaring "modify specs before code" ensures document-code synchronization (Document-Code Synchronization).
3. **Improve PR Quality (提高 PR 质量)**: When AI generates Pull Requests, the implementation will be highly aligned with business logic because it's developed based on the acceptance criteria defined in the specs.

## Project-Specific Notes

This project implements compression algorithms (Huffman, Arithmetic Coding, Range Coder, RLE) in multiple languages (C++17, Go, Rust). Key considerations:

- **Cross-Language Compatibility**: All implementations must share identical binary file formats
- **Unified CLI Interface**: `./binary <encode|decode> <input> <output>`
- **Security Constraints**: Max input 4 GiB, max output 1 GiB to prevent decompression bombs
- **Testing Focus**: Cross-language encode/decode verification is mandatory

For detailed product requirements, see [specs/product/encoding-project.md](specs/product/encoding-project.md).
For architecture decisions, see [specs/rfc/0001-core-architecture.md](specs/rfc/0001-core-architecture.md).
For testing specifications, see [specs/testing/cross-language.md](specs/testing/cross-language.md).
