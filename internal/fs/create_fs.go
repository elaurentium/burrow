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
	"strings"
	"sync"

	"github.com/elaurentium/burrow/internal/helper"
)

type Creator struct {
	Perm    os.FileMode
	Workers int
	Wg      *sync.WaitGroup
}

func NewCreator() *Creator {
	return &Creator{
		Perm:    0755,
		Workers: 0,
		Wg:      &sync.WaitGroup{},
	}
}

func isFile(path string) bool {
	if filepath.Ext(path) != "" {
		return true
	}
	base := filepath.Base(path)
	for _, file := range helper.FilesWithoutExtension {
		if strings.EqualFold(base, file) {
			return true
		}
	}
	return false
}

func (c *Creator) Create(paths []string) error {
	for _, path := range paths {
		parent := filepath.Dir(path)
		if parent != "." && parent != "" {
			if err := os.MkdirAll(parent, c.Perm); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory %s: %v\n", parent, err)
				continue
			}
		}
		if isFile(path) {
			f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL, c.Perm)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating file %s: %v\n", path, err)
				continue
			}
			f.Close()
		} else {
			if err := os.MkdirAll(path, c.Perm); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory %s: %v\n", path, err)
			}
		}
	}
	return nil
}
