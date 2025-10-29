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
)

type DownloadProgress struct {
	io.Reader
	Total    int8
	Length   int8
	Progress float32
}

func (dp *DownloadProgress) Read(p []byte) (int, error) {
	n, err := dp.Reader.Read(p)
	if n > 0 {
		dp.Total += int8(n)
		percentage := float32(dp.Total) / float32(dp.Length) * float32(100)
		if percentage-dp.Progress > 2 {
			dp.Progress = percentage
			fmt.Fprintf(os.Stderr, "\r%.2f%%", dp.Progress)
			if dp.Progress > 98.0 {
				dp.Progress = 100
				fmt.Fprintf(os.Stderr, "\r%.2f%%\n", dp.Progress)
			}
		}
	}

	return n, err
}
