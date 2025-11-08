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

package command

import (
	"context"
	"os"
	"sync"

	"github.com/elaurentium/burrow/cmd/command/streams"
	"github.com/moby/moby/client"
)

// Streams is an interface which exposes the standard input and output streams
type Streams interface {
	In() *streams.In
	Out() *streams.Out
	Err() *streams.Out
}

type Cli interface {
	Streams
	SetIn(in *streams.In)
	CurrentVersion() string
}

type BurrowCli struct {
	in     *streams.In
	out    *streams.Out
	err    *streams.Out
	client client.APIClient
	init   sync.Once
}

func NewCli() *BurrowCli {
	return &BurrowCli{
		in:   streams.NewIn(os.Stdin),
		out:  streams.NewOut(os.Stdout),
		err:  streams.NewOut(os.Stderr),
		init: sync.Once{},
	}
}

// CurrentVersion returns the API version currently negotiated, or the default
// version otherwise.
func (cli *BurrowCli) CurrentVersion() string {
	_ = cli.initialize()
	if cli.client == nil {
		return client.MaxAPIVersion
	}
	return cli.client.ClientVersion()
}

func (cli *BurrowCli) SetIn(in *streams.In) {
	cli.in = in
}

func (cli *BurrowCli) Out() *streams.Out {
	return cli.out
}

func (cli *BurrowCli) Err() *streams.Out {
	return cli.err
}

func (cli *BurrowCli) In() *streams.In {
	return cli.in
}

// initialize initializes the Burrow API connection.
// This method uses sync.Once to ensure it only runs once.
func (cli *BurrowCli) initialize() error {
	var err error
	cli.init.Do(func() {
		// Create a new Burrow with default options
		cli.client, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			// Handle error appropriately
			return
		}

		// Ping the server to ensure the client is working
		_, err = cli.client.Ping(context.Background(), client.PingOptions{})
		if err != nil {
			// If ping fails, set client to nil
			_ = cli.client.Close()
			cli.client = nil
			return
		}
	})

	return err
}
