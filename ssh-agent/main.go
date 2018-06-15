package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/cobra"
)

func main() {
	var deviceID string
	var mqttServer string
	var sshServer string
	var sshServerPort string
	var identityFile string

	var rootCmd = &cobra.Command{
		Use: "ssh-agent",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	rootCmd.Flags().StringVarP(&deviceID, "device-id", "d", "", "device identity")
	rootCmd.Flags().StringVarP(&mqttServer, "mqtt-server", "m", "", "mqtt server")
	rootCmd.Flags().StringVarP(&sshServer, "ssh-server", "s", "", "ssh server")
	rootCmd.Flags().StringVarP(&sshServerPort, "ssh-server-port", "p", "2221", "ssh server port")
	rootCmd.Flags().StringVarP(&identityFile, "identity-file", "i", "", "identity file")

	rootCmd.MarkFlagRequired("device-id")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

	helpCalled, err := rootCmd.Flags().GetBool("help")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if helpCalled {
		os.Exit(1)
	}

	fmt.Println("Starting agent")
	fmt.Printf("device-id=%s\n", deviceID)

	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s", mqttServer))

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	if token := client.Subscribe(fmt.Sprintf("device/%s", deviceID), 0, func(client mqtt.Client, msg mqtt.Message) {
		go func() {
			port := string(msg.Payload())

			fmt.Printf("Reverse port=%s\n", port)

			args := []string{
				"ssh",
				"-i", identityFile,
				"-o", "StrictHostKeyChecking=no",
				"-nNT",
				"-p", sshServerPort,
				"-R", fmt.Sprintf("%s:localhost:22", port),
				fmt.Sprintf("ssh@%s", sshServer),
			}

			cmd := exec.Command(args[0], args[1:]...)
			e := cmd.Run()
			fmt.Println(e)
		}()
	}); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	select {}
}
