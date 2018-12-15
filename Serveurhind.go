package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func handleConnection(c net.Conn) { //gérer chaque connexion de chaque clt
	for {
		message, err := bufio.NewReader(c).ReadString('\n') //pour lire le msg des clients, erreurs à gérer après
		if err != nil { // si erreur, fermer la session tcp dédiée à ce client
			c.Close()
			break
		}

		fmt.Print(message)
	}
}

func main() {

	ln, _ := net.Listen("tcp", ":8081") //le serveur écoute sur toutes les interfaces
	fmt.Println("Serveur en cours de synchronisation...")
	for {
		conn, _ := ln.Accept()              // accepte les connexions
		message := "Bienvenue sur TCCHAT, vous êtes connectés :-)"
		newmessage := strings.ToUpper(message)
		conn.Write([]byte(newmessage + "\n"))
		//cmdLine := make([]byte,(1024 * 4)) // buffer pour stocker les données en attente de la part d'autres client
		go handleConnection(conn)


		}

	}



	//fmt.Printf("%q\n", strings.Split("a,b,c", ","))

	//var names []string
	//names = strings.Split("Ta,b,c", ",")
	//fmt.Printf("%q\n", names)

	//for _, name := range names {
	//	fmt.Println(name)

	//}
	//for name := range names {
	//	fmt.Println(name)

	//}


