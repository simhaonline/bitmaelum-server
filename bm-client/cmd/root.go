package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "bm-client",
	Short: "BitMaelum client",
	Long: `This client allows you to manage accounts, read and compose mail.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringArrayP("config", "c", []string{}, "configuration file")
}
