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

package streams

import (
	"errors"
	"io"
	"os"
	"runtime"

	"github.com/moby/term"
)

// In is an input stream to read user input. It implements [io.ReadCloser]
// with additional utilities, such as putting the terminal in raw mode.
type In struct {
	commonStream
	in io.ReadCloser
}

// Read implements the [io.Reader] interface.
func (i *In) Read(p []byte) (int, error) {
	return i.in.Read(p)
}

// Close implements the [io.Closer] interface.
func (i *In) Close() error {
	return i.in.Close()
}

// SetRawTerminal sets raw mode on the input terminal. It is a no-op if In
// is not a TTY, or if the "NORAW" environment variable is set to a non-empty
// value.
func (i *In) SetRawTerminal() (err error) {
	if !i.isTerminal || os.Getenv("NORAW") != "" {
		return nil
	}
	i.state, err = term.SetRawTerminal(i.fd)
	return err
}

// CheckTty checks if we are trying to attach to a container TTY
// from a non-TTY client input stream, and if so, returns an error.
func (i *In) CheckTty(attachStdin, ttyMode bool) error {
	// In order to attach to a container tty, input stream for the client must
	// be a tty itself: redirecting or piping the client standard input is
	// incompatible with `docker run -t`, `docker exec -t` or `docker attach`.
	if ttyMode && attachStdin && !i.isTerminal {
		const eText = "the input device is not a TTY"
		if runtime.GOOS == "windows" {
			return errors.New(eText + ".  If you are using mintty, try prefixing the command with 'winpty'")
		}
		return errors.New(eText)
	}
	return nil
}

// NewIn returns a new [In] from an [io.ReadCloser].
func NewIn(in io.ReadCloser) *In {
	i := &In{in: in}
	i.fd, i.isTerminal = term.GetFdInfo(in)
	return i
}
