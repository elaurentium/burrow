mod cmd;
mod init;
mod shell;

use std::process::ExitCode;

use clap::Parser;
use cmd::cli::Args;
use cmd::mkdir;
use init::init;
use cmd::cmd::InitHook;

fn main() -> ExitCode {
    let args = Args::parse();

    if let Err(e) = mkdir::run(args.paths) {
        eprintln!("Error: {}", e);
        return ExitCode::FAILURE
    }

    let opts = shell::Opts {
        cmd: args.shell,
        hook: InitHook::None,
        echo: args.flags,
    };

    match init(&opts.cmd, &opts) {
        Ok(output) => {
            println!("{}", output);
            ExitCode::SUCCESS
        }
        Err(e) => {
            eprintln!("Init error: {}", e);
            ExitCode::FAILURE
        }
    }
}
