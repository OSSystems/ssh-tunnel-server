package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	sshserver "github.com/gliderlabs/ssh"
)

func getTargetSessionInfo(s string) (username string, host string, port int) {
	parts := strings.SplitN(s, "@", 2)

	username = parts[0]

	parts = strings.SplitN(parts[1], ":", 2)
	switch len(parts) {
	case 1:
		host = parts[0]
		port = 22
	default:
		host = parts[0]
		port, _ = strconv.Atoi(parts[1])
	}

	return username, host, port
}

func readPassword(stream io.Reader) string {
	reader := bufio.NewReader(stream)
	passwd, _ := reader.ReadString('\r')
	return strings.TrimSpace(passwd)
}

func GetFreePort() (port int, err error) {
	ln, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		return 0, err
	}
	port = ln.Addr().(*net.TCPAddr).Port
	err = ln.Close()
	return
}

func handleSession(s sshserver.Session) {
	username, host, port := getTargetSessionInfo(s.User())

	if host == "" {
		io.WriteString(s, "Invalid target session: device id is missing\n")
		s.Close()
		return
	}

	freePort, _ := GetFreePort()

	fmt.Printf("PUBLISH device/%s %d\n", host, freePort)

	if token := client.Publish(fmt.Sprintf("device/%s", host), 0, false, fmt.Sprintf("%d", freePort)); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	time.Sleep(5 * time.Second)

	io.WriteString(s, "password: ")
	passwd := readPassword(s)
	io.WriteString(s, "\n")

	time.Sleep(5 * time.Second)

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(passwd),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	port = freePort
	host = "localhost"

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		fmt.Println(err)
		io.WriteString(s, fmt.Sprintf("Failed to connect to %s@%s:%d: %s\n", username, host, port, err.Error()))
		s.Close()
		return
	}

	sshClient, err := conn.NewSession()
	if err != nil {
		fmt.Println(err)
	}

	pty, winCh, isPty := s.Pty()

	if isPty {
		err = sshClient.RequestPty(pty.Term, pty.Window.Height, pty.Window.Width, ssh.TerminalModes{})
		if err != nil {
			fmt.Println(err)
		}

		go func() {
			for win := range winCh {
				if err = sshClient.WindowChange(win.Height, win.Width); err != nil {
					fmt.Println(err)
				}
			}
		}()

		stdin, err := sshClient.StdinPipe()
		if err != nil {
			fmt.Println(err)
		}

		stdout, err := sshClient.StdoutPipe()
		if err != nil {
			fmt.Println(err)
		}

		go func() {
			if _, err = io.Copy(stdin, s); err != nil {
				fmt.Println(err)
			}
		}()

		go func() {
			if _, err = io.Copy(s, stdout); err != nil {
				fmt.Println(err)
			}
		}()

		if err = sshClient.Shell(); err != nil {
			fmt.Println(err)
		}

		if err = sshClient.Wait(); err != nil {
			fmt.Println(err)
		}
	}
}

var client mqtt.Client

func main() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetAutoReconnect(true)
	client = mqtt.NewClient(opts)

	for {
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Println(token.Error())
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}

	sshserver.Handle(handleSession)

	log.Println("Listening in 22")

	log.Fatal(sshserver.ListenAndServe(":22", nil, sshserver.HostKeyFile("key.pem")))
}
