mod cmd;
mod init;
mod shell;

use std::process::ExitCode;

use clap::Parser;
use cmd::cli::Args;
use cmd::mkdir;


fn main() -> ExitCode {
    let args = Args::parse();

    if let Err(e) = mkdir::run(args.paths) {
        eprintln!("Error: {}", e);
        ExitCode::FAILURE
    } else {
        ExitCode::SUCCESS
    }
}
