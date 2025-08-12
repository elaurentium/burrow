use std::path::PathBuf;

use clap::Parser;

#[derive(Parser)]
#[command(author, version, about, long_about = None)]
pub struct Args {
    #[arg(required = true)]
    pub paths: Vec<PathBuf>,
    #[arg(short, long, default_value = "zsh")]
    pub shell: String,
    #[arg(short, long)]
    pub flags: bool
}
