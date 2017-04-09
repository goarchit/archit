// ArchIt main routine
// Originally work created on 1/3/2017
//

package main

import (
	"github.com/goarchit/archit/cmd"
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"

	// #include <unistd.h>
	"C"
	"runtime"
	"time"
	"unsafe"
)

func init() {
	// Note that the Parser controls all flow, it initiates and processes all commands BEFORE returning
	// Bacially making it the next and final meaningful thing done
	config.ParseCmdLine() // Go parse flags, env values, & config file
}

func main() {
	const wordSize = int(unsafe.Sizeof(uintptr(0)))
	//  All work is done via commands invoked by the parser in function init()
	// Print some basic debug info
	log.Debug("Started ArchIt at ", time.Now())
	log.Debug("Compiled with ", runtime.Version())
	log.Debug("This system has ", runtime.NumCPU(), "cores")
	log.Debug("Word size", wordSize)
	log.Debug("GOARCH ", runtime.GOARCH)
	memSize := C.sysconf(C._SC_PHYS_PAGES) * C.sysconf(C._SC_PAGE_SIZE) / (1024 * 1024)
	log.Debug("System Memory", memSize, "MB")
	if memSize < 2048 {
		log.Error(memSize, "MB of memory detected.  This program utilizes 2GB+ memory arrays.  Running on a system with less than 2.5GB of real memory will likely cause severe system performance impact as you swap to death")
	}

	log.Debug("Wrapping up archit.init()")
	// Must invoke something in package cmd in order for all init()s to be called
	cmd.GoodBye()
}
