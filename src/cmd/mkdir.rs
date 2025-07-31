#![allow(non_snake_case, non_upper_case_globals)]
use std::error::Error;
use std::fs;
use std::io::{self, Write};
use std::path::{Path, PathBuf};


pub fn run(paths: Vec<PathBuf>) -> Result<(), Box<dyn Error>> {

    for path in paths {

        if path.exists() {
            writeln!(io::stdout(), "Path already exists: {}", path.display())?;
            continue;
        }

        if let Some(parent) = path.parent() {
            if !parent.exists() {
                if let Err(e) = fs::create_dir_all(parent) {
                    writeln!(
                        io::stderr(),
                        "Error creating directory {}: {}",
                        parent.display(),
                        e
                    )
                    .unwrap();
                    continue;
                }
            }
        }

        if path.extension().is_some() {
            match fs::File::create(&path) {
                Ok(_) => writeln!(io::stdout())?,
                Err(e) => writeln!(
                    io::stderr(),
                    "Error creating file {}: {}",
                    path.display(),
                    e
                )?,
            }
        } else {
            match fs::create_dir(&path) {
                Ok(_) => writeln!(io::stdout())?,
                Err(e) => writeln!(
                    io::stderr(),
                    "Error creating directory {}: {}",
                    path.display(),
                    e
                )?,
            }
        }
    }
    Ok(())
}
