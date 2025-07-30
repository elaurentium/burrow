mod cmd;

use clap::Parser;
use cmd::cli::Args;
use cmd::operations::processPaths;

fn main() {
    let args = Args::parse();
    processPaths(args.paths);
}