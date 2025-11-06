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
	"syscall"
	"time"
)

// POSIX mask file types(compatibility with Mode)
const (
	S_IFMT   = 0o170000
	S_IFIFO  = 0o010000
	S_IFCHR  = 0o020000
	S_IFDIR  = 0o040000
	S_IFBLK  = 0o060000
	S_IFREG  = 0o100000
	S_IFLNK  = 0o120000
	S_IFSOCK = 0o140000
)

// FileType is textual type to consume superior layers
type FileType string

type Stat struct {
	Dev     uint64    // device ID
	Ino     uint64    // inode number
	Mode    uint32    // include type + permissions (bits S_IF* + 0777)
	Nlink   uint64    // number of hard links
	Uid     uint32    // user ID
	Gid     uint32    // group ID
	Rdev    uint64    // device ID (for special file)
	Size    int64     // file size in bytes
	Blksize int64     // preferrd block size for file system io
	Blocks  int64     // number of blocks allocated
	Atime   time.Time // last access time
	Mtime   time.Time // last modification time
	Ctime   time.Time // last status change time
	Path    string    // file path
}

// FileStat return all information about file
func FileStat(path string) (Stat, error) {
	var st syscall.Stat_t
	if err := syscall.Lstat(path, &st); err != nil {
		return Stat{Path: path}, err
	}

	return Stat{
		Dev:     st.Dev,
		Ino:     st.Ino,
		Mode:    st.Mode,
		Nlink:   st.Nlink,
		Uid:     st.Uid,
		Gid:     st.Gid,
		Rdev:    st.Rdev,
		Size:    st.Size,
		Blksize: st.Blksize,
		Blocks:  st.Blocks,
		Atime:   time.Unix(int64(st.Atim.Sec), int64(st.Atim.Nsec)),
		Mtime:   time.Unix(int64(st.Mtim.Sec), int64(st.Mtim.Nsec)),
		Ctime:   time.Unix(int64(st.Ctim.Sec), int64(st.Ctim.Nsec)),
		Path:    path,
	}, nil
}

func StatInfo(path string) error {
	st, err := FileStat(path)
	if err != nil {
		return err
	}

	fmt.Printf("total: %s\n"+
		"%o"+" %s"+" %s"+" %s"+" %s"+" %s"+" %s"+" %s"+" %s\n",
		st.Path, st.Mode, st.Uid, st.Gid, st.Size, st.Atime, st.Mtime, st.Ctime, st.Path,
	)

	return nil
}
