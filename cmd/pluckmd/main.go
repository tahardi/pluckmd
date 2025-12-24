package main

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/tahardi/pluckmd/internal/run"
)

const (
	defaultDocDir  = "/must/provide/a/doc/dir/value"
	defaultTimeout = 60
)

var mainCmd = &cobra.Command{
	Use:          "pluckmd",
	Short:        "CLI tool for downloading and inserting Go code into markdown docs",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		runner, err := run.NewRunner()
		if err != nil {
			return err
		}

		ctx, _ := context.WithTimeout(
			context.Background(),
			time.Duration(timeout)*time.Second,
		)
		return runner.Run(ctx, dir)
	},
}

var dir string
var timeout int

func init() {
	mainCmd.PersistentFlags().StringVarP(
		&dir,
		"dir",
		"d",
		defaultDocDir,
		"directory containing markdown files to process (recursive)",
	)
	mainCmd.PersistentFlags().IntVarP(
		&timeout,
		"timeout",
		"t",
		defaultTimeout,
		"max allowed run time of pluckmd in seconds",
	)
}

func main() {
	if err := mainCmd.Execute(); err != nil {
		fmt.Printf("Command failed: %v\n", err)
	}
}
