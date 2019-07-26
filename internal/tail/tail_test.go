package tail_test

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
	"gitlab.dusk.network/dusk-core/node-monitor/internal/tail"

	"github.com/stretchr/testify/assert"
)

var tlogs = []struct {
	msg    string
	limit  int
	values []string
	err    error
}{
	{"level=warning\nlevel=pluto\n", 2, []string{"level=warning", "level=pluto"}, nil},
	{"level=pippo\nlevel=pluto\n", 3, []string{"level=pippo", "level=pluto"}, nil},
	{"level=pippo\nlevel=pluto\nlevel=paperino", 2, []string{"level=pluto", "level=paperino"}, nil},
}

func TestFetchTail(t *testing.T) {
	l := tail.New("")
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

	l := tail.New(fn)
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

var testLogLines = []string{
	`time="2019-06-26T15:15:30Z" level=debug  msg="received regeneration message" process=selection round=13 step=3`,
	`time="2019-06-26T15:15:30Z" level="debug msg="starting selection" process=selection selector state="round: 13 / step: 2`,
	`time="2019-06-26T15:15:31Z" Level = "debug" msg="sending proof" collector round=13 process=generation`,
	`time="2019-06-26T15:15:31Z" level:debug msg="score does not exceed threshold" process=selection`,
}

func TestLevelRegexp(t *testing.T) {
	for _, s := range testLogLines {
		assert.Equal(t, []byte("debug"), tail.LevelRegexp.FindSubmatch([]byte(s))[1])
	}
}

func testJsonReception(d *json.Decoder, test string) error {
	m := monitor.NewParam("tail")
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
	l := tail.New(fName)

	go l.TailLog(w)
	//giving some time
	select {
	case <-time.After(5 * time.Millisecond):
		_, err = f.Write([]byte("level=pippo\n"))
		time.Sleep(5 * time.Millisecond)
		if assert.NoError(t, err) {
			d := json.NewDecoder(r)
			assert.NoError(t, testJsonReception(d, "level=pippo"))

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
	l := tail.New(fName)

	go l.TailLog(w)

	counter := 0
	for {
		select {
		case <-time.After(5 * time.Millisecond):
			_, _ = f.Write([]byte("level=pippo\n"))
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
