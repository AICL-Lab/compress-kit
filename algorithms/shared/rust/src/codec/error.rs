use std::error::Error;
use std::fmt;

/// Lifecycle state of an encoder or decoder.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum State {
    /// Codec is ready to accept input.
    Ready,
    /// Codec is actively processing input.
    Streaming,
    /// Codec has been flushed and buffered output has been emitted.
    Flushing,
    /// Codec has finished and emitted the end-of-stream marker.
    Finished,
    /// Codec encountered an error and must be reset.
    Error,
}

/// Error codes for codec operations.
#[derive(Debug, Clone, PartialEq, Eq)]
pub enum CodecError {
    /// Output buffer is too small. Operation is transactional: state unchanged.
    BufTooSmall,
    /// Input stream ends prematurely during decode.
    Truncated,
    /// Data corruption or integrity check failed.
    Corrupt,
    /// Operation is not valid in the current lifecycle state.
    InvalidState,
    /// Input or output exceeds security limits (4 GiB in / 1 GiB decode out).
    SizeLimit,
    /// Frame version byte is not supported.
    VersionUnsupported,
    /// Unknown algorithm ID in frame header.
    UnknownAlgo,
    /// Other error with message.
    Other(String),
}

impl fmt::Display for CodecError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            CodecError::BufTooSmall => write!(f, "output buffer too small"),
            CodecError::Truncated => write!(f, "input stream truncated"),
            CodecError::Corrupt => write!(f, "data corrupted"),
            CodecError::InvalidState => write!(f, "invalid state for operation"),
            CodecError::SizeLimit => write!(f, "size limit exceeded"),
            CodecError::VersionUnsupported => write!(f, "unsupported version"),
            CodecError::UnknownAlgo => write!(f, "unknown algorithm"),
            CodecError::Other(msg) => write!(f, "{}", msg),
        }
    }
}

impl Error for CodecError {}

/// Security limits.
pub const MAX_INPUT_SIZE: usize = 4 * 1024 * 1024 * 1024; // 4 GiB
pub const MAX_OUTPUT_SIZE: usize = 1 * 1024 * 1024 * 1024; // 1 GiB
