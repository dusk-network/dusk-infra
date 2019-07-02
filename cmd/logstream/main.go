package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// var tjson = `{"data": {"code": "warn", "level": "error", "msg": "this is a message", "error": "my error"}}`

func main() {

	uri, err := url.Parse("/var/tmp/logmon.sock")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	c, err := net.Dial("unix", uri.Path)

	if err != nil {
		log.WithError(err).Errorln("problems in connecting to the unix socket")
		return
	}
	defer c.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		j, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		j = strings.Trim(j, "\n")
		switch j {
		case "help":
			help()
		case "quit":
			fmt.Print("$ BYE!\n")
			return
		case "block":
			sendBlock(c)
		case "log":
			sendLog(j, c)
		default:
			if len(strings.Trim(j, " ")) > 0 {
				fmt.Fprintf(os.Stdout, "unrecognized command %s \n", j)
			}
		}
	}
}

var round = 0

func help() {
	fmt.Fprintf(os.Stdout, "log [error | warning | fatal | panic ] | block | quit | help\n")
}

func sendBlock(w io.Writer) {
	t := time.Now()
	msg := make([]byte, 64)
	_, _ = rand.Read(msg)
	blockHash := hex.EncodeToString(msg)
	blockTime := rand.Float64()*3 + 3
	round++

	s := fmt.Sprintf(`{ "data": { "blockTime": %.2f, "blockHash": "%s", "round": %d, "timestamp": "%s" } }`, blockTime, blockHash, round, t.Format(time.RFC3339))
	mwrite(s, w)
}

func sendLog(level string, w io.Writer) {
	t := time.Now()
	level = strings.Trim(level, " ")
	timestamp := t.Format(time.RFC3339)
	s := fmt.Sprintf(`{ "data"" { "error": "this is an error", level": "%s", "msg": "this is a mock message", "timestamp": "%s" } }`, level, timestamp)
	mwrite(s, w)
}

func mwrite(s string, w io.Writer) {
	mw := io.MultiWriter(w, os.Stdout)
	fmt.Fprintln(mw, s)
}
