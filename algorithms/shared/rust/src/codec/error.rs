use std::error::Error;
use std::fmt;
use std::io;

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
pub const MAX_OUTPUT_SIZE: usize = 1024 * 1024 * 1024; // 1 GiB

/// Converts an io::Error to a CodecError.
///
/// This is a standard conversion used by algorithms that internally use io::Error
/// but need to return CodecError for the streaming API.
///
/// # Example
///
/// ```ignore
/// use compresskit_codec::codec::io_error_to_codec_error;
///
/// let io_err = io::Error::new(io::ErrorKind::UnexpectedEof, "truncated");
/// let codec_err = io_error_to_codec_error(io_err);
/// assert_eq!(codec_err, CodecError::Truncated);
/// ```
pub fn io_error_to_codec_error(e: io::Error) -> CodecError {
    match e.kind() {
        io::ErrorKind::UnexpectedEof => CodecError::Truncated,
        io::ErrorKind::InvalidData => {
            let msg = e.to_string();
            if msg.contains("truncated") || msg.contains("too short") {
                CodecError::Truncated
            } else if msg.contains("invalid") || msg.contains("bad") {
                CodecError::Corrupt
            } else {
                CodecError::Other(msg)
            }
        }
        _ => CodecError::Other(e.to_string()),
    }
}

/// Converts an io::Error to a CodecError using a wrapper function.
///
/// This is useful for mapping errors in Result types.
pub fn map_io_error<T>(result: Result<T, io::Error>) -> Result<T, CodecError> {
    result.map_err(io_error_to_codec_error)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn io_error_truncated() {
        let err = io::Error::new(io::ErrorKind::UnexpectedEof, "eof");
        assert_eq!(io_error_to_codec_error(err), CodecError::Truncated);
    }

    #[test]
    fn io_error_invalid_data_truncated() {
        let err = io::Error::new(io::ErrorKind::InvalidData, "truncated data");
        assert_eq!(io_error_to_codec_error(err), CodecError::Truncated);
    }

    #[test]
    fn io_error_invalid_data_corrupt() {
        let err = io::Error::new(io::ErrorKind::InvalidData, "invalid magic");
        assert_eq!(io_error_to_codec_error(err), CodecError::Corrupt);
    }

    #[test]
    fn io_error_invalid_data_other() {
        let err = io::Error::new(io::ErrorKind::InvalidData, "unknown issue");
        assert!(matches!(io_error_to_codec_error(err), CodecError::Other(_)));
    }

    #[test]
    fn io_error_other() {
        let err = io::Error::new(io::ErrorKind::PermissionDenied, "access denied");
        assert!(matches!(io_error_to_codec_error(err), CodecError::Other(_)));
    }
}
