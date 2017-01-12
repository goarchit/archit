// ArchIt main routine
// Originally work created on 1/3/2017
//

package main

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/parser"
	"github.com/goarchit/archit/cmd"
	"net"
	"strconv"
)

func init() {
	cmd.Define()	// Go define all commands
	config.Read()   // Go parse flags, env values, & config file
	log.Debug("Wrapping up archit.init()")
}

func main() {

	// Print Configuration settings
	log.Info("Final configuration results:")
	log.Info("Port =", parser.Archit.Port)
	log.Info("PublicIP =", parser.Archit.PublicIP)
	log.Debug("Looking up PublicIP...")
	ips, err := net.LookupIP(parser.Archit.PublicIP)
	if err != nil {
		log.Critical("Fatal error: ", err)
		return
	}
	log.Debug("IP(s): ", ips)
	serverip := net.JoinHostPort(ips[0].String(), strconv.Itoa(parser.Archit.Port))
	log.Info("Using server address", serverip)
//	go server.Server(serverip)

//	var s, toip, filename string
//	for {
//		fmt.Printf("Command: ")
//		fmt.Scanf("%s", &s)
//		switch s {
//		case "Hi":
//			fmt.Println("Hello!")
//		case "Send":
//			fmt.Printf("To: ")
//			fmt.Scanf("%s", &toip)
//			fmt.Printf("Filename: ")
//			fmt.Scanf("%s", &filename)
//			go server.SendFile(toip, filename)
//		case "Ping":
//			fmt.Println("Sorry, Ping not implemented")
//		case "Pong":
//			fmt.Println("Pongs are for servers!  Not people!")
//		case "Exit":
//			fmt.Println("Bye now")
//			os.Exit(0)
//		case "Help":
//			fmt.Println("I understand Hi, Send, Ping, Pong, Exit, and Help")
//		default:
//			fmt.Println("Unknown command ", s)
//
//		}
//	}
}
