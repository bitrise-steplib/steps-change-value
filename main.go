package main

import (
	"fmt"
	"log"
	"os"

	"strings"

	"github.com/bitrise-io/go-utils/fileutil"
)

func main() {
	var (
		inputFile                     = os.Getenv("file")
		inputOldValue                 = os.Getenv("old_value")
		inputNewValue                 = os.Getenv("new_value")
		inputIsShowFileContent        = os.Getenv("show_file") == "true"
		inputIsFailIfOldValueNotFound = os.Getenv("notfound_exit") == "false"
	)

	if inputFile == "" {
		log.Fatal("No file input specified")
	}
	if inputOldValue == "" {
		log.Fatal("No old_value input specified")
	}
	if inputNewValue == "" {
		log.Fatal("No new_value input specified")
	}

	origContent, err := fileutil.ReadStringFromFile(inputFile)
	if err != nil {
		log.Fatalf("Failed to read from specified file, error: %s", err)
	}

	if inputIsShowFileContent {
		fmt.Println()
		fmt.Println("------------------------------------------")
		fmt.Println("-------------OLD  FILE--------------------")
		fmt.Println("------------------------------------------")
		fmt.Print(origContent)
		fmt.Println()
		fmt.Println("------------------------------------------")
	}

	if inputIsFailIfOldValueNotFound {
		if !strings.Contains(origContent, inputOldValue) {
			log.Fatalf("Provided old value (%s) was not found in the file (%s)", inputOldValue, inputFile)
		}
	}

	// replace
	fmt.Println(" (i) Replacing...")
	replacedContent := strings.Replace(origContent, inputOldValue, inputNewValue, -1)

	if inputIsShowFileContent {
		fmt.Println()
		fmt.Println("------------------------------------------")
		fmt.Println("-------------NEW  FILE--------------------")
		fmt.Println("------------------------------------------")
		fmt.Print(replacedContent)
		fmt.Println()
		fmt.Println("------------------------------------------")
	}

	// write back to file
	if err := fileutil.WriteStringToFile(inputFile, replacedContent); err != nil {
		log.Printf("Failed to write replaced content back to file, error: %s", err)
	}
	fmt.Println(" (i) Done")
}
