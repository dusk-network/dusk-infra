package main

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"strings"
	"syscall"
	"time"

	logstream "gitlab.dusk.network/dusk-core/node-monitor/api"

	"github.com/sirupsen/logrus"

	"gitlab.dusk.network/dusk-core/node-monitor/web"

	"github.com/namsral/flag"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/aggregator"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/cpu"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/disk"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/ip"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/latency"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/mem"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/tail"
)

type cfg struct {
	debug     bool
	skipAggro bool
	httpAddr  string
	latencyIP string
	logfile   string
	u         *url.URL
	b         *url.URL
	bToken    string
	memProf   bool
	cpuProf   bool
	hostName  string
	hostIP    string
}

var c cfg
var defaultMemProf = "monitor_mem.prof"
var defaultCPUProf = "monitor_cpu.prof"

type logURL struct {
	*url.URL
	defaultLogAddr string
}

func (l *logURL) String() string {
	if l.URL != nil {
		return l.URL.String()
	}
	return l.defaultLogAddr
}

func (l *logURL) Set(value string) error {
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
		defaultLogAddr         = "unix:///var/tmp/logmon.sock"
		defaultAggroAddr       = "https://duskbert.dusk.network"
		WSAddrDesc             = "http service address"
		latencyIPDesc          = "preferred voucher seeder"
		logfileDesc            = "location of the node log file"
		debugMode              = "start in debug mode"
	)

	logURLDesc := "URI of the log monitoring server"
	aggroURLDesc := "URI of the log aggregator"

	var err error
	c.u, err = url.Parse(defaultLogAddr)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	c.b, err = url.Parse(defaultAggroAddr)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	hname, err := os.Hostname()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	ipv4, err := ip.Retrieve()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Server host address setup
	flag.StringVar(&c.httpAddr, "address", defaultWSAddress, WSAddrDesc)
	flag.StringVar(&c.httpAddr, "a", defaultWSAddress, WSAddrDesc+" (shorthand)")
	flag.StringVar(&c.hostName, "hostname", hname, "set the hostname manually")
	flag.StringVar(&c.hostIP, "hostip", ipv4, "set the IP manually")

	// Latency option
	flag.StringVar(&c.latencyIP, "seeder", defaultLatencyProberIP, latencyIPDesc)
	flag.StringVar(&c.latencyIP, "s", defaultLatencyProberIP, latencyIPDesc+" (shorthand)")

	// Logfile Tail option
	flag.StringVar(&c.logfile, "logfile", defaultLogfile, logfileDesc)
	flag.StringVar(&c.logfile, "l", defaultLogfile, logfileDesc+" (shorthand)")

	// Logstream UNIX-SOCKET option
	flag.Var(&logURL{c.u, defaultLogAddr}, "uri-logserver", logURLDesc)
	flag.Var(&logURL{c.u, defaultLogAddr}, "u", logURLDesc+"(shorthand)")

	// Debug options
	flag.BoolVar(&c.debug, "verbose", defaultDebugMode, debugMode)
	flag.BoolVar(&c.debug, "v", defaultDebugMode, debugMode+" (shorthand)")
	flag.BoolVar(&c.cpuProf, "cpu", false, fmt.Sprintf("profile monitor cpu on %s", defaultCPUProf))
	flag.BoolVar(&c.memProf, "mem", false, fmt.Sprintf("profile monitor memory on %s", defaultMemProf))

	// Aggregator options
	flag.Var(&logURL{c.b, defaultAggroAddr}, "bot-aggregator", aggroURLDesc)
	flag.Var(&logURL{c.b, defaultAggroAddr}, "b", aggroURLDesc+"(shorthand)")
	flag.StringVar(&c.bToken, "token", "", "token to authenticate with the bot")
	flag.StringVar(&c.bToken, "t", "", "token to authenticate with the bot (shorthand)")
	flag.BoolVar(&c.skipAggro, "d", false, "disable aggregator (shorthand)")
	flag.BoolVar(&c.skipAggro, "disable-aggregator", false, "disable aggregator")

	flag.Parse()
	if c.debug {
		logrus.SetLevel(logrus.DebugLevel)
		fmt.Println("Setting level to DEBUG")
	}
}

func parseURL(uri string) *url.URL {
	uri = strings.Trim(uri, " ")
	res, err := url.Parse(uri)
	if err != nil {
		fmt.Printf("Malformed url %s: %s\n", uri, err.Error())
		os.Exit(1)
	}
	return res
}

func main() {
	var srv *web.Srv

	if c.u.Scheme == "" {
		fmt.Printf("Unrecognized URL %v\n", c.u.String())
		os.Exit(1)
	}

	checkLatencyProberIP(c.latencyIP)
	m := initMonitors(c)

	if c.bToken != "" && !c.skipAggro {
		wa := aggregator.New(c.b, c.httpAddr, c.bToken, c.hostName, c.hostIP)
		srv = web.New(m, wa)
	} else {
		fmt.Println("Running without aggregator forwarding")
		srv = web.New(m, nil)
	}

	fmt.Printf("Starting up the server at %s\n", c.httpAddr)
	// Handle common process-killing signals so we can gracefully shut down:
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func(c chan os.Signal, monitors []monitor.Mon) {
		cpuprof, cpuActive := profileCPU()
		if cpuActive {
			defer cpuprof.Close()
		}
		memprof, memActive := profileMem()
		// Wait for a SIGINT or SIGKILL:
		sig := <-c
		fmt.Printf("Caught signal %s: shutting down.\n", sig)
		if memActive {
			defer memprof.Close()
			runtime.GC()
			if err := pprof.WriteHeapProfile(memprof); err != nil {
				fmt.Printf("Cannot write memory profiling: %s", err.Error())
			}
		}

		// Stop listening (and unlink the socket if unix type):
		for _, mon := range monitors {
			mon.Shutdown()
		}

		// And we're done:
		os.Exit(0)
	}(sigc, m)

	if err := srv.Serve(c.httpAddr); err != nil {
		fmt.Printf("Error in serving the monitoring data")
		os.Exit(1)
	}
}

func profileCPU() (*os.File, bool) {
	if !c.cpuProf {
		return nil, false
	}

	f, err := os.Create(defaultCPUProf)
	if err != nil {
		fmt.Printf("Problems in creating the CPU profile file %s\n", err.Error())
		return nil, false
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Printf("Problems in starting the CPU profiler: %s", err.Error())
		f.Close()
		return nil, false
	}

	return f, true
}

func profileMem() (*os.File, bool) {
	if !c.memProf {
		return nil, false
	}

	f, err := os.Create(defaultMemProf)
	if err != nil {
		fmt.Printf("Problems in creating the memory profile file %s\n", err.Error())
		return nil, false
	}

	return f, true
}

func checkLatencyProberIP(l string) {
	ip := net.ParseIP(l)
	if ip == nil {
		fmt.Printf("invalid IP for the voucher seeder: %s", l)
		os.Exit(1)
	}
}

func initMonitors(c cfg) []monitor.Mon {
	mons := make([]monitor.Mon, 0)
	// if the url is specified we create the logstream server
	mons = append(
		mons,
		logstream.New(c.u),
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
		monitor.New(
			&disk.Disk{},
			5*time.Second,
			"disk",
		),
	)

	l := latency.New(c.latencyIP)
	if err := l.(*latency.Latency).ProbePriviledges(); err == nil {
		m := monitor.New(l, 10*time.Second, "latency")
		mons = append(mons, m)
	} else {
		fmt.Println("Cannot setup the latency prober. Are you running with enough proviledges?")
		os.Exit(3)
	}

	// if the logfile does not exist we don't add it to the processes
	if l := tail.New(c.logfile); l != nil {
		mons = append(mons, l)
	} else {
		fmt.Println("Logfile not found. Log screening cannot be started")
	}

	return mons
}
