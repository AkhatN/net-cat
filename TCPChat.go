package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	clients  = make(map[net.Conn]string)
	join     = make(chan net.Conn)
	conns    = make(chan net.Conn)
	delconns = make(chan net.Conn)
	resent   = make(chan net.Conn)
	mutex    sync.Mutex
	counter  int = 0
)

func main() {
	args := os.Args[1:]

	if len(args) != 1 && len(args) != 0 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	port := ":8989"
	if len(args) == 1 {
		n, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Cannot use this port")
			fmt.Println("[USAGE]: ./TCPChat $port")
			return
		}
		if n < 0 {
			fmt.Println("Cannot use this port")
			fmt.Println("[USAGE]: ./TCPChat $port")
			return
		}
		port = fmt.Sprint(":" + args[0])
	}

	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Println(err)
		return
	}
	defer ln.Close()

	log.Printf("Server is listening on %s port\n", port)

	message, err := ioutil.ReadFile("welcome.txt")
	if err != nil {
		log.Println(err)
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err.Error())
				conn.Close()
			}
			conns <- conn
		}
	}()

	for {
		select {
		//Handle incoming connections
		case conn := <-conns:
			mutex.Lock()
			if counter >= 10 {
				log.Println("Room if full!")
				conn.Write([]byte("Room is full, please connect later."))
				conn.Close()
				mutex.Unlock()
				continue
			}
			counter++
			mutex.Unlock()

			//After we have the connection, we accept name of a new user
			go func(conn net.Conn) {
				conn.Write(message)
				n := bufio.NewReader(conn)
				var name string
				for {
					conn.Write([]byte("[ENTER YOUR NAME:]"))
					name, err = n.ReadString('\n')
					if err != nil {
						log.Println(err.Error())
						mutex.Lock()
						counter--
						mutex.Unlock()
						return
					}

					if !isNameValid(conn, name) {
						continue
					}

					name = name[:len(name)-1]
					mutex.Lock()
					clients[conn] = string(name)
					mutex.Unlock()
					break
				}

				showHistory(conn)
				join <- conn
				go handleMessages(conn, string(name))

			}(conn)

			//Deletes users who left the chat
		case dlconn := <-delconns:
			m := fmt.Sprint("\n" + clients[dlconn] + " has left our chat...\n")
			for con := range clients {
				if dlconn != con {
					con.Write([]byte(m))
					t1 := time.Now().Format("2006-01-02 15:04:05")
					m1 := fmt.Sprintf("[%v][%v]:", t1, clients[con])
					con.Write([]byte(m1))
				}
			}
			mutex.Lock()
			counter--
			delete(clients, dlconn)
			mutex.Unlock()

			//Informs everyone that new user has joined the chat
		case joined := <-join:
			m := fmt.Sprint(clients[joined] + " has joined our chat...\n")
			log.Print(m)
			for con := range clients {
				if joined != con {
					con.Write([]byte("\n" + m))
					t1 := time.Now().Format("2006-01-02 15:04:05")
					m1 := fmt.Sprintf("[%v][%v]:", t1, clients[con])
					con.Write([]byte(m1))
				}
			}
		}
	}
}

//handleMessages
func handleMessages(conn net.Conn, name string) {
	defer conn.Close()
	for {
		t1 := time.Now().Format("2006-01-02 15:04:05")
		m1 := fmt.Sprintf("[%v][%v]:", t1, name)
		conn.Write([]byte(m1))
		rd := bufio.NewReader(conn)
		mess, err := rd.ReadString('\n')

		if err != nil {
			log.Println(name+" has left our chat", err)
			break
		}

		if !isMsValid(mess) {
			continue
		}

		mess = mess[:(len(mess) - 1)]
		t2 := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("\n"+"[%v][%v]:%v\n", t2, name, mess)
		m2 := fmt.Sprintf("[%v][%v]:%v\n", t2, name, mess)

		fmt.Print(m2)

		saveHistory(m2)

		for con := range clients {
			if conn != con {
				con.Write([]byte(m))
				t3 := time.Now().Format("2006-01-02 15:04:05")
				m3 := fmt.Sprintf("[%v][%v]:", t3, clients[con])
				con.Write([]byte(m3))
			}
		}
	}

	delconns <- conn
}

//savehistory saves history of every message of the chat
func saveHistory(message string) {
	f, err := os.OpenFile("history.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println(err)
	}
	f.WriteString(message)
	if err := f.Close(); err != nil {
		log.Println(err)
	}
}

//showhistory shows history of the chat to new clients
func showHistory(conn net.Conn) {
	f, err := os.OpenFile("history.txt", os.O_APPEND|os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println(err)
	}

	temp := make([]byte, 1024*4)
	for {
		_, err := f.Read(temp)
		if err == io.EOF { // если конец файла
			break // выходим из цикла
		}
		conn.Write(temp)
	}
	if err := f.Close(); err != nil {
		log.Println(err)
	}
}

//isNameValid checks a name for validation
func isNameValid(conn net.Conn, name string) bool {
	if name == "\n" {
		log.Println("Invalid name")
		conn.Write([]byte("Invalid name, cannot accept empty name.\n"))
		return false
	}

	isname := false
	for i, v := range name {
		if i == len(name)-1 {
			break
		}

		if v != ' ' {
			isname = true
			break
		}
	}

	if !isname {
		log.Println("Invalid name")
		conn.Write([]byte("Invalid name, cannot accept empty name.\n"))
		return false
	}

	name = name[:len(name)-1]

	exists := false
	for _, n := range clients {
		if n == string(name) {
			exists = true
			break
		}
	}

	if exists {
		log.Println("Invalid name")
		conn.Write([]byte("This name already exists.\nTry another one.\n"))
		return false
	}

	return true
}

//isMsValid checks message for validation
func isMsValid(mess string) bool {
	isok := false
	for i, v := range mess {
		if i == len(mess)-1 {
			break
		}
		if v != ' ' {
			isok = true
			break
		}
	}

	if !isok {
		return false
	}

	return true
}
