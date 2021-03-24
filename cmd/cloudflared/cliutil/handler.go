package cliutil

import (
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/cloudflare/cloudflared/config"
	"github.com/cloudflare/cloudflared/logger"
)

func Action(actionFunc cli.ActionFunc) cli.ActionFunc {
	return WithErrorHandler(actionFunc)
}

func ConfiguredAction(actionFunc cli.ActionFunc) cli.ActionFunc {
	return WithErrorHandler(func(c *cli.Context) error {
		if err := setFlagsFromConfigFile(c); err != nil {
			return err
		}
		return actionFunc(c)
	})
}

func setFlagsFromConfigFile(c *cli.Context) error {
	const errorExitCode = 1
	log := logger.CreateLoggerFromContext(c, logger.EnableTerminalLog)
	inputSource, err := config.ReadConfigFile(c, log)
	if err != nil {
		if err == config.ErrNoConfigFile {
			return nil
		}
		return cli.Exit(err, errorExitCode)
	}

	if err := altsrc.ApplyInputSource(c, inputSource); err != nil {
		return cli.Exit(err, errorExitCode)
	}
	return nil
}
