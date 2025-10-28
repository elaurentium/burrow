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
	"os"
	"time"
)

type Stat_t struct {
	Dev     uint64
	Ino     uint64
	Mode    uint32 // include type + permissions (bits S_IF* + 0777)
	Nlink   uint64
	Uid     uint32
	Gid     uint32
	Rdev    uint64
	Size    int64
	Blksize int64
	Blocks  int64
	Atime   time.Time
	Mtime   time.Time
	Ctime   time.Time
	Path    string
}

type CreateFs interface {
	Create(paths []string) error
}

func GetFileStat(path string) (Stat_t, error) {
	info, err := os.Lstat(path)
	if  err != nil {
		return Stat_t{Path: path}, err
	}
	return Stat_t{
		Dev:     0,
		Ino:     0,
		Mode:    uint32(info.Mode()),
		Nlink:   0,
		Uid:     0,
		Gid:     0,
		Rdev:    0,
		Size:    info.Size(),
		Blksize: 0,
		Blocks:  0,
		Atime:   info.ModTime(),
		Mtime:   info.ModTime(),
		Ctime:   info.ModTime(),
		Path:    path,
	}, nil
}

func Create(paths []string, perm os.FileMode) error {
	for _, path := range paths {
		if err := os.MkdirAll(path, perm); err != nil {
			return err
		}
	}
	return nil
}
