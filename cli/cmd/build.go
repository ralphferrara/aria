package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the Aria app",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Building the project...")
		out, err := exec.Command("go", "build", "./...").CombinedOutput()
		fmt.Print(string(out))
		if err != nil {
			fmt.Println("Build failed:", err)
		} else {
			fmt.Println("Build succeeded.")
		}
	},
}
