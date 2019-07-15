package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"

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
		cmds := strings.Split(j, " ")
		switch cmds[0] {
		case "help":
			help()
		case "thread":
			sendThreads(cmds[1:], c)
		case "quit":
			fmt.Print("$ BYE!\n")
			return
		case "block":
			sendBlock(c)
		case "log":
			sendLog(cmds[1:], c)
		default:
			if len(strings.Trim(j, " ")) > 0 {
				fmt.Fprintf(os.Stdout, "unrecognized command %s \n", j)
			}
		}
	}
}

var round = 0

func help() {
	fmt.Fprintf(os.Stdout, "log [error | warning | fatal | panic ] | thread [nr] | block | quit | help\n")
}

func sendThreads(params []string, w io.Writer) {
	var n = rand.Intn(1000)
	var err error
	if len(params) > 0 {
		if n, err = strconv.Atoi(params[0]); err != nil {
			fmt.Fprintf(os.Stdout, "invalid parameter for `thread`")
			return
		}
	}
	s := fmt.Sprintf(`{"code": "goroutine", "nr": %d}`, n)
	mwrite(s, w)
}

func sendBlock(w io.Writer) {
	msg := make([]byte, 64)
	_, _ = rand.Read(msg)
	blockHash := hex.EncodeToString(msg)
	blockTime := rand.Float64()*3 + 3
	round++

	s := fmt.Sprintf(`{"code": "round", "blockTime": %.2f, "blockHash": "%s", "round": %d}`, blockTime, blockHash, round)
	mwrite(s, w)
}

func sendLog(params []string, w io.Writer) {
	level := "error"
	if len(params) > 0 {
		level = strings.Trim(params[0], " ")
	}

	s := fmt.Sprintf(`{"code": "warn", "process": "logstream", "error": "this is an error", "level": "%s", "msg": "this is a mock message"}`, level)
	mwrite(s, w)
}

func mwrite(s string, w io.Writer) {
	buf := bytes.NewBufferString(s)

	mw := io.MultiWriter(w, os.Stdout)
	fmt.Fprintln(mw, buf)
}
