package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

const logo = `
░░░    ░░  ░░░░░░  ░░░░░░░░ ░░ ░░░░░░   ░░░░░░  ░░░░░░░░ 
▒▒▒▒   ▒▒ ▒▒    ▒▒    ▒▒    ▒▒ ▒▒   ▒▒ ▒▒    ▒▒    ▒▒    
▒▒ ▒▒  ▒▒ ▒▒    ▒▒    ▒▒    ▒▒ ▒▒▒▒▒▒  ▒▒    ▒▒    ▒▒    
▓▓  ▓▓ ▓▓ ▓▓    ▓▓    ▓▓    ▓▓ ▓▓   ▓▓ ▓▓    ▓▓    ▓▓    
██   ████  ██████     ██    ██ ██████   ██████     ██    

Version: %s
Build Date: %s
Commit Hash: %s
`

var (
	build   string
	commit  string
	version string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show Version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(logo, version, build, commit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
