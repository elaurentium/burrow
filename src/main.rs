mod cmd;
mod shell;

use std::process::ExitCode;

use clap::Parser;
use cmd::burrow;
use cmd::cli::Args;

fn main() -> ExitCode {
    let args = Args::parse();

    if let Err(e) = burrow::run(args.paths) {
        eprintln!("Error: {}", e);
        ExitCode::FAILURE
    } else {
        ExitCode::SUCCESS
    }
}
