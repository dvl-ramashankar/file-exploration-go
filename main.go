package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	op = `What operation do you want to perform? Please choose given following options (like A or a):
		A. List all files
		B. Search file
		C. Copy or move file
		D. Delete file
		Q. Exit`
	SourcePath      = "Enter source file path:"
	DestinationPath = "Enter destination file path:"
	FileName        = "Enter file name:"
)

func main() {
	for {
		var sourceDirectory, destinationDirectory, fileName, operation string
		fmt.Println(op)
		fmt.Scanln(&operation)
		operation = strings.ToLower(operation)
		switch operation {
		case "a":
			allFileOperation(sourceDirectory)
		case "b":
			searchOperation(sourceDirectory, fileName)
		case "c":
			copyAndMoveOperation(sourceDirectory, destinationDirectory, fileName)
		case "d":
			deleteOperation(sourceDirectory, fileName)
		case "q":
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("invalid option")
		}
	}
}

func deleteOperation(sourceDirectory, fileName string) {
	fmt.Println(SourcePath)
	fmt.Scanln(&sourceDirectory)
	fmt.Println(FileName)
	fmt.Scanln(&fileName)
	df, err := deleteFile(sourceDirectory, fileName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(df)
}

func copyAndMoveOperation(sourceDirectory, destinationDirectory, fileName string) {
	fmt.Println(SourcePath)
	fmt.Scanln(&sourceDirectory)
	fmt.Println(DestinationPath)
	fmt.Scanln(&destinationDirectory)
	fmt.Println(FileName)
	fmt.Scanln(&fileName)
	cpOrMV := ""
	fmt.Println("Which operation do you want to perform copy or move:")
	fmt.Scanln(&cpOrMV)
	cp, err := copyAndMoveFile(sourceDirectory, destinationDirectory, fileName, cpOrMV)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cp)
}

func searchOperation(sourceDirectory, fileName string) {
	fmt.Println(SourcePath)
	fmt.Scanln(&sourceDirectory)
	fmt.Println(FileName)
	fmt.Scanln(&fileName)
	flag, err := searchFiles(sourceDirectory, fileName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(flag)
}

func allFileOperation(sourceDirectory string) {
	fmt.Println(SourcePath)
	fmt.Scanln(&sourceDirectory)
	files, err := listAllFiles(sourceDirectory)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(files)
}

// List all files and nested also present in given path
func listAllFiles(dirPath string) ([]string, error) {
	var files []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("directory does not exist")
		}

		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

// search file name in the given directory
func searchFiles(dirPath, fileName string) (string, error) {
	flag := "file is not present in the given directory"

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("directory does not exist")
		}

		if !info.IsDir() && info.Name() == fileName {
			flag = "file is present in the given directory"
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return flag, nil
}

// copy and move files from one directory to given directory
func copyAndMoveFile(sourceDirectory, destinationDirectory, fileName, operation string) (string, error) {
	var flag string
	var err error

	switch operation {
	case "copy":
		flag, err = copyFile(sourceDirectory, destinationDirectory, fileName)
	case "move":
		flag, err = moveFile(sourceDirectory, destinationDirectory, fileName)
	default:
		return "", fmt.Errorf("invalid operation %s", operation)
	}
	if err != nil {
		return "", err
	}
	return flag, nil
}

// copy file
func copyFile(sourceDirectory, destinationDirectory, fileName string) (string, error) {
	if fileName == "" {
		err := copyDirectory(sourceDirectory, destinationDirectory, "copy")
		if err != nil {
			return "", err
		}
		return "Directory copied successfully", nil
	}
	sourceFile, err := os.Open(sourceDirectory + "/" + fileName)
	if err != nil {
		return "", fmt.Errorf("failed to open %s file in given directory %s", fileName, sourceDirectory)
	}
	defer sourceFile.Close()
	err = createFolderIfNotExist(destinationDirectory)
	if err != nil {
		return "", err
	}
	destinationFile, err := os.Create(destinationDirectory + "/" + fileName)
	if err != nil {
		return "", err
	}
	defer destinationFile.Close()
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return "", err
	}

	return "File copied successfully!", nil
}

// Copy all files and directories in the given directory to the destination
func copyDirectory(sourceDirectory, destinationDirectory, operation string) error {
	err := filepath.Walk(sourceDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(sourceDirectory, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(destinationDirectory, relPath)
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}
		if err := copyFileDirectory(path, destPath); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	if operation == "move" {
		if err := os.RemoveAll(sourceDirectory); err != nil {
			return err
		}
	}
	return nil
}

func copyFileDirectory(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}
	return nil
}

// move file
func moveFile(sourceDirectory, destinationDirectory, fileName string) (string, error) {
	if fileName == "" {
		err := copyDirectory(sourceDirectory, destinationDirectory, "move")
		if err != nil {
			return "", err
		}
		return "Directory moved successfully", nil
	}
	_, err := copyFile(sourceDirectory, destinationDirectory, fileName)
	if err != nil {
		return "", err
	}
	err = os.Remove(sourceDirectory + "/" + fileName)
	if err != nil {
		return "", err
	}
	return "File moved successfully!", nil
}

// create folder if it doesn't exist
func createFolderIfNotExist(folderPath string) error {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// Delete file from given directory
func deleteFile(filePath, fileName string) (string, error) {
	flag := "File deleted successfully"
	if fileName == "" {
		err := deleteDirectory(filePath)
		if err != nil {
			return "", err
		}
		return "Directory deleted successfully", nil
	}
	ofile, err := os.Open(filePath + "/" + fileName)
	if err != nil {
		//	log.Printf("failed to open %s file in given directory %s: %v", fileName, filePath, err)
		return "", fmt.Errorf("file is not present in the given directory")
	}
	ofile.Close()
	err = os.Remove(filePath + "/" + fileName)
	if err != nil {
		return "", err
	}
	return flag, nil
}

// Delete a folder and its children
func deleteDirectory(filePath string) error {
	df, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("directory not found")
	}
	df.Close()
	err = os.RemoveAll(filePath)
	if err != nil {
		return err
	}
	return nil
}
