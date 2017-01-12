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
	value := "{\"filename\":\""+logfilename+"\"}"
	logger.SetLogger("file", value)
	logLevel = ll 
	verbose = v
	if (verbose == 1) {
		Warning("Verbosity turned on, expect to see Informational messages")
	}
	if (verbose == 2) {
		Warning("VeryVerbose turned on, be prepared for a LOT of messages")
	}
	Info("Logging start in",logfilename)
}

func Critical(s ...interface{}) {
	logger.Critical(fmt.Sprint(s))
 	panic(fmt.Sprint(s))	
}

func Error(s ...interface{}) {
	if logLevel >= 1 {
		logger.Error(fmt.Sprint(s))
		fmt.Println("Error:", s)
	}
}

func Warning(s ...interface{}) {
	if logLevel >= 2 {
		logger.Warning(fmt.Sprint(s))
		fmt.Println("Warning:", s)
	}
}

func Info(s ...interface{}) {
	if logLevel >= 3 {
		logger.Info(fmt.Sprint(s))
		if verbose >= 1 {
			fmt.Println(s)
		}
	}
}

func Debug(s ...interface{}) {
	if logLevel >= 4 {
		logger.Debug(fmt.Sprint(s))
		if verbose >= 2 {
			fmt.Println(s)
		}
	}
}
