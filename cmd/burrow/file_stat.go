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

package burrow

import (
	"fmt"
	"os"

	create "github.com/elaurentium/burrow/internal/fs"
	"github.com/spf13/cobra"
)

type fileStatOptions struct {
	*ProjectOptions
	All bool // show all files
}

func (f *fileStatOptions) buildFileStat() *cobra.Command {
	opts := &fileStatOptions{
		ProjectOptions: &ProjectOptions{},
		All:            f.All,
	}
	cmd := &cobra.Command{
		Use:   "stat [OPTIONS] [FILE...]",
		Short: "File Stat CLI Tool",
		Long:  "File Stat CLI Tool. Gets information about a file.",
		RunE: func(_ *cobra.Command, args []string) error {
			if opts.All {
				entries, err := os.ReadDir(".")
				if err != nil {
					fmt.Printf("error to read dir: %s", err)
				}
				var paths []string
				for _, entry := range entries {
					paths = append(paths, entry.Name())
				}
				return create.StatInfo(paths)
			}
			return create.StatInfo(args)
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&opts.All, "all", false, "show all files")

	return cmd
}

func statCommand(p *ProjectOptions) *cobra.Command {
	opts := &fileStatOptions{
		ProjectOptions: p,
	}

	return opts.buildFileStat()
}
