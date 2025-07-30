#![allow(non_snake_case, non_upper_case_globals)]

#[derive(Debug, Eq, PartialEq)]
pub struct Opts<'a> {
    pub cmd: Option<&'a str>,
    pub echo: bool
}

#[warn(unused_macros)]
macro_rules! make_template {
    ($name:ident, $path:expr) => {
        #[derive(::std::fmt::Debug, ::askama::Template)]
        #[template(path = $path)]
        pub struct $name<'a>(pub &'a self::Opts<'a>);

        impl<'a> ::std::ops::Deref for $name<'a> {
            type Target = self::Opts<'a>;
            fn deref(&self) -> &Self::Target {
                self.0
            }
        }
    };
}

//make_template!(Bash, "bash.txt");
//make_template!(Elvish, "elvish.txt");
//make_template!(Fish, "fish.txt");
//make_template!(Nushell, "nushell.txt");
//make_template!(Posix, "posix.txt");
//make_template!(Powershell, "powershell.txt");
//make_template!(Tcsh, "tcsh.txt");
//make_template!(Xonsh, "xonsh.txt");
//make_template!(Zsh, "zsh.txt");