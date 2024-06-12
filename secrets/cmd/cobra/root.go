package cobra

import (
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"path/filepath"
)

var RootCmd = &cobra.Command{
	Use:   "secret",
	Short: "Secret is an API key and other secrets manager",
}

var encodingKey string

func init() {
	RootCmd.PersistentFlags().StringVarP(&encodingKey, "key", "k",
		"", "the encoding and decoding key")

}

func secretsPath() string {
	home, _ := homedir.Dir()
	return filepath.Join(home, ".secrets")
}
