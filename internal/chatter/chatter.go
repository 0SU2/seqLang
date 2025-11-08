package chatter

import (
	"log"
	"net"
)

type MessageType int

const (
	ClientConnected MessageType = iota + 1
	NewMessage
	DeleteClient
)

type Message struct {
	Type MessageType
	Conn net.Conn
	Text string
}

// funcion para manejar los mensajes de los clientes
func Server(messages chan Message) {
	conns := map[string]net.Conn{}
	for {
		msg := <-messages
		switch msg.Type {
		case ClientConnected:
			// nuevo usuario, agregarlo a la lista de clientes
			log.Printf("Client %s connected", msg.Conn.RemoteAddr().String())
			conns[msg.Conn.RemoteAddr().String()] = msg.Conn
		case DeleteClient:
			log.Printf("Client %s disconnected", msg.Conn.RemoteAddr().String())
			delete(conns, msg.Conn.RemoteAddr().String())
			msg.Conn.Close()
		case NewMessage:
			// usuario manda mensaje
			for _, conn := range conns {
				if conn.RemoteAddr().String() != msg.Conn.RemoteAddr().String() {
					_, err := conn.Write([]byte(msg.Conn.LocalAddr().String() + ` ` + msg.Text))
					if err != nil {
						log.Printf("Error mientras se intenta leer el mensaje del cliente.\n [ERROR]: %s", err)
					}
				}
			}
		}
	}

}

func Client(conn net.Conn, messages chan Message) {
	// Empezar la coneccion para el cliente
	// Buffer para la entrada de la connexion. Este buffer puede estar limitado a 512
	buff := make([]byte, 512)
	for {
		// Se debe empezar un slicer para el buffer porque puede ser mayor al tamaño que especificamos
		r, err := conn.Read(buff[:])
		if err != nil {
			conn.Close()
			messages <- Message{
				Type: DeleteClient,
				Text: string(buff[0:r]),
				Conn: conn,
			}
			return
		}
		// Tenemos el tamaño del mensaje y el buffer. Realizar una conversión del tamaño que recivimos del cliente
		// a string y mandarlo al canal
		messages <- Message{
			Type: NewMessage,
			Text: string(buff[0:r]),
			Conn: conn,
		}

		log.Printf("user send message")
	}

}
