package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/diamondburned/l4d2lb/pages/leaderboard"
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

	r := chi.NewMux()
	r.Mount("/", leaderboard.Mount(s))

	var httpAddr = os.Getenv("HTTP_FADDRESS")
	log.Println("Starting up at", httpAddr)

	if strings.HasPrefix(httpAddr, "unix://") {
		httpAddr = strings.TrimPrefix(httpAddr, "unix://")

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
