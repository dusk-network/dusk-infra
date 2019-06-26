package nodelogger_test

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"testing"
	"time"

	nodelogger "gitlab.dusk.network/dusk-core/node-monitor/api"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/test"

	"github.com/stretchr/testify/assert"
)

const file = "/tmp/dusk.soc"

var unixSoc = fmt.Sprintf("unix://%s", file)

func TestPipe(t *testing.T) {
	_ = os.Remove(file)
	defer os.Remove(file)

	h, err := url.Parse(unixSoc)
	if err != nil {
		panic(err)
	}
	nl := nodelogger.New(h)

	w := test.NewJsonWriter()
	go nl.Pipe(w)

	await(t, nl, 1*time.Millisecond)

	packet := `{ "Pippo": "pluto" }`
	if assert.NoError(t, send(packet)) {
		await(t, nl, 1*time.Millisecond)
		p := w.Get()
		assert.Equal(t, "pluto", p.Data["Pippo"])
	}
}

func send(data string) error {
	u, _ := url.Parse(unixSoc)
	c, err := net.Dial(u.Scheme, u.Path)
	if err != nil {
		return err
	}
	if _, err := c.Write([]byte(data)); err != nil {
		return err
	}
	return nil

}

func await(t *testing.T, nl *nodelogger.LogProxy, d time.Duration) {
	//giving enough time to the server to start
	select {
	case err := <-nl.ErrChan:
		assert.FailNow(t, "%s\n", err)
	case <-time.After(d):
		return
	}
}
