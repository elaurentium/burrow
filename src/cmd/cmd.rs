#![allow(clippy::module_inception)]

#[derive(Debug, Copy, Clone, Eq, PartialEq)]
pub enum InitHook {
    None,
    Prompt,
    Pwd
}

