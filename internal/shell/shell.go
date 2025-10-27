/*

	MIT License

	Copyright (c) 2025 Evandro

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all
	copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	SOFTWARE.

*/

package shell

import (
	_ "embed"
	"strings"
	"text/template"
)

var (
	//go:embed../../templates/bash.txt
	bashTemplate string
	//go:embed../../templates/zsh.txt
	zshTemplate string
)

type Opts struct {
	Cmd  string
	Hook InitHook
	Echo bool
}

type ShellTemplate struct {
	*Opts
	tmpl *template.Template
}

type Bash struct {
	ShellTemplate
}

type Zsh struct {
	ShellTemplate
}

func (st *ShellTemplate) Execute() (string, error) {
	var buf strings.Builder
	if err := st.tmpl.Execute(&buf, st.Opts); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (b *Bash) Render() (string, error) { return b.Execute() }
func (z *Zsh) Render() (string, error)  { return z.Execute() }

func Default(val, fallback string) string {
	if strings.TrimSpace(val) == "" {
		return fallback
	}
	return val
}

func Replace(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

func funcMap() template.FuncMap {
	return template.FuncMap{
		"Default": Default,
		"Replace": Replace,
	}
}

func NewBash(opts *Opts) (*Bash, error) {
	tmpl, err := template.New("bash").Funcs(funcMap()).Parse(bashTemplate)
	if err != nil {
		return nil, err
	}
	return &Bash{ShellTemplate{Opts: opts, tmpl: tmpl}}, nil
}

func NewZsh(opts *Opts) (*Zsh, error) {
	tmpl, err := template.New("zsh").Funcs(funcMap()).Parse(zshTemplate)
	if err != nil {
		return nil, err
	}
	return &Zsh{ShellTemplate{Opts: opts, tmpl: tmpl}}, nil
}
