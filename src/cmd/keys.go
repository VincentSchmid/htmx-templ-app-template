//go:build dev

package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"
)

var keysCmd = &cobra.Command{
	Use:   "generate-keys",
	Short: "Generate new session keys",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		authKey := generateKey(64)
		encryptKey := generateKey(32)

		fmt.Println("Auth Key:", authKey)
		fmt.Println("Encrypt Key:", encryptKey)
	},
}

func generateKey(length int) string {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(key)
}

func init() {
	rootCmd.AddCommand(keysCmd)
}
