package log_test

import (
	"bytes"
	"encoding/json"
	"testing"

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

func TestWriteLastLines(t *testing.T) {
	l := log.New("")
	w := new(bytes.Buffer)
	for _, tt := range tlogs {
		test := bytes.NewBufferString(tt.msg)
		err := l.WriteLastLines(test, w, tt.limit)
		if tt.err != err {
			assert.FailNow(t, "expected error %s but got %s", tt.err, err)
			return
		}

		m := monitor.Param{}
		d := json.NewDecoder(w)
		for _, v := range tt.values {
			if assert.NoError(t, d.Decode(&m)) {
				assert.Equal(t, v, m.Value)
			}
		}
	}
}

func TestTailLog(t *testing.T) {}
