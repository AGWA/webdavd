// Copyright (C) 2021-2023 Andrew Ayer
//
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
// THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.
//
// Except as contained in this notice, the name(s) of the above copyright
// holders shall not be used in advertising or otherwise to promote the
// sale, use or other dealings in this Software without prior written
// authorization.

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/webdav"

	"src.agwa.name/go-listener"
	_ "src.agwa.name/go-listener/tls"
)

func usageError(message string) {
	fmt.Fprintln(os.Stderr, message)
	flag.Usage()
	os.Exit(2)
}

func isLogNoise(message string) bool {
	return strings.HasPrefix(message, "http: TLS handshake error")
}

type httpServerLogWriter struct{}

func (httpServerLogWriter) Write(p []byte) (int, error) {
	if message := string(p); !isLogNoise(message) {
		log.Print(message)
	}
	return len(p), nil
}

func main() {
	var flags struct {
		readwrite bool
		root      string
		users     string
		public    bool
		listen    []string
	}
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "For go-listener syntax, see https://pkg.go.dev/src.agwa.name/go-listener#readme-listener-syntax\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Each line of the users file should contain a username and password separated by whitespace\n")
	}
	flag.BoolVar(&flags.readwrite, "readwrite", false, "Allow read/write access (read-only is the default)")
	flag.StringVar(&flags.root, "root", "", "Path to root directory (required)")
	flag.StringVar(&flags.users, "users", "", "Path to users file (required unless -public is used)")
	flag.BoolVar(&flags.public, "public", false, "Don't require authentication")
	flag.Func("listen", "Socket to listen on, in go-listener syntax (repeatable)", func(arg string) error {
		flags.listen = append(flags.listen, arg)
		return nil
	})
	flag.Parse()

	if flags.root == "" {
		usageError("-root flag required")
	}
	if len(flags.listen) == 0 {
		usageError("At least one -listen flag required")
	}
	if flags.users == "" && !flags.public {
		usageError("Either -users or -public must be specified")
	}
	if flags.users != "" && flags.public {
		usageError("-users and -public can't both be specified")
	}

	handler := &webdav.Handler{
		LockSystem: webdav.NewMemLS(),
	}

	if flags.readwrite {
		handler.FileSystem = webdav.Dir(flags.root)
	} else {
		handler.FileSystem = readOnlyFileSystem{webdav.Dir(flags.root)}
	}

	httpServer := http.Server{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
		ErrorLog:     log.New(httpServerLogWriter{}, "", 0),
	}

	if flags.public {
		httpServer.Handler = handler
	} else {
		users, err := loadUsersFile(flags.users)
		if err != nil {
			log.Fatal(err)
		}
		httpServer.Handler = authHandler(users, handler)
	}

	listeners, err := listener.OpenAll(flags.listen)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.CloseAll(listeners)

	for _, l := range listeners {
		go func(l net.Listener) {
			log.Fatal(httpServer.Serve(l))
		}(l)
	}

	select {}
}
