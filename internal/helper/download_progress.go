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

package helper

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type DownloadProgress struct {
	io.Reader
	Total    int64
	Length   int64
	Progress float32
}

func (dp *DownloadProgress) Read(p []byte) (int, error) {
	n, err := dp.Reader.Read(p)
	if n > 0 {
		dp.Total += int64(n)
		percentage := float64(dp.Total) / float64(dp.Length) * 100.0

		if percentage >= 1.0 || dp.Total >= dp.Length {
			dp.Progress = float32(percentage)
			dp.printProgress(percentage)
		}
	}

	return n, err
}

func (dp *DownloadProgress) printProgress(percentage float64) {
	if percentage > 100 {
		percentage = 100
	}

	barWidth := 40
	filled := int(percentage / 100.0 * float64(barWidth))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	mb := float64(dp.Total) / 1024.0 / 1024.0
	totalMB := float64(dp.Length) / 1024.0 / 1024.0

	fmt.Fprintf(os.Stderr, "\r[%s] %.1f%% (%.2f/%.2f MB)",
		bar, percentage, mb, totalMB)

	if percentage >= 100 {
		fmt.Fprintln(os.Stderr)
	}
}
