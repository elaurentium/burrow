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
	"sync"
)

type Creator struct {
	Perm os.FileMode
	wg   *sync.WaitGroup
}

// Multi path with paralelism
func (c *Creator) Create(paths []string) []error {
	jobs := make(chan string, len(paths))
	errors := make(chan error, len(paths))

	// Init workers
	for i := 0; i < len(paths); i++ {
		c.wg.Add(1)
		go func() {
			for path := range jobs {
				if err := c.create_fs(path); err != nil {
					errors <- fmt.Errorf("Failed to create %s: %w", path, err)
				}
			}
		}()
	}

	// Send jobs
	for _, path := range paths {
		jobs <- path
	}
	close(jobs)

	// Wait for workers to finish
	c.wg.Wait()
	close(errors)

	//Colect errors
	var errs []error
	for err := range errors {
		errs = append(errs, err)
	}

	return errs
}

// Handle creation fs
func (c *Creator) create_fs(path string) error {
	// Verify path if already exists
	if _, err := os.Stat(path); err == nil {
		fmt.Fprintf(os.Stdout, "Path already exists: %s\n", path)
		return nil
	}

	parent := filepath.Dir(path)
	if parent != "." && parent == "" {
		// Verify parent dir exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// Create parent dir if not exists
			if err := os.MkdirAll(parent, c.Perm); err != nil {
				return fmt.Errorf("Error creating directory %s: %w", parent, err)
			}
		}
	}

	// Create file or dir
	if filepath.Ext(path) != "" {
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("Error creating file %s: %w", path, err)
		}
		file.Close()
	} else {
		if err := os.MkdirAll(path, c.Perm); err != nil {
			return fmt.Errorf("Error creating directory %s: %w", path, err)
		}
	}

	return nil
}
