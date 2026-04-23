# Project Philosophy: OpenSpec-Driven Development

This project uses **OpenSpec** for spec-driven development. All code implementations must be based on specifications in the `openspec/specs/` directory as the Single Source of Truth.

## OpenSpec Setup

### Installation

```bash
# Requires Node.js 20.19.0+
npm install -g @fission-ai/openspec@latest
```

### Verification

```bash
openspec --version
openspec list
```

## Directory Structure

| Directory | Purpose |
|-----------|---------|
| `openspec/specs/` | Main specification documents |
| `openspec/changes/` | Active change proposals |
| `openspec/changes/archive/` | Completed and archived changes |
| `openspec/config.yaml` | OpenSpec configuration |
| `docs/` | User-facing documentation (VitePress) |

## AI Agent Workflow

When developing a new feature, fixing a bug, or making architectural changes, follow this workflow:

### Step 1: Review Specs (审查规范)

Before any code changes:
```bash
# Check current specs
openspec list

# View specific spec
cat openspec/specs/<capability>/spec.md
```

- Read relevant specs in `openspec/specs/`
- If user request conflicts with specs, **stop and ask** whether to create a change proposal
- 如果用户指令与现有 Spec 冲突，请立即停止，询问是否需要创建变更提案

### Step 2: Create Proposal (创建提案)

For new features or significant changes, use OpenSpec workflow:

```bash
# Start a new change proposal
/opsx:propose "add-lz77-compression"
```

This creates `openspec/changes/add-lz77-compression/` with:
- `proposal.md` - Rationale and scope
- `specs/` - Delta specs (ADDED/MODIFIED/REMOVED requirements)
- `design.md` - Technical approach
- `tasks.md` - Implementation checklist

等待用户确认提案后再进入实现阶段。

### Step 3: Implement (代码实现)

After proposal approval:

```bash
# Start implementing tasks
/opsx:apply add-lz77-compression
```

- Implement tasks from `tasks.md` in order
- Mark tasks complete: `- [ ]` → `- [x]`
- Follow requirement keywords: `SHALL`, `MUST`
- Keep code changes minimal and focused

### Step 4: Test Against Spec (测试验证)

Write tests based on scenarios in specs:

```bash
# Run cross-language tests
make test

# Run benchmarks
make bench
```

Each requirement has scenarios in GIVEN/WHEN/THEN format - ensure tests cover these.

### Step 5: Archive (归档)

When implementation is complete:

```bash
# Archive the change
/opsx:archive add-lz77-compression
```

This:
1. Syncs delta specs to main specs (if any)
2. Moves change to `openspec/changes/archive/YYYY-MM-DD-add-lz77-compression/`

## OpenSpec Commands Reference

| Command | Purpose |
|---------|---------|
| `/opsx:propose <name>` | Create new change proposal with all artifacts |
| `/opsx:apply <name>` | Implement tasks from a change |
| `/opsx:archive <name>` | Archive completed change |
| `/opsx:sync <name>` | Sync delta specs to main specs without archiving |

## Code Generation Rules

- Binary format changes MUST update `openspec/specs/core-architecture/spec.md`
- New algorithms MUST update `openspec/specs/encoding-project/spec.md`
- Test changes MUST update `openspec/specs/cross-language-testing/spec.md`
- All error messages MUST be in English
- Follow language-specific conventions:
  - **C++17**: Google Style Guide, snake_case, PascalCase classes
  - **Go 1.21+**: gofmt, go vet, Effective Go
  - **Rust 1.70+**: rustfmt, clippy, Rust API Guidelines
  - **Python 3.8+**: PEP 8

## Why OpenSpec Matters

1. **Structured Change Management**: Each change is a self-contained unit with proposal, specs, design, and tasks
2. **Delta Spec Tracking**: Track what requirements are ADDED/MODIFIED/REMOVED
3. **Archive History**: Completed changes preserved with timestamps
4. **AI-Native Workflow**: Commands designed for AI agent execution

## Project-Specific Notes

This project implements compression algorithms (Huffman, Arithmetic Coding, Range Coder, RLE) in multiple languages (C++17, Go, Rust). Key considerations:

- **Cross-Language Compatibility**: All implementations must share identical binary file formats
- **Unified CLI Interface**: `./binary <encode|decode> <input> <output>`
- **Security Constraints**: Max input 4 GiB, max output 1 GiB to prevent decompression bombs
- **Testing Focus**: Cross-language encode/decode verification is mandatory

## Current Specs

| Spec | Description |
|------|-------------|
| [encoding-project](openspec/specs/encoding-project/spec.md) | Product requirements, algorithm status, security/quality requirements |
| [core-architecture](openspec/specs/core-architecture/spec.md) | Directory structure, CLI patterns, frequency table format, CI/CD design |
| [cross-language-testing](openspec/specs/cross-language-testing/spec.md) | Correctness tests, benchmarks, known issues |
