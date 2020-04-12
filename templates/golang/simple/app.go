package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli"
)

func NewApp() *cli.App {
	config := NewDefaultConfig()
	app := cli.NewApp()
	app.Name = prog
	app.Usage = description
	app.Version = version
	app.Flags = getCommandLineOptions()
	app.UsageText = fmt.Sprintf("%s [options]", prog)

	app.Action = func(ctx *cli.Context) error {
		if err := parseCLIOptions(ctx, config); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		server, err := NewAPIServer(config)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		defer server.logger.Sync()
		if err := server.Run(); err != nil {
			return cli.NewExitError(err.Error(), 1)

		}

		// Setup the termination signals
		signalChannel := make(chan os.Signal)
		signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-signalChannel

		return nil
	}

	return app
}
