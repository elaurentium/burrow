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
	"os/user"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/unix"
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
	Nlink   uint16    // number of hard links
	Uid     uint32    // user ID
	Gid     uint32    // group ID
	Rdev    uint64    // device ID (for special file)
	Size    int64     // file size in bytes
	Blksize int32     //preferrd block size for file system io
	Blocks  int64     // number of blocks allocated
	Atime   time.Time // last access time
	Mtime   time.Time // last modification time
	Ctime   time.Time // last status change time
	Path    string    // file path
}

// FormatPermissions converts a file mode to a string like "drwxr-xr-x"
func formatPermissions(mode uint32) string {
	var buf strings.Builder

	// File type
	switch mode & S_IFMT {
	case S_IFDIR:
		buf.WriteByte('d')
	case S_IFREG:
		buf.WriteByte('-')
	case S_IFLNK:
		buf.WriteByte('l')
	case S_IFIFO:
		buf.WriteByte('p')
	case S_IFSOCK:
		buf.WriteByte('s')
	case S_IFCHR:
		buf.WriteByte('c')
	case S_IFBLK:
		buf.WriteByte('b')
	default:
		buf.WriteByte('?')
	}

	// Permissions: rwx for user, group, other
	perms := []string{"r", "w", "x"}
	for i := 0; i < 3; i++ {
		// user, group, other
		for j := 0; j < 3; j++ {
			if mode&(1<<(8-3*i-j)) != 0 {
				buf.WriteString(perms[j])
			} else {
				buf.WriteByte('-')
			}
		}
	}

	return buf.String()
}

// FileStat return all information about file
func FileStat(path string) (Stat, error) {
	var st unix.Stat_t
	if err := unix.Lstat(path, &st); err != nil {
		return Stat{Path: path}, err
	}

	return Stat{
		Dev:     uint64(st.Dev),
		Ino:     st.Ino,
		Mode:    uint32(st.Mode),
		Nlink:   uint16(st.Nlink),
		Uid:     st.Uid,
		Gid:     st.Gid,
		Rdev:    uint64(st.Rdev),
		Size:    st.Size,
		Blksize: int32(st.Blksize),
		Blocks:  st.Blocks,
		Atime:   time.Unix(int64(st.Atim.Sec), int64(st.Atim.Nsec)),
		Mtime:   time.Unix(int64(st.Mtim.Sec), int64(st.Mtim.Nsec)),
		Ctime:   time.Unix(int64(st.Ctim.Sec), int64(st.Ctim.Nsec)),
		Path:    path,
	}, nil
}

func StatInfo(paths []string) error {
	for _, path := range paths {
		st, err := FileStat(path)
		if err != nil {
			return err
		}

		// Get username
		u, err := user.LookupId(strconv.Itoa(int(st.Uid)))
		username := strconv.Itoa(int(st.Uid))
		if err == nil {
			username = u.Username
		}

		// Get group name
		g, err := user.LookupGroupId(strconv.Itoa(int(st.Gid)))
		groupname := strconv.Itoa(int(st.Gid))
		if err == nil {
			groupname = g.Name
		}

		// Format time
		mtime := st.Mtime
		now := time.Now()
		sixMonthsAgo := now.AddDate(0, -6, 0)
		timeStr := mtime.Format("Jan _2  2006")
		if mtime.After(sixMonthsAgo) {
			timeStr = mtime.Format("Jan _2 15:04")
		}

		fmt.Printf("%s %d %s %s %d %s %s\n",
			formatPermissions(st.Mode), st.Nlink, username, groupname, st.Size, timeStr, st.Path,
		)
	}

	return nil
}
