// ArchIt main routine
// Originally work created on 1/3/2017
//

package main

import (
	"github.com/goarchit/archit/config"
	//	"github.com/goarchit/archit/wallet"
	"github.com/goarchit/archit/cmd"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/parser"
	// #include <unistd.h>
	"C"
	"net"
	"runtime"
	"strconv"
	"time"
	"unsafe"
)

func init() {
	cmd.Define()  // Go define all commands
	config.Read() // Go parse flags, env values, & config file
	log.Debug("Wrapping up archit.init()")
}

func main() {
	const wordSize = int(unsafe.Sizeof(uintptr(0)))
	// Print some basic debug info
	log.Debug("Starting ArchIt at ", time.Now())
	log.Debug("Compiled with ", runtime.Version())
	log.Debug("This system has ", runtime.NumCPU(), "cores")
	log.Debug("Word size", wordSize)
	log.Debug("GOARCH ", runtime.GOARCH)
	memSize := C.sysconf(C._SC_PHYS_PAGES) * C.sysconf(C._SC_PAGE_SIZE) / (1024 * 1024)
	log.Debug("System Memory", memSize, "MB")
	if memSize < 2048 {
		log.Error(memSize, "MB of memory detected.  This program utilizes 1GB+ memory arrays.  Running on a system with less than 2GB of real memory will likely cause severe system performance impact")
	}

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
	//	go wallet.Server()
	//	time.Sleep(2000 * time.Millisecond)
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
