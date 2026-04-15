use rangecoder::{decode, encode};
use std::env;
use std::fs;
use std::process;

// Range coder CLI 封装。
// Read entire file into memory，调用 rangecoder 库执行编解码，写出结果。
// File format is fully compatible with C++/Go implementations，supports cross-language encode/decode verification。

fn main() {
    let args: Vec<String> = env::args().collect();
    if args.len() != 4 {
        eprintln!("usage: {} encode|decode input output", args[0]);
        process::exit(1);
    }
    let mode = &args[1];
    let input_path = &args[2];
    let output_path = &args[3];

    let result = match mode.as_str() {
        "encode" => run_encode(input_path, output_path),
        "decode" => run_decode(input_path, output_path),
        _ => {
            eprintln!("unknown mode, expected encode or decode");
            process::exit(1);
        }
    };

    if let Err(e) = result {
        eprintln!("execution failed: {e}");
        process::exit(1);
    }
}

fn run_encode(input_path: &str, output_path: &str) -> Result<(), Box<dyn std::error::Error>> {
    let data = fs::read(input_path)?;
    let encoded = encode(&data)?;
    fs::write(output_path, &encoded)?;
    Ok(())
}

fn run_decode(input_path: &str, output_path: &str) -> Result<(), Box<dyn std::error::Error>> {
    let data = fs::read(input_path)?;
    let decoded = decode(&data)?;
    fs::write(output_path, &decoded)?;
    Ok(())
}
