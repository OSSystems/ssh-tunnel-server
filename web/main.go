package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"golang.org/x/crypto/ssh"
	"golang.org/x/net/websocket"
)

func copyWorker(dst io.Writer, src io.Reader, doneCh chan<- bool) {
	io.Copy(dst, src)
	doneCh <- true
}

func relayHandler(ws *websocket.Conn) {
	user := ws.Request().URL.Query().Get("user")
	passwd := ws.Request().URL.Query().Get("passwd")

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(passwd),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", "localhost:22", config)
	if err != nil {
		fmt.Println(err)
		ws.Close()
		return
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	session, err := client.NewSession()
	if err != nil {
		session.Close()
		ws.Close()
		return
	}

	sshOut, err := session.StdoutPipe()
	if err != nil {
		session.Close()
		ws.Close()
		return
	}

	sshIn, err := session.StdinPipe()
	if err != nil {
		session.Close()
		ws.Close()
		return
	}

	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		session.Close()
		ws.Close()
		return
	}
	if err := session.Shell(); err != nil {
		session.Close()
		ws.Close()
		return
	}

	doneCh := make(chan bool)

	go copyWorker(sshIn, ws, doneCh)
	go copyWorker(ws, sshOut, doneCh)

	<-doneCh

	client.Close()
	ws.Close()

	<-doneCh
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui", http.StatusMovedPermanently)
	})
	http.Handle("/ws", websocket.Handler(relayHandler))
	http.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir("ui"))))

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}
}
