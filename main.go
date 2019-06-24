package main

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"gitlab.dusk.network/dusk-core/node-monitor/web"

	"github.com/namsral/flag"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/cpu"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/disk"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/log"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/mem"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

const defaultLogAddress = "unix:///var/logmon.sock"

type cfg struct {
	debug     bool
	httpAddr  string
	latencyIP string
	logfile   string
	u         *url.URL
}

var c cfg

type logUrl struct {
	*url.URL
}

func (l *logUrl) String() string {
	if l.URL != nil {
		return l.URL.String()
	}
	return defaultLogAddress
}

func (l *logUrl) Set(value string) error {
	lURL, err := url.Parse(value)
	if err != nil {
		return err
	}
	*l.URL = *lURL
	return nil
}

func init() {
	c = cfg{}
	const (
		defaultDebugMode       = false
		defaultLogfile         = "/var/log/node.log"
		defaultWSAddress       = "localhost:8080"
		defaultLatencyProberIP = "178.62.193.89"
		WSAddrDesc             = "http service address"
		latencyIPDesc          = "preferred voucher seeder"
		logfileDesc            = "location of the node log file"
		debugMode              = "start in debug mode"
	)

	logUrlDesc := "URI of the log monitoring server (default \"" + defaultLogAddress + "\")"

	var err error
	c.u, err = url.Parse(defaultLogAddress)
	if err != nil {
		panic(err)
	}

	flag.StringVar(&c.httpAddr, "address", defaultWSAddress, WSAddrDesc)
	flag.StringVar(&c.httpAddr, "a", defaultWSAddress, WSAddrDesc+" (shorthand)")
	flag.StringVar(&c.latencyIP, "seeder", defaultLatencyProberIP, latencyIPDesc)
	flag.StringVar(&c.latencyIP, "s", defaultLatencyProberIP, latencyIPDesc+" (shorthand)")
	flag.StringVar(&c.logfile, "logfile", defaultLogfile, logfileDesc)
	flag.StringVar(&c.logfile, "l", defaultLogfile, logfileDesc+" (shorthand)")
	flag.Var(&logUrl{c.u}, "uri-logserver", logUrlDesc)
	flag.Var(&logUrl{c.u}, "u", logUrlDesc+"(shorthand)")
	flag.BoolVar(&c.debug, "verbose", defaultDebugMode, debugMode)
	flag.BoolVar(&c.debug, "v", defaultDebugMode, debugMode+" (shorthand)")

	flag.Parse()
	if c.debug {
		logrus.SetLevel(logrus.DebugLevel)
		fmt.Println("Setting level to DEBUG")
	}
}

func main() {
	// checkPrivileges()
	if c.u.Scheme == "" {
		fmt.Printf("Unrecognized URL %v\n", c.u.String())
		os.Exit(1)
	}

	checkLatencyProberIP(c.latencyIP)
	m := initMonitors(c)
	srv := &web.Srv{
		Monitors: m,
	}
	fmt.Printf("Starting up the server at %v\n", c.httpAddr)
	if err := srv.Serve(strings.Trim(c.httpAddr, " ")); err != nil {
		fmt.Printf("Error in serving the monitoring data")
		os.Exit(1)
	}
}

func checkLatencyProberIP(l string) {
	ip := net.ParseIP(l)
	if ip == nil {
		fmt.Printf("invalid IP for the voucher seeder: %s", l)
		os.Exit(1)
	}
}

func checkPrivileges() {
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()

	if err != nil {
		fmt.Printf("Unexpected error in checking priviledges: %v", err)
		os.Exit(1)
	}

	// 0 = root, 501 = non-root user
	if i, err := strconv.Atoi(string(output[:len(output)-1])); err != nil {
		fmt.Printf("Unexpected error in checking priviledges: %v", err)
		os.Exit(1)
	} else if i != 0 {
		fmt.Println("Not enough priviledges. Are you running as root?")
		os.Exit(3)
	}
}

func initMonitors(c cfg) []monitor.Mon {
	mons := make([]monitor.Mon, 0)
	mons = append(
		mons,
		monitor.New(
			&cpu.CPU{},
			5*time.Second,
			"cpu",
		),
		monitor.New(
			&mem.Mem{},
			8*time.Second,
			"mem",
		),
		// monitor.New(
		// 	latency.New(c.latencyIP),
		// 	10,
		// 	"latency",
		// ),
		monitor.New(
			&disk.Disk{},
			5*time.Second,
			"disk",
		),
	)

	// if the logfile does not exist we don't add it to the processes
	if l := log.New(c.logfile); l != nil {
		mons = append(mons, l)
	} else {
		fmt.Println("Logfile not found. Log screening cannot be started")
	}
	return mons
}
