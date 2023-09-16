package neviraide

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/RAprogramm/neviraide-install/internal/ui"
	"github.com/RAprogramm/neviraide-install/internal/utils"
)

const (
	repoURL  = "https://github.com/RAprogramm/NEVIRAIDE.git"
	cloneDir = "/tmp/neovim-config"
)

func Install() {
	if _, err := os.Stat(cloneDir); !os.IsNotExist(err) {
		os.RemoveAll(cloneDir)
	}

	fmt.Println(utils.Color("grey", "italic", "Cloning NEVIRAIDE repository..."))

	_, err := exec.Command("git", "clone", "--depth", "1", repoURL, cloneDir).Output()
	if err != nil {
		fmt.Print(utils.Color("red", "italic", "Cloning repository error: %v\n", err))
		return
	}
	fmt.Println(utils.Color("grey", "italic", "Repository cloned successfully!"))

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Print(utils.Color("red", "italic", "Error getting home directory: %v\n", err))
		return
	}
	configDir := filepath.Join(homeDir, ".config/nvim")

	ui.ExistDir(configDir)

	err = exec.Command("cp", "-r", cloneDir, configDir).Run()
	if err != nil {
		fmt.Print(utils.Color("red", "italic", "Error copying configuration files: %v\n", err))
		return
	}

	os.RemoveAll(cloneDir)

	fmt.Println(utils.Color("green", "italic", "NEVIRAIDE was successfully installed!"))
}
