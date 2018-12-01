package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/crgimenes/goconfig"
)

var (
	in  = os.Stdin
	out = os.Stdout
	buf = make([]byte, 4*1024)
)

func main() {
	type configFlags struct {
		Input   string `json:"i" cfg:"i" cfgDefault:"stdin" cfgHelper:"input from"`
		Output  string `json:"o" cfg:"o" cfgDefault:"stdout" cfgHelper:"output to"`
		OutFile string `json:"f" cfg:"f" cfgRequired:"true" cfgHelper:"output file name"`
	}

	cfg := configFlags{}
	goconfig.PrefixEnv = "YPIPE"
	err := goconfig.Parse(&cfg)
	if err != nil {
		fatal(err.Error())
	}

	f, err := os.Create(cfg.OutFile)
	if err != nil {
		fatal(err.Error())
	}
	fWriter := bufio.NewWriter(f)
	reader := bufio.NewReader(in)
	writer := bufio.NewWriter(out)

	n := 0
	for {
		n, err = reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				fatal(err.Error())
			}
			if n == 0 {
				break
			}
		}
		write(writer, buf[:n])
		write(fWriter, buf[:n])
	}
	flush(writer)
	flush(fWriter)
}

func fatal(msg string) {
	fmt.Fprintf(os.Stderr, "%v\n", msg) // nolint
	os.Exit(-1)
}

func flush(w *bufio.Writer) {
	err := w.Flush()
	if err != nil {
		fatal(err.Error())
	}
}

func write(w io.Writer, buf []byte) {
	n, err := w.Write(buf)
	if err != nil {
		fatal(err.Error())
	}
	if n < len(buf) {
		fatal(err.Error())
	}
}
