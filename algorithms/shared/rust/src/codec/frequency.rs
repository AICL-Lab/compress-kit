use std::error::Error;
use std::fmt;

pub const SYMBOL_LIMIT: usize = 257;
pub const EOF_SYMBOL: u32 = (SYMBOL_LIMIT - 1) as u32;
const MAX_READ_FREQUENCY_COUNT: usize = 1024;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum FrequencyErrorKind {
    Truncated,
    Corrupt,
    Overflow,
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

    fn overflow(message: &'static str) -> Self {
        Self {
            kind: FrequencyErrorKind::Overflow,
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

pub fn build_frequencies(data: &[u8]) -> Result<Vec<u32>, FrequencyError> {
    let counts = count_frequencies(data);
    finalize_exact_frequencies(&counts, "symbol frequency exceeds u32::MAX")
}

pub fn build_scaled_frequencies(data: &[u8], max_total: u32) -> Vec<u32> {
    let counts = count_frequencies(data);
    finalize_scaled_frequencies(&counts, max_total)
}

fn count_frequencies(data: &[u8]) -> Vec<u64> {
    let mut freq = vec![0u64; SYMBOL_LIMIT];
    for &byte in data {
        freq[byte as usize] += 1;
    }
    freq
}

fn finalize_exact_frequencies(
    counts: &[u64],
    overflow_message: &'static str,
) -> Result<Vec<u32>, FrequencyError> {
    let mut freq = Vec::with_capacity(counts.len());
    for (index, &count) in counts.iter().enumerate() {
        let final_count = if index == EOF_SYMBOL as usize {
            count
                .checked_add(1)
                .ok_or_else(|| FrequencyError::overflow(overflow_message))?
        } else {
            count
        };
        if final_count > u32::MAX as u64 {
            return Err(FrequencyError::overflow(overflow_message));
        }
        freq.push(final_count as u32);
    }
    Ok(freq)
}

fn finalize_scaled_frequencies(counts: &[u64], max_total: u32) -> Vec<u32> {
    let mut freq = counts.to_vec();
    freq[EOF_SYMBOL as usize] += 1;
    scale_counts(&freq, max_total)
}

pub fn scale_frequencies(freq: &mut [u32], max_total: u32) {
    let scaled = scale_counts(
        &freq
            .iter()
            .map(|&value| value as u64)
            .collect::<Vec<u64>>(),
        max_total,
    );
    freq.copy_from_slice(&scaled);
}

fn scale_counts(counts: &[u64], max_total: u32) -> Vec<u32> {
    let total: u64 = counts.iter().sum();
    let mut scaled_freq = vec![0u32; counts.len()];
    if total == 0 {
        for value in scaled_freq.iter_mut() {
            *value = 1;
        }
        return scaled_freq;
    }
    if total <= max_total as u64 {
        for (value, &count) in scaled_freq.iter_mut().zip(counts.iter()) {
            *value = count as u32;
        }
        return scaled_freq;
    }

    let mut new_total = 0u64;
    for (value, &count) in scaled_freq.iter_mut().zip(counts.iter()) {
        if count == 0 {
            continue;
        }
        let mut scaled = (count * max_total as u64) / total;
        if scaled == 0 {
            scaled = 1;
        }
        *value = scaled as u32;
        new_total += scaled;
    }

    if new_total == 0 {
        let mut base = max_total / scaled_freq.len() as u32;
        if base == 0 {
            base = 1;
        }
        for value in scaled_freq.iter_mut() {
            *value = base;
        }
    } else {
        while new_total > max_total as u64 {
            let mut adjusted = false;
            for value in scaled_freq.iter_mut() {
                if *value > 1 {
                    *value -= 1;
                    new_total -= 1;
                    adjusted = true;
                    if new_total <= max_total as u64 {
                        break;
                    }
                }
            }
            if !adjusted {
                break;
            }
        }
    }

    scaled_freq
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
        .ok_or_else(|| FrequencyError::truncated(truncated_count_message))?
        as usize;
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

pub fn read_frequencies(
    input: &[u8],
    pos: &mut usize,
    truncated_count_message: &'static str,
    truncated_entries_message: &'static str,
    invalid_count_message: &'static str,
) -> Result<Vec<u32>, FrequencyError> {
    let count = read_u32_le(input, pos)
        .ok_or_else(|| FrequencyError::truncated(truncated_count_message))?
        as usize;
    if count == 0 || count > MAX_READ_FREQUENCY_COUNT {
        return Err(FrequencyError::corrupt(invalid_count_message));
    }

    let mut freq = Vec::with_capacity(count);
    for _ in 0..count {
        freq.push(
            read_u32_le(input, pos)
                .ok_or_else(|| FrequencyError::truncated(truncated_entries_message))?,
        );
    }
    Ok(freq)
}

fn read_u32_le(input: &[u8], pos: &mut usize) -> Option<u32> {
    if *pos + 4 > input.len() {
        return None;
    }
    let value = u32::from_le_bytes([
        input[*pos],
        input[*pos + 1],
        input[*pos + 2],
        input[*pos + 3],
    ]);
    *pos += 4;
    Some(value)
}

#[cfg(test)]
mod tests {
    use super::{
        build_cumulative, build_cumulative_strict, build_scaled_frequencies,
        finalize_exact_frequencies, finalize_scaled_frequencies, read_frequencies,
        read_frequencies_exact, FrequencyErrorKind, EOF_SYMBOL,
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
    fn finalize_exact_frequencies_rejects_symbol_count_above_u32_max() {
        let mut counts = vec![0u64; super::SYMBOL_LIMIT];
        counts[b'a' as usize] = u32::MAX as u64 + 1;

        let err = finalize_exact_frequencies(&counts, "symbol count exceeds u32::MAX").unwrap_err();

        assert_eq!(err.kind, FrequencyErrorKind::Overflow);
        assert_eq!(err.message, "symbol count exceeds u32::MAX");
    }

    #[test]
    fn finalize_scaled_frequencies_accepts_u64_counts_above_u32_max() {
        let mut counts = vec![0u64; super::SYMBOL_LIMIT];
        counts[b'a' as usize] = u32::MAX as u64 + 1;
        counts[b'b' as usize] = 4;

        let freq = finalize_scaled_frequencies(&counts, 8);
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
    fn read_frequencies_accepts_non_exact_count_within_supported_range() {
        let input = [
            2, 0, 0, 0, // count
            1, 0, 0, 0, 3, 0, 0, 0,
        ];
        let mut pos = 0;
        let freq = read_frequencies(
            &input,
            &mut pos,
            "truncated count",
            "truncated entries",
            "wrong symbol count",
        )
        .unwrap();

        assert_eq!(freq, vec![1, 3]);
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
