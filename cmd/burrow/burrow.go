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
	"github.com/elaurentium/burrow/cmd/command"
	create "github.com/elaurentium/burrow/internal/fs"

	"github.com/elaurentium/burrow/internal/helper"
	"github.com/spf13/cobra"
)

type ProjectOptions struct {
	All bool
}

func RootCmd(cli command.Cli) *cobra.Command {
	opts := &ProjectOptions{}
	var (
		version bool
	)
	c := &cobra.Command{
		Use:   helper.Usage,
		Short: "Directory/File Creation CLI Tool",
		Long:  "Create directories and files quickly. Paths with extensions are treated as files; others as directories.",
		Args:  cobra.ArbitraryArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			if version {
				versionCommand(cli)
				return nil
			}
			// Constructor of creator
			creator := create.NewCreator()
			creator.Create(args)
			return nil
		},
	}

	c.AddCommand(
		updateCommand(),
		statCommand(opts),
		versionCommand(cli),
	)

	return c
}

func Execute() error {
	cli := command.NewCli()
	root := RootCmd(cli)
	return root.Execute()
}
