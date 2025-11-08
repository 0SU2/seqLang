package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	_ "github.com/glebarez/go-sqlite"
	"proyecto-final.com/internal/chatter"
	"proyecto-final.com/internal/query"
)

const PORT = "8080"

func main() {
	// Connectando a la base de datos sqlite previamente creada
	db, err := sql.Open("sqlite", "./database/database.db")
	if err != nil {
		log.Fatalf("An error occur trying to open sql: %s\n", err)
	}

	result, err := query.ManyQuery(db)
	if err != nil {
		log.Fatalf("An error while trying query users: %s\n", err)
	}

	for _, f := range result {
		log.Printf("%d \t %s \t %s \n", query.GetUserId(&f), query.GetUserName(&f), query.GetUserUsername(&f))
	}

	net, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("An error occur starting the server: %s\n", err)
	}
	defer db.Close()

	// Inicializacion del servidor
	log.Printf("Listening on %s in port %s ...\n", net.Addr().Network(), PORT)

	messages := make(chan chatter.Message)
	go chatter.Server(messages)

	for {
		conn, err := net.Accept()
		if err != nil {
			log.Printf("An error occur \n%s\n", err)
			os.Exit(1)
		}

		messages <- chatter.Message{
			Type: chatter.ClientConnected,
			Conn: conn,
		}

		log.Printf("Accepted connection from %s", conn.RemoteAddr().String())

		go chatter.Client(conn, messages)
	}
}
