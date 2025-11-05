package main

import (
	"database/sql"
	"fmt"
	"net"
	"os"

	_ "github.com/glebarez/go-sqlite"
	"proyecto-final.com/internal/query"
)

func HandleConn(net *net.Conn) {
	fmt.Print("hello")
}

func main() {
	// Connectando a la base de datos sqlite previamente creada
	c, err := sql.Open("sqlite", "./database/database.db")
	if err != nil {
		fmt.Printf("An error occur \n%s\n", err)
		os.Exit(1)
	}

	result, err := query.ManyQuery(c)
	if err != nil {
		fmt.Printf("An error occur \n%s\n", err)
		os.Exit(1)
	}

	for _, f := range result {
		fmt.Printf("%d \t %s \t %s \n", query.GetUserId(&f), query.GetUserName(&f), query.GetUserUsername(&f))
	}

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("An error occur \n%s\n", err)
		os.Exit(1)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("An error occur \n%s\n", err)
			os.Exit(1)
		}
		defer c.Close()
		go HandleConn(&conn)
	}
}
