package root

import (
	"fmt"
	"os"

	"github.com/SECRYPT-2022/SECRYPT/command/backup"
	"github.com/SECRYPT-2022/SECRYPT/command/genesis"
	"github.com/SECRYPT-2022/SECRYPT/command/helper"
	"github.com/SECRYPT-2022/SECRYPT/command/ibft"
	"github.com/SECRYPT-2022/SECRYPT/command/license"
	"github.com/SECRYPT-2022/SECRYPT/command/monitor"
	"github.com/SECRYPT-2022/SECRYPT/command/peers"
	"github.com/SECRYPT-2022/SECRYPT/command/secrets"
	"github.com/SECRYPT-2022/SECRYPT/command/server"
	"github.com/SECRYPT-2022/SECRYPT/command/status"
	"github.com/SECRYPT-2022/SECRYPT/command/txpool"
	"github.com/SECRYPT-2022/SECRYPT/command/version"
	"github.com/SECRYPT-2022/SECRYPT/command/whitelist"
	"github.com/spf13/cobra"
)

type RootCommand struct {
	baseCmd *cobra.Command
}

func NewRootCommand() *RootCommand {
	rootCommand := &RootCommand{
		baseCmd: &cobra.Command{
			Short: "Polygon Edge is a framework for building Ethereum-compatible Blockchain networks",
		},
	}

	helper.RegisterJSONOutputFlag(rootCommand.baseCmd)

	rootCommand.registerSubCommands()

	return rootCommand
}

func (rc *RootCommand) registerSubCommands() {
	rc.baseCmd.AddCommand(
		version.GetCommand(),
		txpool.GetCommand(),
		status.GetCommand(),
		secrets.GetCommand(),
		peers.GetCommand(),
		monitor.GetCommand(),
		ibft.GetCommand(),
		backup.GetCommand(),
		genesis.GetCommand(),
		server.GetCommand(),
		whitelist.GetCommand(),
		license.GetCommand(),
	)
}

func (rc *RootCommand) Execute() {
	if err := rc.baseCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}
