// ArchIt logging routine
// Originally work created on 1/7/2017
//

package log

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"os"
)

// LogLevel 0 = Critical, die after logging
// LogLevel 1 = Error, but keep running (maybe)
// LogLevel 2 = Wanrting, but keep running (likely)
// Loglevel 3 = Informational
// Loglevel 4 = Debug

var logLevel int
var verbose int
var logger = logs.NewLogger(10000)

func Setup(ll int, logfilename string, v int, resetLog bool) {
	if resetLog {
		os.Remove(logfilename)
	}
	value := "{\"filename\":\"" + logfilename + "\",\"perm\":\"0644\"}"
	logger.SetLogger("file", value)
	logLevel = ll
	verbose = v
	if verbose == 1 {
		Console("Verbosity turned on, expect to see Informational messages")
	}
	if verbose == 2 {
		Console("VeryVerbose turned on, be prepared for a LOT of messages")
	}
	Info("Logging start in", logfilename)
}

func Critical(s ...interface{}) {
	line := fmt.Sprint(s)
	logger.Critical(line)
	fmt.Println("Critical Error:", line)
	os.Exit(1)
}

func Error(s ...interface{}) {
	if logLevel >= 1 {
		line := fmt.Sprint(s)
		logger.Error(line)
		fmt.Println("Error:", line)
	}
}

func Warning(s ...interface{}) {
	if logLevel >= 2 {
		line := fmt.Sprint(s)
		logger.Warning(line)
		fmt.Println("Warning:", line)
	}
}

func Info(s ...interface{}) {
	if logLevel >= 3 {
		line := fmt.Sprint(s)
		logger.Info(line)
		if verbose >= 1 {
			println(line)
		}
	}
}

func Debug(s ...interface{}) {
	if logLevel >= 4 {
		line := fmt.Sprint(s)
		logger.Debug(line)
		if verbose >= 2 {
			println(line)
		}
	}
}

func Trace(s ...interface{}) {
	if logLevel >= 5 {
		line := fmt.Sprint(s)
		logger.Trace(line)
		if verbose >= 2 {
			println(line)
		}
	}
}

func Console(s ...interface{}) {
	line := fmt.Sprint(s)
	if logLevel >= 4 {
		logger.Debug(line)
	}
	println(line)
}

func println(line string) {
	length := len(line)
	if length >= 2 && line[0] == '[' && line[length-1] == ']' {
		fmt.Println(line[1 : length-1])
	}
}
