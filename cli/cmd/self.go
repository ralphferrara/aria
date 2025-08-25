package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

var selfCmd = &cobra.Command{
	Use:   "self-rebuild",
	Short: "Rebuild the aria CLI binary from source (stays in aria-cli only)",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, _ := os.Getwd()

		var buildName string
		if runtime.GOOS == "windows" {
			buildName = "aria_temp.exe"
		} else {
			buildName = "aria_temp"
		}

		buildPath := filepath.Join(cwd, buildName)

		// Build the new binary in the current directory
		buildCmd := exec.Command("go", "build", "-o", buildPath, "-buildvcs=false")
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr
		fmt.Println("Building new binary at:", buildPath)
		if err := buildCmd.Run(); err != nil {
			fmt.Println("Build failed:", err)
			return
		}
		fmt.Println("Build successful.")

		// Remove temp file after build
		if err := os.Remove(buildPath); err != nil {
			fmt.Printf("Warning: failed to remove temp file %s: %v\n", buildPath, err)
		} else {
			fmt.Printf("Temp file %s removed.\n", buildPath)
		}
	},
}
