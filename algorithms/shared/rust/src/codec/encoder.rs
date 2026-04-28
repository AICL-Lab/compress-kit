use crate::codec::error::{CodecError, State};

/// Encoder defines the streaming encoder interface.
/// All CompressKit algorithms implement this trait.
///
/// # Lifecycle
/// READY → STREAMING → FLUSHING → FINISHED
/// Error transitions to ERROR state from any state.
/// Reset() returns to READY from any state.
///
/// # Thread Safety
/// Encoders are NOT thread-safe. Callers must manage concurrency.
pub trait Encoder {
    /// Process encodes input bytes and writes to the output buffer.
    /// Returns the number of bytes written to out.
    /// If out is too small, returns BufTooSmall and state is unchanged (transactional).
    ///
    /// # State Transitions
    /// - READY → STREAMING
    /// - STREAMING → STREAMING
    /// - FLUSHING → STREAMING
    fn process(&mut self, input: &[u8], output: &mut [u8]) -> Result<usize, CodecError>;

    /// Flush writes all buffered output to the output buffer.
    /// Returns the number of bytes written.
    /// If out is too small, returns BufTooSmall and state is unchanged.
    ///
    /// # State Transitions
    /// - READY → READY (no-op)
    /// - STREAMING → FLUSHING
    /// - FLUSHING → FLUSHING (idempotent)
    fn flush(&mut self, output: &mut [u8]) -> Result<usize, CodecError>;

    /// Finish flushes any remaining buffered data and writes the end-of-stream marker.
    /// Returns the number of bytes written.
    /// If out is too small, returns BufTooSmall and state is unchanged.
    ///
    /// # State Transitions
    /// - READY/STREAMING/FLUSHING → FINISHED
    /// - FINISHED → ERROR (invalid)
    fn finish(&mut self, output: &mut [u8]) -> Result<usize, CodecError>;

    /// Reset returns the encoder to READY state, clearing all internal buffers.
    /// Can be called from any state.
    fn reset(&mut self);

    /// State returns the current lifecycle state.
    fn state(&self) -> State;
}

/// Decoder defines the streaming decoder interface.
pub trait Decoder {
    /// Process decodes input bytes and writes to the output buffer.
    /// Returns the number of bytes written to out.
    /// If out is too small, returns BufTooSmall and state is unchanged.
    ///
    /// # State Transitions
    /// Same as Encoder::process
    fn process(&mut self, input: &[u8], output: &mut [u8]) -> Result<usize, CodecError>;

    /// Flush writes all buffered output to the output buffer.
    ///
    /// # State Transitions
    /// Same as Encoder::flush
    fn flush(&mut self, output: &mut [u8]) -> Result<usize, CodecError>;

    /// Finish completes decoding and verifies the end-of-stream marker.
    /// Returns ErrTruncated if input is incomplete.
    ///
    /// # State Transitions
    /// Same as Encoder::finish
    fn finish(&mut self, output: &mut [u8]) -> Result<usize, CodecError>;

    /// Reset returns the decoder to READY state.
    fn reset(&mut self);

    /// State returns the current lifecycle state.
    fn state(&self) -> State;
}
