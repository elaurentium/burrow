use anyhow::{Result, anyhow};
use askama::Template;

use crate::shell::{Bash, Opts, Zsh};

pub fn init(shell: &str, opts: &Opts) -> Result<String> {
    match shell {
        "bash" => Ok(Bash(opts).render()?),
        "zsh" => Ok(Zsh(opts).render()?),
        _ => Err(anyhow!("Unsupported shell: {}", shell)),
    }
}
