package main

import (
"bufio"
"fmt"
"net"
)

func main() {
	for {
		conn, _ := net.Dial("tcp", "127.0.0.1:8081") // recoit connexion
		reader := bufio.NewReader(conn)
		for {
			message, err := reader.ReadString('\n') //pour lire les msg du serveur
			if err != nil { // si erreur, fermer la session tcp dédiée à ce client
				conn.Close()
				break
				}
			fmt.Print(message) // print les messages des autres clients
		}
		conn.Close()
		fmt.Print("Tapez votre pseudo ") //TCCHAT_REGISTER
		fmt.Print("1) ECRIRE MSG, 2) SE DECONNECTER")


	}




	}
