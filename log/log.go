package log

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"
)

type loggingType struct {
	program   string
	logDir    string
	logFile   string
	logStdOut bool
	logDebug  bool
	logQuiet  bool
}

var logging loggingType

func init() {
	flag.StringVar(&logging.logDir, "logdir", "/var/log/letsencrypt-dns/", "")
	flag.BoolVar(&logging.logStdOut, "logstdout", false, "Output log to stdout")
	flag.BoolVar(&logging.logDebug, "logdebug", false, "Show file and line number in log")
	flag.StringVar(&logging.logFile, "logfile", "", "if logfile supplied will override default (logdir/program name)")
	flag.BoolVarP(&logging.logQuiet, "quiet", "q", false, "quiet (mute log)")
}

// InitLogs create log files
func InitLogs() {
	if logging.logQuiet {
		log.SetOutput(ioutil.Discard)
	}
	if logging.logDebug {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}
	if !logging.logStdOut {
		program := filepath.Base(os.Args[0])
		fname := logging.logDir + program + ".log"
		if logging.logFile != "" {
			fname = logging.logFile
		}
		f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("Can't create log file", err)
		}
		log.SetOutput(f)
	}
}
