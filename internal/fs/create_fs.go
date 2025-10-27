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

package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

type Args struct {
	Paths []string
	Shell string
	Flags bool
}

type CreateFs interface {
	Create(paths []string) error
}

func Create(paths []string) error {
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			fmt.Fprintf(os.Stdout, "Path already exists: %s\n", path)
			continue
		}

		parent := filepath.Dir(path)
		if parent != "." && parent != "" {
			if _, err := os.Stat(parent); os.IsNotExist(err) {
				if err := os.MkdirAll(parent, 0755); err != nil {
					fmt.Fprintf(os.Stderr, "Error creating directory %s: %v\n", parent, err)
					continue
				}
			}
		}

		if filepath.Ext(path) != "" {
			if _, err := os.Create(path); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating file %s: %v\n", path, err)
			} else {
				fmt.Fprintln(os.Stdout)
			}
		} else {
			if err := os.MkdirAll(path, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory %s: %v\n", path, err)
			} else {
				fmt.Fprintln(os.Stdout)
			}
		}
	}

	return nil
}
