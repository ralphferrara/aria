package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run tests (auto-starts docker-compose)",
	Run: func(cmd *cobra.Command, args []string) {
		composeFile := "docker-compose.yml"
		fmt.Println("Starting docker-compose for tests...")
		up := exec.Command("docker-compose", "-f", composeFile, "up", "-d", "--remove-orphans")
		up.Stdout = os.Stdout
		up.Stderr = os.Stderr
		if err := up.Run(); err != nil {
			fmt.Println("docker-compose up failed:", err)
			return
		}

		// Wait a few seconds for services to initialize (customize as needed)
		fmt.Print("Waiting for services to start")
		for i := 0; i < 6; i++ {
			time.Sleep(1 * time.Second)
			fmt.Print(".")
		}
		fmt.Println()

		fmt.Println("Running tests...")
		test := exec.Command("go", "test", "./...")
		test.Stdout = os.Stdout
		test.Stderr = os.Stderr
		err := test.Run()
		if err != nil {
			fmt.Println("Tests failed:", err)
		} else {
			fmt.Println("All tests passed.")
		}

		fmt.Println("Stopping docker-compose...")
		down := exec.Command("docker-compose", "-f", composeFile, "down", "-v", "--remove-orphans")
		down.Stdout = os.Stdout
		down.Stderr = os.Stderr
		if err := down.Run(); err != nil {
			fmt.Println("docker-compose down failed:", err)
		} else {
			fmt.Println("Test environment cleaned up.")
		}
	},
}
