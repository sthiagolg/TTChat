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
	var TCCHAT_WELCOME = "TCCHAT_WELCOME\tWelcome to the "+chat_server_name+"!\n"

	/*
	var TCCHAT_USERIN = "TCCHAT_USERIN\t"+User.nickname+"\n"
	var TCHAT_USEROUT = "TCCHAT_USEROUT\t"+User.nickname+"\n"
	var TCCHAT_BCAST = "TCCHAT_BCAST\t"+User.nickname+"\t"+User.message+"\n"


	*/


	//var register = "TTCHAT_REGISTER"

	type User struct {
		nickname string
		connection net.Conn
		message string
	}

	aconns := make(map[User]int)
	conns := make(chan User)
	dconns := make(chan User)
	msgs := make(chan string)
	i := 0
	//Users := make(chan[User])

	fmt.Println("Launching server...", chat_server_name)

	ln, err := net.Listen("tcp", ":8081")

	go func() {
		for{
			conn , err := ln.Accept()
			conn.Write([]byte(TCCHAT_WELCOME + "\n"))
			if err != nil {
				log.Println(err.Error())
			}
			rd := bufio.NewReader(conn)
			name,err := rd.ReadString('\n')
			user := User{nickname:name,connection:conn}
			fmt.Printf( user.nickname+ " '-> has connected\n")
			conns <- user

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
			//msgs <- TCCHAT_WELCOME
			aconns[conn] = i
			i++
			user := conn
			go func(user User, i int) {
				rd := bufio.NewReader(user.connection)

				for {
					m,err := rd.ReadString('\n')
					if err != nil {
						break
					}
					//fmt.Print("Message Received:", string(m))
					//conn.Write([]byte(m + "\n"))
					msgs <- fmt.Sprintf("%v '-> %v",user.nickname,m)
					//msgs <- m
				}
				dconns <-user
			}(user,i)

		case msg := <- msgs:
			for user := range aconns{
				user.connection.Write([]byte(msg))
				fmt.Println(msg)

			}
		case dconn := <- dconns:
			fmt.Printf("%v '-> is gone \n", dconn.nickname)
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
