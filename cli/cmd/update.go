package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	bumpMajor bool
	bumpMinor bool
	bumpPatch bool
)

var updateCmd = &cobra.Command{
	Use:   "update [commit message]",
	Short: "Update git repo: bump version, commit, and push",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		versionFile := "VERSION"
		oldVersion, err := os.ReadFile(versionFile)
		if err != nil {
			fmt.Println("Failed to read VERSION file:", err)
			return
		}
		parts := strings.Split(strings.TrimSpace(string(oldVersion)), ".")
		if len(parts) != 3 {
			fmt.Println("VERSION file must be in x.y.z format")
			return
		}
		major, err1 := strconv.Atoi(parts[0])
		minor, err2 := strconv.Atoi(parts[1])
		patch, err3 := strconv.Atoi(parts[2])
		if err1 != nil || err2 != nil || err3 != nil {
			fmt.Println("Invalid version number:", parts)
			return
		}

		switch {
		case bumpMajor:
			major++
			minor = 0
			patch = 0
		case bumpMinor:
			minor++
			patch = 0
		default:
			patch++
		}

		newVersion := fmt.Sprintf("%d.%d.%d", major, minor, patch)
		if err := os.WriteFile(versionFile, []byte(newVersion+"\n"), 0644); err != nil {
			fmt.Println("Failed to write VERSION file:", err)
			return
		}
		fmt.Printf("Bumped version to %s\n", newVersion)

		// Build commit message
		var msg string
		if len(args) > 0 {
			msg = strings.Join(args, " ")
		} else {
			msg = fmt.Sprintf("Bump version to %s", newVersion)
		}

		// Run git add, commit, push, showing errors/output
		commands := [][]string{
			{"git", "add", "-A"},
			{"git", "commit", "-m", msg},
			{"git", "push"},
		}
		for _, c := range commands {
			fmt.Printf("> %s\n", strings.Join(c, " "))
			cmd := exec.Command(c[0], c[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Printf("Command failed: %v\n", err)
				return
			}
		}
		fmt.Println("Changes committed and pushed.")
	},
}

func init() {
	updateCmd.Flags().BoolVar(&bumpMajor, "major", false, "Bump major version (x.0.0)")
	updateCmd.Flags().BoolVar(&bumpMinor, "minor", false, "Bump minor version (x.y.0)")
	updateCmd.Flags().BoolVar(&bumpPatch, "patch", false, "Bump patch version (default)")
}
