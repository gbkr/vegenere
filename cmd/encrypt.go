package cmd

import (
	"fmt"
	"vegenere/vegenerelib"

	"github.com/spf13/cobra"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt plaintext",
	Long:  `Encrypt plaintext`,
	Run: func(cmd *cobra.Command, args []string) {
		result := vegenerelib.Encrypt(Source, EncryptionKey)
		fmt.Println(result)
	},
}

func init() {
	RootCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().StringVarP(&EncryptionKey, "key", "k", "", "Key to use for encryption / decryption")
}
