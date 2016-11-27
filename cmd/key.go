package cmd

import (
	"fmt"
	"vegenere/vegenerelib"

	"github.com/spf13/cobra"
)

var keyCmd = &cobra.Command{
	Use:   "key",
	Short: "Display the key",
	Long:  `Use this command to attempt to break the cipher and display the key.`,
	Run: func(cmd *cobra.Command, args []string) {
		result := vegenerelib.DecryptKey(Source)
		fmt.Println(result)
	},
}

func init() {
	decryptCmd.AddCommand(keyCmd)
}
