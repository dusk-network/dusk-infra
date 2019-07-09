package log

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

var lg = logrus.WithField("process", "logtail")
var LinesToRetain = 10

// LogProc is a convenience wrapper over a log file tailing
type LogProc struct {
	logFile   string
	file      *os.File
	lock      sync.RWMutex
	closed    bool
	QuitChan  chan error
	TailProc  *tail.Tail
	lastLines []*monitor.Param
}

// New creates a *LogProc
func New(logFile string) *LogProc {
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		return nil
	}

	s := &LogProc{
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

func (l *LogProc) open() error {
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

// StreamLog first Writes the tail (last 10 lines) of a file and thus spawns a goroutine that pipes the tail of a file to a writer.
// In case of errors within the tailing goroutine it notifies the parent process through an error channel before exiting
// panics if it fails to setup the Tail process. Otherwise it gracefully exits and writes the reason on the LogProc.QuitChan channel.
// Note: since the process is blocking, it should run on a goroutine
func (l *LogProc) Wire(w io.Writer) {
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

func (l *LogProc) InitialState(conn io.Writer) error {
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
func (l *LogProc) Monitor(conn io.Writer, m *monitor.Param) error {

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
func (l *LogProc) FetchTail(r io.Reader, nrLines int) []*monitor.Param {
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

func (l *LogProc) IsOpen() bool {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return !l.closed
}

// TailLog tails a file and writes on a writer
func (l *LogProc) TailLog(w io.Writer) {
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

	l.TailProc, err = tail.TailFile(logfile, cfg)

	if err != nil {
		l.QuitChan <- err
		return
	}

	for line := range l.TailProc.Lines {
		m := newParam(line.Text)
		if err := l.Monitor(w, m); err != nil {
			l.QuitChan <- err
			return
		}
	}
	l.Shutdown()
}

func (l *LogProc) close() {
	l.lock.Lock()
	defer l.lock.Unlock()
	if !l.closed {
		_ = l.file.Close()
		l.closed = true
	}
}

func (l *LogProc) Shutdown() {
	lg.Debugln("shutting down")
	l.close()
	l.QuitChan <- errors.New("Tail process stopped")
	// this triggers a race condition. However we don't care as the process is shutdown anyway
	_ = l.TailProc.Stop()
	lg.Debugln("bye")
}

func newParam(m string) *monitor.Param {
	return &monitor.Param{
		Metric:    "tail",
		Value:     m,
		Timestamp: time.Now(),
	}
}
