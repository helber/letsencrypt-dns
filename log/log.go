package log

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

type loggingType struct {
	program   string
	logDir    string
	logStdOut bool
	logDebug  bool
}

var logging loggingType

func init() {
	flag.StringVar(&logging.logDir, "logdir", "/var/log/letsencrypt-dns/", "")
	flag.BoolVar(&logging.logStdOut, "logstdout", false, "Output log to stdout")
	flag.BoolVar(&logging.logDebug, "logdebug", true, "Show file and line number in log")
}

// InitLogs create log files
func InitLogs() {
	if logging.logDebug {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}
	if !logging.logStdOut {
		program := filepath.Base(os.Args[0])
		fname := logging.logDir + program + ".log"
		f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("Can't create log file", err)
		}
		log.SetOutput(f)
	}
}
