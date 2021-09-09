package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/robbiemcmichael/processpool"
)

var options processpool.Options

var rootCmd = &cobra.Command{
	Use:   "processpool [flags] -- command [args]",
	Short: "Maintains a pool of processes which can be fetched via HTTP.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("mussing")
		}
		run(args[0], args[1:])
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&options.Address, "address", "a", "127.0.0.1:8080", "address for processpool to listen on")
	rootCmd.PersistentFlags().IntVarP(&options.Port, "port", "p", 0, "starting port for the first process to use (required)")
	rootCmd.PersistentFlags().IntVarP(&options.Processes, "processes", "n", 0, "number of processes (required)")

	rootCmd.MarkPersistentFlagRequired("port")
	rootCmd.MarkPersistentFlagRequired("processes")
}

func run(path string, args []string) {
	pool := processpool.NewPool(options, path, args)

	if err := pool.Serve(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
