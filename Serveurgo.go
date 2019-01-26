package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// Initialisation de la map qui établit la correspondance entre l'utilisateur et son id. de connexion
var mapUsers = make(map[string]net.Conn)

// Fonction qui supprime le retour chariot et qui sépare les composants du message
func splitter(msg string) []string {
	s := strings.Replace(msg, "\n", "", 2)
	s2 := strings.Split(s, "\t")
	return s2

}

func handleConnection(c net.Conn) { //gérer chaque connexion de chaque clt
	for {
		message, err := bufio.NewReader(c).ReadString('\n') //pour lire TCCHAT_REGISTER DU CLIENT
		if err != nil { // si erreur, fermer la session tcp dédiée à ce client
			c.Close()
			break
		}
		go detectionOfTypeServer(message, c)
	}
}

// Traitement de tout le message reçu
func detectionOfTypeServer(msg string, conn net.Conn) {
	message := splitter(msg)  // On le split
	contentType := message[0] // Le type est le premier champ

	if contentType == "TCCHAT_REGISTER" {

		//On récupère le pseudo de l'utilisateur
		nickname := message[1]

		//Vérifier qu'il l'a bien saisi
		if nickname == "" {
			message := "TCCHAT_WELCOME\tTCCHAT\tVous devez saisir un nom d'utilisateur, merci de réessayer\n"
			print (message)
			conn.Write([]byte(message + "\n"))
		}

		// Vérifier que le pseudo n'est pas déjà utilisé
		for k,_ := range mapUsers {
			if nickname == k {
				print (k)
				message := "TCCHAT_WELCOME\tTCCHAT\tPseudo déjà pris, merci de réessayer\n"
				conn.Write([]byte(message + "\n"))
			}
		}
		// S'il n y a pas de problème:
		mapUsers[nickname] = conn
		messageIn := "TCCHAT_USERIN\t" + nickname
		for _, v := range mapUsers {
			v.Write([]byte(messageIn + "\n")) // Envoyer en broadcast un userin

		}

	}

	if contentType == "TCCHAT_MESSAGE" {
		messagepayload := message[1]
		// Il faut trouver l'utilisateur qui l'a envoyé
		for k, v := range mapUsers {
			if v == conn {
				nickname := k
				messageIn := "TCCHAT_BCAST" +"\t" + nickname + "\t" + messagepayload
				//conn.Write([]byte(newmessage + "\n"))
				for _, v := range mapUsers {
					v.Write([]byte(messageIn + "\n")) // Envoyer en broadcast le message de l'utilisateur

				}
			}
		}
	}

	if contentType == "TCCHAT_DISCONNECT" {
		// Il faut trouver l'utilisateur qui l'a envoyé
		for k, v := range mapUsers {
			if v == conn {
				delete(mapUsers, k)
				conn.Close()
				nickname := k
				messageIn := "TCCHAT_USEROUT\t" + nickname + "\t"
				//conn.Write([]byte(newmessage + "\n"))
				for _, v := range mapUsers {
					v.Write([]byte(messageIn + "\n")) // Envoyer en broadcast le message de l'utilisateur

				}
			}
		}
	}

}

func main() {

	fmt.Println("Serveur en cours de synchronisation...")
	//le serveur écoute sur toutes les interfaces
	ln, err := net.Listen("tcp", "127.0.0.1:8081")
	if err != nil {
		print (err)
	}
	fmt.Println("Serveur synchronisé! Attente de clients...")

	// Initialisation du nombre de connections
	numberofConnections := 0

	for {
		conn, _ := ln.Accept() // accepte les connexions
		numberofConnections += 1

		if numberofConnections>= 30 { //Refuse la connexion si on a atteint 30 clients connectés
			conn.Write([]byte("Nombre maximum d'utilisateurs atteint, veuillez réessayer ultérieurement\n"))
			conn.Close()
		}

		print("Nouvelle connexion:", conn , "\n") // Identifier les différentes connexions

		message := "TCCHAT_WELCOME\tTCCHAT"
		conn.Write([]byte(message + "\n"))
		go handleConnection(conn)

	}


}
