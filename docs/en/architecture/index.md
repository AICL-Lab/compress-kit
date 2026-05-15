---
title: System Architecture
description: CompressKit layered architecture and module design
---

# System Architecture Design

CompressKit employs a clear layered architecture ensuring code maintainability, testability, and cross-language consistency.

## Architecture Overview

```mermaid
graph TB
    subgraph "CLI Layer"
        CLI[CLI Entry Point]
    end
    
    subgraph "Buffer Layer"
        BE[BufferedEncoder]
        BD[BufferedDecoder]
    end
    
    subgraph "Streaming Layer"
        SE[StreamingEncoder]
        SD[StreamingDecoder]
        FSM[5-State FSM]
    end
    
    subgraph "Algorithm Core"
        H[Huffman]
        A[Arithmetic]
        R[Range]
        RLE[RLE]
    end
    
    subgraph "Shared Utilities"
        Codec[codec.Encoder/Decoder]
        Err[errors]
        Bits[bits]
        Freq[frequency]
        Buf[buffer]
    end
    
    CLI --> BE
    CLI --> BD
    BE --> SE
    BD --> SD
    SE --> FSM
    SD --> FSM
    FSM --> H
    FSM --> A
    FSM --> R
    FSM --> RLE
    H --> Codec
    A --> Codec
    R --> Codec
    RLE --> Codec
    Codec --> Err
    Codec --> Bits
    Codec --> Freq
    Codec --> Buf
```

## Layer Descriptions

### 1. CLI Layer

Unified command-line entry supporting all algorithms and languages:

```bash
compress-kit encode --algo huffman --lang go input.bin output.bin
compress-kit decode --algo huffman --lang rust output.bin decoded.bin
```

**Design highlight**: 94% boilerplate reduction through unified launcher.

### 2. Buffer Layer

Stateless convenience wrapper for simple use cases:

```go
// Go example
encoder := huffman.NewBufferedEncoder()
output, err := encoder.Encode(input)
```

```rust
// Rust example
let encoder = huffman::BufferedEncoder::new();
let output = encoder.encode(&input)?;
```

**Features**:
- Each call is independent
- Automatic buffer management
- Simplified error handling

### 3. Streaming Layer

Core state machine implementation supporting incremental processing:

```go
encoder := huffman.NewStreamingEncoder()

// Incremental processing
encoder.Process(chunk1)
encoder.Process(chunk2)
encoder.Process(chunk3)

// Finish and get result
output, err := encoder.Finish()
```

**Features**:
- 5-state finite state machine
- Transactional error handling
- Flush and reset support

### 4. Algorithm Core

Implementations of four compression algorithms:

| Algorithm | File | Core Functions |
|-----------|------|----------------|
| Huffman | `huffman/encode.go` | `encodeBlock()`, `buildTree()` |
| Arithmetic | `arithmetic/encode.go` | `encodeSymbol()`, `normalize()` |
| Range | `range/encode.go` | `encodeSymbol()`, `shiftBytes()` |
| RLE | `rle/encode.go` | `encodeRun()` |

### 5. Shared Utilities

Cross-algorithm shared infrastructure:

| Module | Function |
|--------|----------|
| `codec` | Encoder/Decoder interface definitions |
| `errors` | Unified error types and codes |
| `bits` | Bit writer/reader |
| `frequency` | Frequency table processing |
| `buffer` | Buffer management |

## Binary Format Specification

### Common Structure

```
| Magic (4 bytes) | Header | Payload |
```

### Frequency Table Format

**Cross-language unified**:

- Order: symbols 0-255 (byte values), symbol 256 (EOF)
- Byte order: Little-Endian
- Total size: 4 bytes (symbol count) + 257 × 4 bytes = 1032 bytes

## Security Boundaries

| Limit | Value | Purpose |
|-------|-------|---------|
| Max input size | 4 GiB | Prevent frequency overflow and decompression bomb attacks |
| Max output size (decode only) | 1 GiB | Prevent decompression bomb attacks |

## Deep Module Design

CompressKit follows the Deep Module principle:

```
Deep Module = Simple interface + Complex implementation

BufferedEncoder.Encode(input) → output
    ↓
Hidden complexity:
- State machine management
- Buffer expansion
- Error handling
- Bit alignment
```

## Further Reading

- [Streaming API](/en/api/streaming) - 5-state FSM details
- [Streaming API](/en/api/streaming) - Complete API documentation
- [Cross-Language Testing](/en/testing/cross-language) - Conformance verification
