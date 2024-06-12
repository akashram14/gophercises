package cobra

import (
	"github.com/spf13/cobra"
	"secrets"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get commands",
	Run: func(cmd *cobra.Command, args []string) {
		v := secrets.File(encodingKey, secretsPath())
		v.Get(args[0])
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}
