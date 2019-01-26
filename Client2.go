package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

// On va utiliser cette liste pour print les messages des clients.
var liste2 []string

// Lire fichier où on stocke la donnée reçue et la retourner.
func lirefichier() (string) {
	file, err := os.Open("client.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		liste2 = append(liste2, scanner.Text())
	}
	return liste2[len(liste2)-1]
}

// Fonction qui supprime le retour chariot et qui sépare les composants du message.
func splitter5(msg string) []string {
	s := strings.Replace(msg, "\n", "", 2)
	s2 := strings.Split(s, "\t")
	return s2

}

// Traitement de tout le message reçu.
func detectionOfTypeUser3(msg string, conn net.Conn) {
	// On le split
	message := splitter5(msg)
	// Le type est le premier champ
	contentType := message[0]
	tentativedeco := 0

	// On reçoit un client, on lui demande son pseudo et on l'enregistre au près du serveur.
	if contentType == "TCCHAT_WELCOME" {
		if tentativedeco > 1 { // Si le client n'arrive pas à entrer dans la chatroom dès le premier essai, ça veut dire qu'il y aura un message d'erreur du serveur

			ioutil.WriteFile("client.txt", []byte(message[2]), 0777)
			msgerreur := (lirefichier())
			fmt.Println(msgerreur)
		} else {
			servern := message[1]
			ioutil.WriteFile("client.txt", []byte(servern), 0777)
			servername := (lirefichier())
			fmt.Println("Welcome to " + servername + "!\n")

			// Récupérer le pseudo de l'utilisateur --> TC CHAT REGISTER
			lecture := bufio.NewReader(os.Stdin)
			print("Tapez votre pseudo:")
			pseudo, err := lecture.ReadString('\n')
			if err != nil {
				print(err)
			}
			tentativedeco += 1
			conn.Write([]byte("TCCHAT_REGISTER\t" + pseudo + "\n"))
			fmt.Println("--------------------------" + servername + " INSTRUCTIONS ------------------------------------\n" + "           Ecrivez votre message, ou tapez 'exit' pour vous déconnecter            " + "-----------------------------------------------------------------------------------\n")
		}
	}
	// Si on reçoit un type autre que TCCHAT WELCOME, ça veut dire que le client est déjà authentifié et qu'il peut écrire un message.
	if contentType != "TCCHAT_WELCOME" {

		// Si on reçoit du serveur un USERIN, on avertit tous les utilisateurs que user1 est entré dans la room;
		if contentType == "TCCHAT_USERIN" {
			nickname := message[1]
			ioutil.WriteFile("client.txt", []byte(nickname), 0777)
			messagefichier := lirefichier()
			fmt.Print(messagefichier + " has entered the room.\n")

		}

		// On récupère les 2 paramètres qui suivent TCCHAT_BCAST pour afficher le message (variable lemessage) d'un utilisateur (nickname)
		if contentType == "TCCHAT_BCAST" {
			nickname := message[1]
			lemessage := message[2]

			// On injecte les deux variables récupérées dans un fichier pour les lire et les print sur le terminal client.
			ioutil.WriteFile("client.txt", []byte(nickname), 0777)
			nickname2 := (lirefichier())
			ioutil.WriteFile("client.txt", []byte(lemessage), 0777)
			sending := lirefichier()
			fmt.Println(nickname2 + " a écrit : " + sending)
		}

		// Affiche les clients qui se sont déconnectés du chat
		if contentType == "TCCHAT_USEROUT" {

			ioutil.WriteFile("client.txt", []byte(message[1]), 0777)
			nickname := lirefichier()

			fmt.Println(nickname + " has left the room.\n")
		}

		// Go routine qui demande à l'utilisateur de saisir son message car le fait d'en recevoir ne devrait pas nous empêcher d'écrire
		ioutil.WriteFile("client.txt", []byte(""), 0777)
		go writeamessage(conn, msg)

	}
}

func writeamessage(conn net.Conn, msg string) {
	lecture3 := bufio.NewReader(os.Stdin)

	messagepayload, _ := lecture3.ReadString('\n')
	ioutil.WriteFile("client.txt", []byte(messagepayload), 0777)
	messagepayload2 := lirefichier()
	if messagepayload2 == "exit" {
		conn.Write([]byte("TCCHAT_DISCONNECT\t" + "\n"))
		print("Exiting!")
		os.Exit(0)
	}

	if messagepayload2 != "exit" {
		conn.Write([]byte("TCCHAT_MESSAGE\t" + messagepayload + "\n"))
	}

}

func main() {
	os.Create("client.txt")
	// Se connecte au serveur
	conn, _ := net.Dial("tcp", "127.0.0.1:8081")

	for {
		// Lire les entrées recues dans le tunnel tcp établi entre le serveur et le client
		reader := bufio.NewReader(conn)

		//Pour lire les msg du serveur
		message, err := reader.ReadString('\n')
		// Si erreur, fermer la session tcp dédiée à ce client
		if err != nil {
			conn.Close()
			break
		}

		err2 := ioutil.WriteFile("client.txt", []byte(message), 0777)
		if err2 != nil {
			// Print l'erreur
			fmt.Println(err2)
		}

		readfromfile, err := ioutil.ReadFile("client.txt")
		if err != nil {
			fmt.Println(err)
		}

		go detectionOfTypeUser3(string(readfromfile), conn)
	}

}
