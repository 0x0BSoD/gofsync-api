package utils

import (
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func whatIoChecker(param string) io.Writer {
	switch param {
	case "STDOUT":
		return os.Stderr
	case "STDERR":
		return os.Stderr
	case "DISCARD":
		return ioutil.Discard
	default:
		tmp := strings.Split(param, "/")
		folder := strings.Join(tmp[:len(tmp)-1], "/")
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			err = os.Mkdir(folder, 0666)
			if err != nil {
				log.Fatalf("Error on mkdir: %s", err)
			}
		}
		f, err := os.OpenFile(param, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(errors.WithMessage(err, "error opening file"))
		}
		return f
	}
}

func Init(traceHandle *string,
	infoHandle *string,
	warningHandle *string,
	errorHandle *string) {
	config(whatIoChecker(*traceHandle),
		whatIoChecker(*infoHandle),
		whatIoChecker(*warningHandle),
		whatIoChecker(*errorHandle),
	)
}

func config(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(infoHandle, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(warningHandle, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorHandle, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
