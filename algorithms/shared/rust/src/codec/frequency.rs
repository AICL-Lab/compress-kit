use std::error::Error;
use std::fmt;

pub const SYMBOL_LIMIT: usize = 257;
pub const EOF_SYMBOL: u32 = (SYMBOL_LIMIT - 1) as u32;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum FrequencyErrorKind {
    Truncated,
    Corrupt,
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct FrequencyError {
    pub kind: FrequencyErrorKind,
    pub message: &'static str,
}

impl FrequencyError {
    fn truncated(message: &'static str) -> Self {
        Self {
            kind: FrequencyErrorKind::Truncated,
            message,
        }
    }

    fn corrupt(message: &'static str) -> Self {
        Self {
            kind: FrequencyErrorKind::Corrupt,
            message,
        }
    }
}

impl fmt::Display for FrequencyError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{}", self.message)
    }
}

impl Error for FrequencyError {}

pub fn build_frequencies(data: &[u8]) -> Vec<u32> {
    let mut freq = vec![0u32; SYMBOL_LIMIT];
    for &byte in data {
        freq[byte as usize] += 1;
    }
    freq[EOF_SYMBOL as usize] = 1;
    freq
}

pub fn build_scaled_frequencies(data: &[u8], max_total: u32) -> Vec<u32> {
    let mut freq = build_frequencies(data);
    scale_frequencies(&mut freq, max_total);
    freq
}

pub fn scale_frequencies(freq: &mut [u32], max_total: u32) {
    let total: u64 = freq.iter().map(|&value| value as u64).sum();
    if total == 0 {
        for value in freq.iter_mut() {
            *value = 1;
        }
        return;
    }
    if total <= max_total as u64 {
        return;
    }

    let mut new_total = 0u64;
    for value in freq.iter_mut() {
        if *value == 0 {
            continue;
        }
        let mut scaled = (*value as u64 * max_total as u64) / total;
        if scaled == 0 {
            scaled = 1;
        }
        *value = scaled as u32;
        new_total += scaled;
    }

    if new_total == 0 {
        let mut base = max_total / freq.len() as u32;
        if base == 0 {
            base = 1;
        }
        for value in freq.iter_mut() {
            *value = base;
        }
    }
}

pub fn build_cumulative(freq: &[u32]) -> Vec<u32> {
    let mut cumulative = vec![0u32; freq.len() + 1];
    for (index, &value) in freq.iter().enumerate() {
        cumulative[index + 1] = cumulative[index] + value;
    }
    if cumulative.last().copied().unwrap_or(0) == 0 {
        for index in 0..freq.len() {
            cumulative[index + 1] = (index + 1) as u32;
        }
    }
    cumulative
}

pub fn build_cumulative_strict(
    freq: &[u32],
    zero_table_message: &'static str,
) -> Result<Vec<u32>, FrequencyError> {
    if freq.iter().all(|&value| value == 0) {
        return Err(FrequencyError::corrupt(zero_table_message));
    }
    Ok(build_cumulative(freq))
}

pub fn write_frequencies(out: &mut Vec<u8>, freq: &[u32]) {
    out.extend_from_slice(&(freq.len() as u32).to_le_bytes());
    for &value in freq {
        out.extend_from_slice(&value.to_le_bytes());
    }
}

pub fn read_frequencies_exact(
    input: &[u8],
    pos: &mut usize,
    expected_count: usize,
    truncated_count_message: &'static str,
    truncated_entries_message: &'static str,
    invalid_count_message: &'static str,
) -> Result<Vec<u32>, FrequencyError> {
    let count = read_u32_le(input, pos)
        .ok_or_else(|| FrequencyError::truncated(truncated_count_message))? as usize;
    if count != expected_count {
        return Err(FrequencyError::corrupt(invalid_count_message));
    }

    let mut freq = vec![0u32; expected_count];
    for value in freq.iter_mut() {
        *value = read_u32_le(input, pos)
            .ok_or_else(|| FrequencyError::truncated(truncated_entries_message))?;
    }
    Ok(freq)
}

fn read_u32_le(input: &[u8], pos: &mut usize) -> Option<u32> {
    if *pos + 4 > input.len() {
        return None;
    }
    let value =
        u32::from_le_bytes([input[*pos], input[*pos + 1], input[*pos + 2], input[*pos + 3]]);
    *pos += 4;
    Some(value)
}

#[cfg(test)]
mod tests {
    use super::{
        build_cumulative, build_cumulative_strict, build_scaled_frequencies, read_frequencies_exact,
        EOF_SYMBOL, FrequencyErrorKind,
    };

    #[test]
    fn build_scaled_frequencies_clamps_when_feasible_and_preserves_eof() {
        let mut data = vec![b'a'; 12];
        data.extend_from_slice(&[b'b'; 6]);

        let freq = build_scaled_frequencies(&data, 8);
        let total: u32 = freq.iter().sum();

        assert!(total <= 8, "total = {total}, want <= 8");
        assert_eq!(freq[EOF_SYMBOL as usize], 1);
        assert!(freq[b'a' as usize] > freq[b'b' as usize]);
    }

    #[test]
    fn read_frequencies_exact_rejects_wrong_count_with_caller_supplied_error_message() {
        let input = [
            3, 0, 0, 0, // count
            1, 0, 0, 0, 2, 0, 0, 0, 3, 0, 0, 0,
        ];
        let mut pos = 0;
        let err = read_frequencies_exact(
            &input,
            &mut pos,
            4,
            "truncated table",
            "truncated entries",
            "wrong symbol count",
        )
        .unwrap_err();

        assert_eq!(err.kind, FrequencyErrorKind::Corrupt);
        assert_eq!(err.message, "wrong symbol count");
    }

    #[test]
    fn read_frequencies_exact_uses_entry_message_for_truncated_entries() {
        let input = [
            2, 0, 0, 0, // count
            1, 0, 0, 0, // first entry only
        ];
        let mut pos = 0;
        let err = read_frequencies_exact(
            &input,
            &mut pos,
            2,
            "truncated count",
            "truncated entries",
            "wrong symbol count",
        )
        .unwrap_err();

        assert_eq!(err.kind, FrequencyErrorKind::Truncated);
        assert_eq!(err.message, "truncated entries");
    }

    #[test]
    fn build_cumulative_uses_sequential_fallback_for_all_zero_table() {
        let cumulative = build_cumulative(&[0, 0, 0]);
        assert_eq!(cumulative, vec![0, 1, 2, 3]);
    }

    #[test]
    fn build_cumulative_strict_rejects_all_zero_table() {
        let err = build_cumulative_strict(&[0, 0, 0], "invalid table").unwrap_err();

        assert_eq!(err.kind, FrequencyErrorKind::Corrupt);
        assert_eq!(err.message, "invalid table");
    }
}
