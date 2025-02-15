package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const keptnFolderName = ".keptn"

// GetKeptnDirectory returns a path, which is used to store logs and possibly creds
func GetKeptnDirectory() (string, error) {

	keptnDir := UserHomeDir() + string(os.PathSeparator) + keptnFolderName + string(os.PathSeparator)

	if _, err := os.Stat(keptnDir); os.IsNotExist(err) {
		err := os.MkdirAll(keptnDir, os.ModePerm)
		fmt.Println("keptn creates the folder " + keptnDir + " to store logs and possibly creds.")
		if err != nil {
			return "", err
		}
	}

	return keptnDir, nil
}

// UserHomeDir returns the HOME directory by taking into account the operating system
func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

// ExpandTilde expands ~ to HOME
func ExpandTilde(fileName string) string {
	if fileName == "~" {
		return UserHomeDir()
	} else if strings.HasPrefix(fileName, "~/") {
		return filepath.Join(UserHomeDir(), fileName[2:])
	}
	return fileName
}

// GetFiles returns a list of files in a directory filtered by the provided suffix
func GetFiles(workingPath string, suffixes ...string) ([]string, error) {
	var files []string
	err := filepath.Walk(workingPath, func(path string, info os.FileInfo, err error) error {
		for _, suffix := range suffixes {
			if strings.HasSuffix(path, suffix) {
				files = append(files, path)
				break
			}
		}
		return nil
	})
	return files, err
}
