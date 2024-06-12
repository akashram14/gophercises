package cmd

import (
	"cmd/task/db"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a task to your task list.",
	Run: func(cmd *cobra.Command, args []string) {
		taskvalue := strings.Join(args, " ")
		err := db.CreateTask(taskvalue)
		if err != nil {
			fmt.Println("Something went wrong:", err.Error())
			os.Exit(1)
		}
		fmt.Printf("Added \"%s\" to your task list.\n", taskvalue)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
