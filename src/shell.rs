#![allow(non_snake_case, non_upper_case_globals)]

use crate::cmd::cmd::InitHook;

#[derive(Debug, Eq, PartialEq)]
pub struct Opts {
    pub cmd: String,
    pub hook: InitHook,
    pub echo: bool
}

#[warn(unused_macros)]
macro_rules! make_template {
    ($name:ident, $path:expr) => {
        #[derive(::std::fmt::Debug, ::askama::Template)]
        #[template(path = $path)]
        pub struct $name<'a>(pub &'a self::Opts);

        impl<'a> ::std::ops::Deref for $name<'a> {
            type Target = self::Opts;
            fn deref(&self) -> &Self::Target {
                self.0
            }
        }
    };
}

make_template!(Bash, "bash.txt");
make_template!(Zsh, "zsh.txt");