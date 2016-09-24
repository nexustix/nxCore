package main

import (
	"fmt"
	"net"
	"time"

	bp "github.com/nexustix/boilerplate"
)

type Server struct {
	Remotes []*Remote
	//MessageQueue Queue
	started bool
}

func (s *Server) Start() {
	if !s.started {
		fmt.Printf("[%v]<-> Starting nxCore\n", time.Now().Unix())
		s.listen()
		fmt.Printf("[%v]<-> Stopping nxCore\n", time.Now().Unix())
		s.started = true
	}
}

func (s *Server) listen() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
		fmt.Printf("[%v]<!> ERROR listening on port\n", time.Now().Unix())
		return
	}
	fmt.Printf("[%v]<-> Started nxCore\n", time.Now().Unix())
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			fmt.Printf("[%v]<!> ERROR accepting client connection\n", time.Now().Unix())
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(connection net.Conn) {
	tmpRemote := Remote{server: s}
	go tmpRemote.Run(connection)
	s.Remotes = append(s.Remotes, &tmpRemote)
}

func (s *Server) Bridgecast(remote *Remote, message string) {
	for _, pReceiver := range s.Remotes {
		for _, sChannel := range remote.Channels {
			//if pReceiver != remote && bp.StringInSlice(pReceiver.Channels, sChannel) {
			if bp.StringInSlice(pReceiver.Channels, sChannel) {
				pReceiver.Send(message)
				break
			}
		}
	}

}
