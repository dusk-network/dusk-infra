package log

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"

	"github.com/hpcloud/tail"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

// LogProc is a convenience wrapper over a log file tailing
type LogProc struct {
	logFile  string
	QuitChan chan error
	TailProc *tail.Tail
}

// New creates a *LogProc
func New(logFile string) *LogProc {
	s := &LogProc{
		logFile:  logFile,
		QuitChan: make(chan error),
	}

	return s
}

// StreamLog first Writes the tail (last 10 lines) of a file and thus spawns a goroutine that pipes the tail of a file to a writer.
// In case of errors within the tailing goroutine it notifies the parent process through an error channel before exiting
// Returns error if the synchronous operation of opening the file and reading the last 10 lines fails
func (l *LogProc) StreamLog(w io.Writer, closeOnExit bool) error {
	file, err := os.Open(l.logFile)
	if err != nil {
		return err
	}

	r := bufio.NewReader(file)
	if err := l.WriteLastLines(r, w, 10); err != nil {
		return err
	}

	readyChan := make(chan struct{})
	go l.TailLog(file, w, closeOnExit, readyChan)
	return nil
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
func (l *LogProc) TailLog(f *os.File, w io.Writer, closeOnExit bool, ready chan struct{}) {
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

	ready <- struct{}{}

	for line := range l.TailProc.Lines {
		m := newParam(line.Text)
		if err := l.Monitor(w, m); err != nil {
			l.QuitChan <- err
			return
		}
	}
	l.QuitChan <- errors.New("Tail process stopped")
}

func newParam(m string) *monitor.Param {
	return &monitor.Param{
		Metric:    "log",
		Value:     m,
		Timestamp: time.Now(),
	}
}
