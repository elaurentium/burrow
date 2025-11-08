//go:build linux
// +build linux

package fs

import (
	"syscall"
	"time"
)

func timesFromStat(st syscall.Stat_t) (time.Time, time.Time, time.Time) {
	return time.Unix(int64(st.Atim.Sec), int64(st.Atim.Nsec)),
		time.Unix(int64(st.Mtim.Sec), int64(st.Mtim.Nsec)),
		time.Unix(int64(st.Ctim.Sec), int64(st.Ctim.Nsec))
}
