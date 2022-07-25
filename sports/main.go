package main

import (
	"database/sql"
	"flag"
	"log"
	"net"
	"sports/proto/sports"

	"google.golang.org/grpc"
	"sports/db"
	"sports/service"
)

var (
	sportEndpoint = flag.String("sport-endpoint", "localhost:9200", "sport server endpoint")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("failed running sport server: %s\n", err)
	}
}

func run() error {
	conn, err := net.Listen("tcp", ":9200")
	if err != nil {
		return err
	}

	sportDB, err := sql.Open("sqlite3", "./db/sports.db")
	if err != nil {
		return err
	}

	sportRepo := db.NewSportsRepo(sportDB)
	if err := sportRepo.Init(); err != nil {
		return err
	}

	sportServer := grpc.NewServer()

	sports.RegisterSportsServer(
		sportServer,
		service.NewSportService(
			sportRepo,
		),
	)

	log.Printf("sport server listening on: %s\n", *sportEndpoint)

	if err := sportServer.Serve(conn); err != nil {
		return err
	}

	return nil
}
