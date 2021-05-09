package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/webdav"

	"src.agwa.name/go-listener"
	_ "src.agwa.name/go-listener/tls"
)

func main() {
	var flags struct {
		readwrite bool
		root      string
		users     string
		public    bool
		listen    []string
	}
	flag.BoolVar(&flags.readwrite, "readwrite", false, "Allow read/write access")
	flag.StringVar(&flags.root, "root", "", "Root directory")
	flag.StringVar(&flags.users, "users", "", "Path to users file")
	flag.BoolVar(&flags.public, "public", false, "Don't require authentication")
	flag.Func("listen", "Socket to listen on (repeatable)", func(arg string) error {
		flags.listen = append(flags.listen, arg)
		return nil
	})
	flag.Parse()

	if flags.root == "" {
		log.Fatal("-root flag required")
	}
	if len(flags.listen) == 0 {
		log.Fatal("At least one -listen flag required")
	}
	if flags.users == "" && !flags.public {
		log.Fatal("Either -users or -public must be specified")
	}
	if flags.users != "" && flags.public {
		log.Fatal("-users and -public can't both be specified")
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
