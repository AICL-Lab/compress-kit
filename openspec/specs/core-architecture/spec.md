# Core Architecture Specification

## Overview

Technical architecture for the multi-language encoding project, defining directory structure, component patterns, and system design decisions.

## Purpose

Establish consistent architecture across C++17, Go, and Rust implementations while maintaining cross-language binary compatibility.
## Requirements

### REQ-ARCH-001: Directory Structure

The project SHALL follow algorithm-first organization with language subdirectories.

#### Scenario: Algorithm directory layout
- **GIVEN** any algorithm (huffman, arithmetic, range, rle)
- **WHEN** checking directory structure
- **THEN** algorithm directory SHALL contain cpp/, go/, and rust/ subdirectories

#### Scenario: Shared test location
- **GIVEN** project root
- **WHEN** looking for test data
- **THEN** tests/ directory SHALL contain gen_testdata.py and data/

### REQ-ARCH-002: CLI Pattern Consistency

All implementations SHALL follow identical CLI argument pattern.

#### Scenario: C++ CLI pattern
- **GIVEN** C++ implementation
- **WHEN** processing command line
- **THEN** main(argc, argv) SHALL accept: mode, input, output

#### Scenario: Go CLI pattern
- **GIVEN** Go implementation
- **WHEN** processing command line
- **THEN** os.Args SHALL provide: mode, input, output

#### Scenario: Rust CLI pattern
- **GIVEN** Rust implementation
- **WHEN** processing command line
- **THEN** env::args() SHALL provide: mode, input, output

### REQ-ARCH-003: Frequency Table Format

All static-model algorithms SHALL use consistent frequency table binary format.

#### Scenario: Table header
- **GIVEN** frequency table encoding
- **WHEN** writing to output
- **THEN** first 4 bytes SHALL be symbol count (uint32 LE)

#### Scenario: Frequency entries
- **GIVEN** frequency table with 256 symbols
- **WHEN** encoding
- **THEN** each frequency SHALL be 4 bytes LE (uint32)
- **AND** EOF frequency SHALL be final entry

### REQ-ARCH-004: Error Handling

All implementations SHALL follow consistent error handling strategy.

#### Scenario: Input validation
- **GIVEN** any input file
- **WHEN** opening file
- **THEN** system SHALL check file size before reading

#### Scenario: Error messages
- **GIVEN** any error condition
- **WHEN** reporting error
- **THEN** message SHALL be in English and actionable

#### Scenario: Exit codes
- **GIVEN** program completion
- **WHEN** exiting
- **THEN** exit code 0 SHALL indicate success, 1 SHALL indicate error

### REQ-ARCH-005: Performance Optimizations

Implementations SHALL use performance-optimized I/O patterns.

#### Scenario: Buffered I/O
- **GIVEN** file processing
- **WHEN** reading/writing
- **THEN** buffered I/O SHALL be used (bufio, ifstream/ofstream)

#### Scenario: Frequency scaling
- **GIVEN** arithmetic or range coding
- **WHEN** building frequency table
- **THEN** frequencies SHALL scale to maxTotal (2^24) for stability

### REQ-ARCH-006: CI/CD Workflow

The project SHALL implement multi-stage CI/CD pipeline.

#### Scenario: Build stage
- **GIVEN** CI trigger (push/PR)
- **WHEN** build stage runs
- **THEN** C++, Go, Rust builds SHALL execute in parallel

#### Scenario: Test stage
- **GIVEN** successful builds
- **WHEN** test stage runs
- **THEN** cross-language correctness tests SHALL execute

### REQ-ARCH-007: Documentation Architecture

Documentation SHALL be built using VitePress.

#### Scenario: VitePress configuration
- **GIVEN** docs/.vitepress/config.mts
- **WHEN** building documentation
- **THEN** site SHALL include guide/, public/, and index.md

#### Scenario: Algorithm documentation
- **GIVEN** documentation site
- **WHEN** navigating to algorithms section
- **THEN** complexity analysis and implementation notes SHALL be available

### REQ-ARCH-008: Development Tool Configuration

The project SHALL maintain development tool configurations.

#### Scenario: Code formatting
- **GIVEN** project root
- **WHEN** checking for formatting configs
- **THEN** .clang-format, gofmt (via go vet), and rustfmt SHALL be configured

#### Scenario: Dependency management
- **GIVEN** project dependencies
- **WHEN** checking for lock files
- **THEN** requirements.txt (Python), go.sum, and Cargo.lock SHALL exist

#### Scenario: Security scanning
- **GIVEN** GitHub workflows
- **WHEN** checking for security tools
- **THEN** CodeQL SHALL be configured
- **AND** Dependabot is deliberately disabled to maintain a clean, single-branch finalized state

### Requirement: REQ-ARCH-009 Streaming Encoder/Decoder Interface

All algorithm implementations SHALL expose a streaming interface with a defined lifecycle.

#### Scenario: Streaming encoder lifecycle
- **GIVEN** a freshly constructed encoder
- **WHEN** caller calls `process(chunk)` one or more times, then `finish()`
- **THEN** all input bytes SHALL be encoded and emitted across the `process` and `finish` calls in order
- **AND** state SHALL transition READY → STREAMING → FINISHED

#### Scenario: Flush produces stable output
- **GIVEN** an encoder in STREAMING state
- **WHEN** caller calls `flush()`
- **THEN** all currently buffered output SHALL be emitted
- **AND** state SHALL transition to FLUSHING
- **AND** a subsequent `process()` call SHALL transition back to STREAMING

#### Scenario: finish auto-flushes
- **GIVEN** an encoder in any non-FINISHED state
- **WHEN** caller calls `finish()`
- **THEN** remaining buffered bytes SHALL be flushed before the end-of-stream marker is written

#### Scenario: Reset from error
- **GIVEN** an encoder in ERROR state
- **WHEN** caller calls `reset()`
- **THEN** encoder SHALL return to READY state with no residual data

### Requirement: REQ-ARCH-010 Buffer-Layer Convenience API

All implementations SHALL expose a stateless buffer-layer API wrapping the streaming interface.

#### Scenario: Buffer encode
- **GIVEN** a byte slice as input
- **WHEN** caller invokes `encode_buffer(algo, input)`
- **THEN** the function SHALL return the complete encoded output as a new byte slice
- **AND** the call SHALL be equivalent to `new encoder → process(input) → finish()`

#### Scenario: Buffer decode
- **GIVEN** a valid encoded byte slice
- **WHEN** caller invokes `decode_buffer(algo, input)`
- **THEN** function SHALL return fully decoded output
- **AND** an invalid or truncated input SHALL return `ERR_TRUNCATED` or `ERR_CORRUPT`

### Requirement: REQ-ARCH-011 Output Buffer Contract

Implementations SHALL define and document maximum output expansion ratios.

#### Scenario: BUF_TOO_SMALL is transactional
- **GIVEN** an output buffer smaller than required
- **WHEN** `process()`, `flush()`, or `finish()` is called
- **THEN** the call SHALL return `BUF_TOO_SMALL`
- **AND** internal encoder state SHALL be unchanged (caller may retry with larger buffer)

#### Scenario: Security limits enforced at streaming boundary — output
- **GIVEN** a streaming decoder accumulating output
- **WHEN** cumulative decoded output would exceed 1 GiB
- **THEN** decoder SHALL return `ERR_SIZE_LIMIT` and enter ERROR state

#### Scenario: Security limits enforced at streaming boundary — input
- **GIVEN** a streaming encoder or decoder receiving input chunks
- **WHEN** cumulative input bytes would exceed 4 GiB
- **THEN** the implementation SHALL return `ERR_SIZE_LIMIT` and enter ERROR state

