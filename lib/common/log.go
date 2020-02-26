package common

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

var Log = createLogger()

func createLogger() *logrus.Logger {
	l := logrus.New()

	l.SetOutput(os.Stderr)

	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	l.SetFormatter(customFormatter)

	lls := os.Getenv("TERRAFMT_LOG")
	if lls == "" {
		lls = "WARN"
	}

	ll, err := logrus.ParseLevel(lls)
	if err != nil {
		l.SetLevel(logrus.TraceLevel)
		l.Errorf("defaulting to TRACE: unable to parse `TERRAFMT_LOG` into a valid log level %v", err)
	} else {
		l.SetLevel(ll)
	}

	return l
}

func CaptureRun(f func() error) (stdout, stderr string, err error) {
	// setup stdout
	outpipeR, outpipeW, err := os.Pipe()
	if err != nil {
		return "", "", fmt.Errorf("failed to create pipe for stdout: %w", err)
	}
	outCh := make(chan string)
	go func() {
		buf := bytes.NewBufferString("")
		io.Copy(buf, outpipeR)
		// close read end, since have read everything out
		outpipeR.Close()
		outCh <- buf.String()
	}()

	// replace stdout
	defer func(o *os.File) { os.Stdout = o }(os.Stdout)
	os.Stdout = outpipeW

	// setup stderr
	errpipeR, errpipeW, err := os.Pipe()
	if err != nil {
		return "", "", fmt.Errorf("failed to create pipe for stderr: %w", err)
	}
	errCh := make(chan string)
	go func() {
		buf := bytes.NewBufferString("")
		io.Copy(buf, errpipeR)
		// close read end, since have read everything out
		errpipeR.Close()
		errCh <- buf.String()
	}()

	// replace stdout
	defer func(o *os.File) { os.Stderr = o }(os.Stderr)
	os.Stderr = errpipeW

	// replace log writer
	defer log.SetOutput(log.Writer())
	log.SetOutput(errpipeR)

	// invoke function
	err = f()

	// close pipe write end so that the copy routines could finish
	outpipeW.Close()
	errpipeW.Close()

	return <-outCh, <-errCh, err
}
