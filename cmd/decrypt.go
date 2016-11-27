package cmd

import (
	"fmt"
	"vegenere/vegenerelib"

	"github.com/spf13/cobra"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt a given cipher",
	Long: `Given a key, this command will decrypt a Vegenere cipher. Called without a key,
it will attempt to break the cipher and display the plaintext. Used with the key command
and supplied a source file it will attempt to break the cipher and display the key.`,
	Run: func(cmd *cobra.Command, args []string) {
		result := vegenerelib.Decrypt(Source, EncryptionKey)
		fmt.Println(result)
	},
}

func init() {
	RootCmd.AddCommand(decryptCmd)
	decryptCmd.Flags().StringVarP(&EncryptionKey, "key", "k", "", "Key to use for encryption / decryption")
}
