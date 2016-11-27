package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Source string
var EncryptionKey string

var RootCmd = &cobra.Command{
	Use:   "vegenere",
	Short: "Encrypt, decrypt and break Vegenere ciphers",
	Long:  `Vegenere is an application for encrypting, decrypting and breaking Vegenere ciphers.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&Source, "source", "s", "", "Source file to be processed")
}
