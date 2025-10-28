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

	"golang.org/x/sys/unix"
)

type FileStat struct {
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

func GetFileStat(path string) (FileStat, error) {
	var st unix.Stat_t
	if err := unix.Lstat(path, &st); err != nil {
		return FileStat{Path: path}, err
	}
	return FileStat{
		Dev:     uint64(st.Dev),
		Ino:     st.Ino,
		Mode:    uint32(st.Mode),
		Nlink:   uint64(st.Nlink),
		Uid:     st.Uid,
		Gid:     st.Gid,
		Rdev:    uint64(st.Rdev),
		Size:    st.Size,
		Blksize: int64(st.Blksize),
		Blocks:  st.Blocks,
		Atime:   time.Unix(int64(st.Atim.Sec), int64(st.Atim.Nsec)),
		Mtime:   time.Unix(int64(st.Mtim.Sec), int64(st.Mtim.Nsec)),
		Ctime:   time.Unix(int64(st.Ctim.Sec), int64(st.Ctim.Nsec)),
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
