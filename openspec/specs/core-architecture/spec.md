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
