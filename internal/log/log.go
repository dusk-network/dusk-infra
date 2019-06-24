package log

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

var lg = logrus.WithField("process", "logtail")

// LogProc is a convenience wrapper over a log file tailing
type LogProc struct {
	logFile  string
	QuitChan chan error
	TailProc *tail.Tail
}

// New creates a *LogProc
func New(logFile string) *LogProc {
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		return nil
	}

	s := &LogProc{
		logFile:  logFile,
		QuitChan: make(chan error),
	}

	return s
}

// StreamLog first Writes the tail (last 10 lines) of a file and thus spawns a goroutine that pipes the tail of a file to a writer.
// In case of errors within the tailing goroutine it notifies the parent process through an error channel before exiting
// panics if it fails to setup the Tail process. Otherwise it gracefully exits and writes the reason on the LogProc.QuitChan channel.
// Note: since the process is blocking, it should run on a goroutine
func (l *LogProc) Wire(w io.WriteCloser) {
	defer w.Close()

	file, err := os.Open(l.logFile)
	if err != nil {
		lg.WithError(err).Errorln(fmt.Sprintf("cannot start tailing log %s. Aborting", l.logFile))
		return
	}

	r := bufio.NewReader(file)
	if err := l.WriteLastLines(r, w, 10); err != nil {
		lg.WithError(err).Errorln(fmt.Sprintf("cannot read last lines of the %s. Aborting", l.logFile))
	}

	l.TailLog(file, w, true)
}

// Monitor simply writes a JSON stream to the writer
func (l *LogProc) Monitor(w io.Writer, m *monitor.Param) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}

// WriteTail writes the tail of a reader (i.e. a file) to a writer
func (l *LogProc) WriteLastLines(r io.Reader, w io.Writer, nrLines int) error {
	s := bufio.NewScanner(r)
	lines := make([]string, 0)
	for s.Scan() {
		txt := s.Text()

		if len(lines) == nrLines {
			_, lines = lines[0], lines[1:]
		}

		lines = append(lines, txt)

	}

	if err := s.Err(); err != nil && err != io.EOF {
		return err
	}

	for _, line := range lines {
		m := newParam(line)
		if err := l.Monitor(w, m); err != nil {
			return err
		}
	}
	return nil
}

// TailLog tails a file and writes on a writer
func (l *LogProc) TailLog(f *os.File, w io.Writer, closeOnExit bool) {
	if closeOnExit {
		defer f.Close()
	}

	logfile := f.Name()
	fi, err := f.Stat()

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
	l.Disconnect()
}

func (l *LogProc) Disconnect() {
	l.QuitChan <- errors.New("Tail process stopped")
}

func newParam(m string) *monitor.Param {
	return &monitor.Param{
		Metric:    "log",
		Value:     m,
		Timestamp: time.Now(),
	}
}
