package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/diamondburned/l4d2lb/pages"
	"github.com/diamondburned/l4d2lb/pages/index"
	"github.com/diamondburned/l4d2lb/stats"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("env"); err != nil {
		log.Fatalln("Failed to load `env':", err)
	}

	s, err := stats.Connect(os.Getenv("MYSQL_ADDRESS"))
	if err != nil {
		log.Fatalln("Failed to connect:", err)
	}

	var siteName string
	if v := os.Getenv("L4D2LB_SVNAME"); v != "" {
		siteName = v
	} else {
		siteName = "Left 4 Dead 2"
	}

	r := chi.NewMux()
	r.Mount("/static", http.StripPrefix("/static", pages.MountStatic()))
	r.Mount("/", index.Mount(&pages.RenderState{
		Database: s,
		SiteName: siteName,
	}))

	var httpAddr = os.Getenv("HTTP_FADDRESS")
	log.Println("Starting up at", httpAddr)

	if strings.HasPrefix(httpAddr, "unix://") {
		httpAddr = strings.TrimPrefix(httpAddr, "unix://")

		// Ensure that the socket is cleaned up because we're not gracefully
		// handling closes.
		if err := os.Remove(httpAddr); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				log.Fatalln("Failed to clean up old socket:", err)
			}
		}

		l, err := net.Listen("unix", httpAddr)
		if err != nil {
			log.Fatalln("Failed to listen to Unix socket:", err)
		}
		defer l.Close()

		if err := http.Serve(l, r); err != nil {
			log.Fatalln("Failed to serve HTTP:", err)
		}
	} else {
		if err := http.ListenAndServe(httpAddr, r); err != nil {
			log.Fatalln("Failed to listen and serve HTTP:", err)
		}
	}
}
