package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var mailInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("info called")
	},
	Annotations: map[string]string{"position": "20"},
}

func init() {
	mailCmd.AddCommand(mailInfoCmd)
}
