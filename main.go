package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// greeting displays a welcome message and asks the user for confirmation
// to continue with the installation.
func greeting() bool {
	fmt.Println("Welcome to the NEVIRAIDE installer!")
	fmt.Println("This script will check for required dependencies and install them if they're missing.")
	fmt.Println("It will also set up the NEVIRAIDE configuration for Neovim.")
	fmt.Println("")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Would you like to continue with the installation? [y/n]: ")
	answer, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}
	answer = strings.TrimSpace(answer)
	if answer != "y" {
		fmt.Println("Installation canceled by the user.")
		return false
	}
	return true
}

// checkCommandsAvailability checks if the required commands are available
// Returns a slice of missing commands.
func checkCommandsAvailability(names map[string]string) []string {
	missing := []string{}
	for name, cmd := range names {
		cmd := exec.Command("command", "-v", cmd)
		if err := cmd.Run(); err != nil {
			missing = append(missing, name)
		}
	}
	return missing
}

// checkSudo checks if sudo is available
func checkSudo() bool {
	cmd := exec.Command("command", "-v", "sudo")
	return cmd.Run() == nil
}

// installWithPacman installs the provided package using pacman.
func installWithPacman(pkg string) bool {
	fmt.Printf("Installing %s...\\n", pkg)
	cmd := exec.Command("sudo", "pacman", "-S", "--noconfirm", pkg)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Installation error for %s: %v\\n", pkg, err)
		return false
	}
	fmt.Printf("%s was installed successfully!\\n", pkg)
	return true
}

func main() {
	if !greeting() {
		os.Exit(0)
	}

	if !checkSudo() {
		fmt.Println("sudo is not available or you don't have sudo privileges.")
		os.Exit(1)
	}

	dependencies := map[string]string{
		"neovim":  "nvim",
		"git":     "git",
		"ripgrep": "rg",
		"fd":      "fd",
		"unzip":   "unzip",
		"tar":     "tar",
		"wget":    "wget",
		"curl":    "curl",
		"npm":     "npm",
	}

	missingDeps := checkCommandsAvailability(dependencies)
	for _, dep := range missingDeps {
		success := installWithPacman(dep)
		if !success {
			fmt.Printf("Failed to install %s. Aborting installation.\\n", dep)
			os.Exit(1)
		}
	}

	repoURL := "https://github.com/RAprogramm/NEVIRAIDE.git"
	cloneDir := "/tmp/neovim-config"

	fmt.Println("Cloning NEVIRAIDE repository...")
	_, err := exec.Command("git", "clone", "--depth", "1", repoURL, cloneDir).Output()
	if err != nil {
		fmt.Printf("Cloning repository error: %v\\n", err)
		return
	}
	fmt.Println("Repository cloned successfully!")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\\n", err)
		return
	}
	configDir := filepath.Join(homeDir, ".config/nvim")

	if _, err = os.Stat(configDir); !os.IsNotExist(err) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("~/.config/nvim already exists. Remove it? [y/n]: ")
		answer, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			os.Exit(1)
		}
		answer = strings.TrimSpace(answer)

		switch answer {
		case "y":
			err = os.RemoveAll(configDir)
			if err != nil {
				fmt.Printf("Error removing directory: %v\\n", err)
			}
		case "n":
			err = os.Rename(configDir, configDir+".old")
			if err != nil {
				fmt.Printf("Error renaming directory: %v\\n", err)
			}
		default:
			fmt.Println("Undefined choice. Abort installation.")
			os.Exit(0)
		}
	}

	// Ensure the destination directory exists
	os.MkdirAll(configDir, os.ModePerm)

	err = exec.Command("cp", "-r", cloneDir, configDir).Run()
	if err != nil {
		fmt.Printf("Error copying configuration files: %v\\n", err)
		return
	}

	// Cleaning up the temporary directory
	os.RemoveAll(cloneDir)

	fmt.Println("NEVIRAIDE was successfully installed!")
}
