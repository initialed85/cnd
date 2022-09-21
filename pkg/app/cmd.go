package app

import (
	"flag"
	"os"
	"path/filepath"
)

func Run() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	_, cwdFolderName := filepath.Split(cwd)

	jobName := flag.String(
		"job",
		cwdFolderName,
		"Name of the check-n-do job",
	)

	checkCommand := flag.String(
		"check",
		"",
		"Check command to execute for gathering the comparison state",
	)

	doCommand := flag.String(
		"do",
		"",
		"Do command to execute if this comparison state doesn't match last comparison state (optional; unset = dry-run check)",
	)

	flag.Parse()

	a, err := New(*jobName, *checkCommand, *doCommand)
	if err != nil {
		return err
	}

	err = a.Run()
	if err != nil {
		return err
	}

	return nil
}
