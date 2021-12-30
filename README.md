# net-cat


### How to run:

1. Run the **go build TCPChat.go** command in the root directory of the project
2. Run the **./TCPChat $port**

### Description

This project consists on recreating the NetCat in a Server-Client Architecture that can run in a server mode on a specified port listening for incoming connections, and it can be used in client mode, trying to connect to a specified port and transmitting information to the server.

   * NetCat, nc system command, is a command-line utility that reads and writes data across network connections using TCP or UDP. It is used for anything involving TCP, UDP, or UNIX-domain sockets, it is able to open TCP connections, send UDP packages, listen on arbitrary TCP and UDP ports and many more.

Net-cat creates a group chatvia TCP. The project has the following features :

   * TCP connection between server and multiple clients (relation of 1 to many).
    * A name requirement to the client.
    * Control connections quantity.
    * Clients must be able to send messages to the chat.
    * Do not broadcast EMPTY messages from a client.
    * Messages sent, must be identified by the time that was sent and the user name of who sent the message, example : [2020-01-20 15:48:41][client.name]:[client.message]
    * If a Client joins the chat, all the previous messages sent to the chat must be uploaded to the new Client.
    * If a Client connects to the server, the rest of the * Clients must be informed by the server that the Client joined the group.
    * If a Client exits the chat, the rest of the Clients must be informed by the server that the Client left.
    * All Clients must receive the messages sent by other Clients.
    * If a Client leaves the chat, the rest of the Clients must not disconnect.
    * If there is no port specified, then set as default the port 8989. Otherwise, program must respond with usage message: [USAGE]: ./TCPChat $port

### Instructions
* Start TCP server, listen and accept connections
* The project has Go-routines
* The project has channels and Mutexes
* Accepts maximum 10 connections


### Packages
* Only standard library was used

### Usage

    Server: 
    $ go run .
    Listening on the port :8989
    $ go run . 2525
    Listening on the port :2525
    $ go run . 2525 localhost
    [USAGE]: ./TCPChat \$port
    $
    
    Client:
    $ nc localhost 2525


