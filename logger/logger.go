package logger

import (
	"flag"
	"gopkg.in/Graylog2/go-gelf.v1/gelf"
	"io"
	"log"
	"os"
)

var (
	INFO  *log.Logger
	WARNING *log.Logger
	ERROR *log.Logger
	host = flag.String("graylog-host", "127.0.0.1:9000", "Host and port of the graylog")
)

func init() {
	fileSettings := os.O_APPEND|os.O_CREATE|os.O_WRONLY
	stdFile, err := os.OpenFile("std_output.log", fileSettings, 0666)
	if err != nil {
		log.Fatal(err)
	}
	errFile, err := os.OpenFile("error_output.log", fileSettings, 0666)
	if err != nil {
		log.Fatal(err)
	}
	gelfWriter, err := gelf.NewWriter(*host)
	if err != nil {
		log.Fatalf("gelf.NewWriter: %s", err)
	}

	stdmw := io.MultiWriter(os.Stdout, stdFile, gelfWriter)
	errmw := io.MultiWriter(os.Stdout, errFile, gelfWriter)
	logSettings := log.Ldate|log.Lmicroseconds|log.LUTC
	INFO = log.New(stdmw, "INFO: ", logSettings)
	WARNING = log.New(stdmw, "WARNING: ", logSettings)
	ERROR = log.New(errmw, "ERROR: ", logSettings)
}
