package foundation

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bootf/foundation/httpserver"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Foundation struct {
	Version string
	Command struct {
		Commands []*cobra.Command
	}
	Http struct {
		Port    string
		Handler http.Handler
	}
	shutdowns []func() error
}

func (f *Foundation) CommandRun() {
	root := &cobra.Command{
		Version: f.Version,
	}

	root.AddCommand(f.Command.Commands...)

	if err := root.Execute(); err != nil {
		logrus.Fatalf("unable to execute command error : %s", err)
	}
}

func (f *Foundation) ServerRun() {
	httpServer := httpserver.New(f.Http.Handler, httpserver.Port(f.Http.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logrus.Info("signal: " + s.String())
	case err := <-httpServer.Notify():
		logrus.Errorf("http server.Notify: %s", err)
	}

	// Shutdown
	err := httpServer.Shutdown()
	if err != nil {
		logrus.Errorf("http server.shutdown: %s", err)
	}

	// another shutdown
	for i := range f.shutdowns {
		f.shutdowns[i]()
	}
}
