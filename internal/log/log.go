package log

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

var lg = logrus.WithField("process", "logtail")

// LinesToRetain represents the number of lines we need to retain from the tail process
var LinesToRetain = 10

// Tailer is a convenience wrapper over a log file tailing
type Tailer struct {
	logFile   string
	file      *os.File
	lock      sync.RWMutex
	closed    bool
	QuitChan  chan error
	Tail      *tail.Tail
	lastLines []*monitor.Param
}

// New creates a *Tailer
func New(logFile string) *Tailer {
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		return nil
	}

	s := &Tailer{
		logFile:   logFile,
		QuitChan:  make(chan error, 1),
		closed:    true,
		lastLines: make([]*monitor.Param, 0, LinesToRetain),
	}

	if err := s.open(); err != nil {
		lg.Panic(err)
	}
	return s
}

func (l *Tailer) open() error {
	var err error
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.closed {
		l.file, err = os.Open(l.logFile)

		if err != nil {
			return err
		}
		l.closed = false
	}
	return nil
}

// Wire first Writes the tail (last 10 lines) of a file and thus spawns a goroutine that pipes the tail of a file to a writer.
// In case of errors within the tailing goroutine it notifies the parent process through an error channel before exiting
// panics if it fails to setup the Tail process. Otherwise it gracefully exits and writes the reason on the Tailer.QuitChan channel.
// Note: since the process is blocking, it should run on a goroutine
func (l *Tailer) Wire(w io.Writer) {
	if err := l.open(); err != nil {
		lg.WithError(err).Errorln(fmt.Sprintf("cannot start tailing log %s. Aborting", l.logFile))
		return
	}

	r := bufio.NewReader(l.file)

	if lines := l.FetchTail(r, LinesToRetain); lines != nil {
		l.lastLines = lines
	}

	l.TailLog(w)
}

// InitialState writes the current state to a Writer.
// It is a connection initialization
func (l *Tailer) InitialState(conn io.Writer) error {
	for _, param := range l.lastLines {
		if param == nil {
			// the first nil signals that there aren't any more lines store (as lastLine is a queue with first elements being the most recent)
			break
		}
		if err := l.Monitor(conn, param); err != nil {
			return err
		}
	}
	return nil
}

// Monitor simply writes a JSON stream to the writer
func (l *Tailer) Monitor(conn io.Writer, m *monitor.Param) error {

	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if _, err := conn.Write(b); err != nil {
		return err
	}

	// saving on the initial state for new incoming writers
	if len(l.lastLines) > 0 {
		_, l.lastLines = l.lastLines[0], l.lastLines[1:]
	}

	l.lastLines = append(l.lastLines, m)
	return nil
}

// FetchTail writes the tail of a reader (i.e. a file) to a writer
func (l *Tailer) FetchTail(r io.Reader, nrLines int) []*monitor.Param {
	s := bufio.NewScanner(r)
	lastLines := make([]*monitor.Param, 0, nrLines)
	for s.Scan() {
		txt := s.Text()
		if len(lastLines) >= nrLines {
			_, lastLines = lastLines[0], lastLines[1:]
		}

		p := newParam(txt)
		lastLines = append(lastLines, p)
	}

	if err := s.Err(); err != nil && err != io.EOF {
		lg.WithError(err).Warnln("could not fetch the log tail. Continuing without")
		return []*monitor.Param{}
	}

	return lastLines
}

// IsOpen checks if the process is open or otherwise
func (l *Tailer) IsOpen() bool {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return !l.closed
}

// TailLog tails a file and writes on a writer
func (l *Tailer) TailLog(w io.Writer) {
	defer l.close()

	logfile := l.file.Name()
	fi, err := l.file.Stat()

	if err != nil {
		l.QuitChan <- err
		return
	}

	cfg := tail.Config{
		Follow: true,
		Poll:   true,
		Location: &tail.SeekInfo{
			Offset: fi.Size(),
			Whence: io.SeekStart,
			// Whence: io.SeekStart,
		},
	}

	l.Tail, err = tail.TailFile(logfile, cfg)

	if err != nil {
		l.QuitChan <- err
		return
	}

	for line := range l.Tail.Lines {
		row := strings.Trim(line.Text, " ")
		if len(row) <= 0 {
			continue
		}
		m := newParam(row)
		if err := l.Monitor(w, m); err != nil {
			l.QuitChan <- err
			return
		}
	}
	l.Shutdown()
}

func (l *Tailer) close() {
	l.lock.Lock()
	defer l.lock.Unlock()
	if !l.closed {
		_ = l.file.Close()
		l.closed = true
	}
}

// Shutdown stops the Tail process
func (l *Tailer) Shutdown() {
	lg.Debugln("shutting down")
	l.close()
	l.QuitChan <- errors.New("Tail process stopped")
	// this triggers a race condition. However we don't care as the process is shutdown anyway
	_ = l.Tail.Stop()
	lg.Debugln("bye")
}

func newParam(m string) *monitor.Param {
	return &monitor.Param{
		Metric:    "tail",
		Value:     m,
		Timestamp: time.Now(),
	}
}
