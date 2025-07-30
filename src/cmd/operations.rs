#![allow(non_snake_case, non_upper_case_globals)]
use std::path::Path;
use std::fs;
use std::io::{self, Write};

pub fn processPaths(paths: Vec<String>) {
    for pathStr in paths {
        let path = Path::new(&pathStr);

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

        if pathStr.ends_with("/") {
            let dirPath = Path::new(&pathStr[..pathStr.len() - 1]);
            if !dirPath.exists() {
                match fs::create_dir(dirPath) {
                    Ok(_) => {
                        writeln!(io::stdout(), "Created directory: {}", dirPath.display()).unwrap()
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
                    Ok(_) => writeln!(io::stdout(), "Created file: {}", path.display()).unwrap(),
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
}
