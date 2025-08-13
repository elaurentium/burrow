# Directory/File Creation CLI Tool Design

## Project Overview
Create a CLI tool in Rust that can create directories and files based on command-line arguments.

## Requirements
- Accept directory and file paths as command-line arguments
- Create directory structures like "src/bin/"
- Create files like "main.rs" in specified paths
- Handle errors gracefully

## Dependencies
- `clap` for command-line argument parsing

## CLI Interface Design
```
b <path1> <path2> ... <pathN>
```

Examples:
```
b src/bin/main.rs
b src/lib.rs tests/
b src/ tests/ README.md
```

## Instalation
```
##Cargo instalation
cargo install burrow-io
```
