#![allow(non_snake_case, non_upper_case_globals)]
use std::error::Error;
use std::fs;
use std::io::{self, Write};
use std::path::{Path, PathBuf};


pub fn run(paths: Vec<PathBuf>) -> Result<(), Box<dyn Error>> {

    for pathStr in paths {
        let path = Path::new(&pathStr);

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

        let strBuffer = pathStr.to_string_lossy();
        if strBuffer.ends_with("/") {
            let dirPath = Path::new(&strBuffer[..strBuffer.len() - 1]);
            if !dirPath.exists() {
                match fs::create_dir(dirPath) {
                    Ok(_) => {
                        writeln!(io::stdout()).unwrap()
                    }
                    Err(e) => writeln!(
                        io::stderr(),
                        "Error creating directory {}: {}",
                        dirPath.display(),
                        e
                    )
                    .unwrap(),
                }
            } else {
                writeln!(
                    io::stdout(),
                    "Directory already exists: {}",
                    dirPath.display()
                )
                .unwrap();
            }
        } else {
            if !path.exists() {
                match fs::File::create(path) {
                    Ok(_) => writeln!(io::stdout()).unwrap(),
                    Err(e) => writeln!(
                        io::stderr(),
                        "Error creating file {}: {}",
                        path.display(),
                        e
                    )
                    .unwrap(),
                }
            } else {
                writeln!(io::stdout(), "File already exists: {}", path.display()).unwrap();
            }
        }
    }
    Ok(())
}
