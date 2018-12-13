package main

import (
	"bufio"
	"log"
	"net"
	_ "strconv"
)
import "fmt"

func main() {

	var chat_server_name = "TCChat"
	//var TCCHAT_WELCOME = "TCCHAT_WELCOME\tWelcome to the "+chat_server_name+"!\n"

	/*
	var TCCHAT_USERIN = "TCCHAT_USERIN\t",User.nickname"\n"
	var TCHAT_USEROUT = "TCCHAT_USEROUT\t",User.nickname"\n"
	var TCCHAT_BCAST = "TCCHAT_BCAST\t",User.nickname"\t",User.message,"\n"*/


	//var register = "TTCHAT_REGISTER"

	type User struct {
		nickname string
		conn net.Conn
		message string
	}

	aconns := make(map[net.Conn]int)
	conns := make(chan net.Conn)
	dconns := make(chan net.Conn)
	msgs := make(chan string)
	i := 0

	fmt.Println("Launching server...", chat_server_name)

	ln, err := net.Listen("tcp", ":8081")

	go func() {
		for{
			conn , err := ln.Accept()
			if err != nil {
				log.Println(err.Error())
			}
			fmt.Printf("Client %v has connected\n", i)
			conns <- conn

		}
	}()

	if err != nil{

		fmt.Print("The server is not listenning, something went wrong..", err)

		return
	}

	for{
		select {
		case conn := <- conns:
			//conn.Write([]byte(TCCHAT_WELCOME + "\n"))
			aconns[conn] = i
			i++

			go func(conn net.Conn, i int) {
				rd := bufio.NewReader(conn)

				for {
					m,err := rd.ReadString('\n')
					if err != nil {
						break
					}
					//fmt.Print("Message Received:", string(m))
					//conn.Write([]byte(m + "\n"))
					msgs <- fmt.Sprintf("Client %v: %v",i,m)
				}
				dconns <-conn
			}(conn,i)

		case msg := <- msgs:
			for conn := range aconns{

				conn.Write([]byte(msg))
			}
		case dconn := <- dconns:
			fmt.Printf("Client %v is gone \n", aconns[dconn])
			delete(aconns,dconn)
		}

	}

	//connection, error := listen.Accept()



	/*for {

		message, _ := bufio.NewReader(connection).ReadString('\n')

		fmt.Print("Message Received:", string(message))

		newmessage := string(message)

		connection.Write([]byte(newmessage + "\n"))

	}*/



}
