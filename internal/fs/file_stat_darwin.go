//go:build darwin
// +build darwin

package fs

import (
    "time"
    "syscall"
)

func timesFromStat(st syscall.Stat_t) (time.Time, time.Time, time.Time) {
    return time.Unix(int64(st.Atimespec.Sec), int64(st.Atimespec.Nsec)),
           time.Unix(int64(st.Mtimespec.Sec), int64(st.Mtimespec.Nsec)),
           time.Unix(int64(st.Ctimespec.Sec), int64(st.Ctimespec.Nsec))
}