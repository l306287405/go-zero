package prometheus

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/l306287405/go-zero/core/logx"
	"github.com/l306287405/go-zero/core/syncx"
	"github.com/l306287405/go-zero/core/threading"
)

var (
	once    sync.Once
	enabled syncx.AtomicBool
)

// Enabled returns if prometheus is enabled.
func Enabled() bool {
	return enabled.True()
}

// StartAgent starts a prometheus agent.
func StartAgent(c Config) {
	if len(c.Host) == 0 {
		return
	}

	once.Do(func() {
		enabled.Set(true)
		threading.GoSafe(func() {
			http.Handle(c.Path, promhttp.Handler())
			addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
			logx.Infof("Starting prometheus agent at %s", addr)
			if err := http.ListenAndServe(addr, nil); err != nil {
				logx.Error(err)
			}
		})
	})
}
