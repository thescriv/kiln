package utilconfig

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// solution based on : https://github.com/joho/godotenv/issues/126

// Load loads the environment variables from the .env file.
func LoadConfig() {
	err := godotenv.Load(dir(".env"))
	if err != nil {
		panic(fmt.Errorf("error loading .env file: %w", err))
	}
}

// dir returns the absolute path of the given environment file (envFile) in the Go module's
// root directory. It searches for the 'go.mod' file from the current working directory upwards
// and appends the envFile to the directory containing 'go.mod'.
// It panics if it fails to find the 'go.mod' file.
func dir(envFile string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			break
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			panic(fmt.Errorf("go.mod not found"))
		}
		currentDir = parent
	}

	return filepath.Join(currentDir, envFile)
}
