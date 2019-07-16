package log_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"

	"github.com/stretchr/testify/assert"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/log"
)

var tlogs = []struct {
	msg    string
	limit  int
	values []string
	err    error
}{
	{"pippo\npluto", 2, []string{"pippo", "pluto"}, nil},
	{"pippo\npluto\n", 3, []string{"pippo", "pluto"}, nil},
	{"pippo\npluto\npaperino", 2, []string{"pluto", "paperino"}, nil},
}

func TestFetchTail(t *testing.T) {
	l := log.New("")
	for _, tt := range tlogs {
		test := bytes.NewBufferString(tt.msg)
		lines := l.FetchTail(test, tt.limit)

		for i, v := range tt.values {
			if !assert.Equal(t, v, lines[i].Value) {
				assert.FailNowf(t, "error in FetchTail from %s", tt.msg)
			}
		}
	}
}

func TestMonitor(t *testing.T) {
	tmpf, _ := ioutil.TempFile("", ".testlogtail.log")
	fn := tmpf.Name()
	defer os.Remove(fn)

	if _, err := tmpf.Write([]byte("line 1\nline 2\n")); err != nil {
		assert.FailNowf(t, "error in writing to tmp file: %s\n", err.Error())
	}

	l := log.New(fn)
	r := new(bytes.Buffer)

	test := "pippo"
	testM := &monitor.Param{
		Value:     test,
		Timestamp: time.Now(),
		Metric:    "log",
	}

	if !assert.NoError(t, l.Monitor(r, testM)) {
		return
	}

	d := json.NewDecoder(r)
	assert.NoError(t, testJsonReception(d, testM.Value))
}

func testJsonReception(d *json.Decoder, test string) error {
	m := &monitor.Param{}
	if err := d.Decode(m); err != nil {
		return err
	}

	if test != m.Value {
		return fmt.Errorf("Expected %v, got %v", test, m.Value)
	}
	return nil
}

func TestTailLog(t *testing.T) {
	r, w := io.Pipe()
	f, err := ioutil.TempFile("", ".test.log")
	fName := f.Name()
	if !assert.NoError(t, err) {
		return
	}

	defer os.Remove(f.Name())
	l := log.New(fName)

	go l.TailLog(w)
	//giving some time
	select {
	case <-time.After(5 * time.Millisecond):
		_, err = f.Write([]byte("pippo\n"))
		time.Sleep(5 * time.Millisecond)
		if assert.NoError(t, err) {
			d := json.NewDecoder(r)
			assert.NoError(t, testJsonReception(d, "pippo"))

			if !assert.NoError(t, l.Tail.Stop()) {
				assert.FailNow(t, "error in stopping tail process")
			}
			e := <-l.QuitChan
			assert.Error(t, io.EOF, e)
		}
	case err = <-l.QuitChan:
		assert.FailNow(t, "%v", err)
	}
}

func TestShutdown(t *testing.T) {
	_, w := io.Pipe()
	f, err := ioutil.TempFile("", ".test.log")
	fName := f.Name()
	if !assert.NoError(t, err) {
		return
	}

	defer os.Remove(f.Name())
	l := log.New(fName)

	go l.TailLog(w)

	counter := 0
	for {
		select {
		case <-time.After(5 * time.Millisecond):
			_, _ = f.Write([]byte("pippo\n"))
			time.Sleep(5 * time.Millisecond)
			counter++
			if counter > 1 {
				t.FailNow()
				return
			}
			l.Shutdown()

		case <-l.QuitChan:
			assert.False(t, l.IsOpen())
			return
		}
	}
}
