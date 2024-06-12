package cmd

import (
	"cmd/task/db"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completedCmd = &cobra.Command{
	Use:   "completed",
	Short: "Lists all completed tasks within 24hrs",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := db.AllCompletedTasks()
		if err != nil {
			fmt.Println("Something went wrong", err)
			os.Exit(1)
		}
		if len(tasks) == 0 {
			fmt.Println("No tasks")
			return
		}
		fmt.Println("You have the following tasks")
		for i, task := range tasks {
			fmt.Printf("%d. %s Completed on %v, Key = %d\n", i+1, task.Value, task.TimeCompleted, task.Key)
		}
	},
}

func init() {
	RootCmd.AddCommand(completedCmd)
}
