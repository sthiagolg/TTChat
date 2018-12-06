package main

import (
	"net"
	_ "strconv"
)
import "fmt"
import "bufio"


func main() {

	var chat_server_name = "TCChat"

	//var register = "TTCHAT_REGISTER"


	fmt.Println("Launching server...", chat_server_name)


	listen, err := net.Listen("tcp", ":8081")

	if err != nil{

		fmt.Print("The server is not listenning, something went wrong..", err)

		return
	}

	connection, error := listen.Accept()

	if error == nil {
		connection.Write([]byte("WELCOME TO " + chat_server_name + "!\n"))
	}



	for {

		message, _ := bufio.NewReader(connection).ReadString('\n')

		fmt.Print("Message Received:", string(message))

		newmessage := string(message)

		connection.Write([]byte(newmessage + "\n"))

	}



}

