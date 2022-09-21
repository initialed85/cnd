package app

import (
	"database/sql"
	"fmt"
	_ "github.com/glebarez/go-sqlite"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	Usage = "usage: cnd [checkCommand | doCommand | run] [command]"
)

type App struct {
	jobName      string
	checkCommand string
	doCommand    string
	db           *sql.DB
}

func New(
	jobName string,
	checkCommand string,
	doCommand string,
) (*App, error) {
	if strings.TrimSpace(jobName) == "" {
		return nil, fmt.Errorf("jobName empty or missing")
	}

	if strings.TrimSpace(checkCommand) == "" {
		return nil, fmt.Errorf("checkCommand empty or missing")
	}

	homePath := os.Getenv("HOME")
	if strings.TrimSpace(homePath) == "" {
		return nil, fmt.Errorf("HOME env var empty or unset")
	}

	stateFolderPath := filepath.Join(homePath, ".cnd")

	stat, err := os.Stat(stateFolderPath)
	if err == nil && !stat.IsDir() {
		return nil, fmt.Errorf("statePath=%#+v already exists but is not a folder", stateFolderPath)
	}

	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(stateFolderPath, 0755)
		if err != nil {
			return nil, err
		}
	}

	stateFilePath := filepath.Join(stateFolderPath, fmt.Sprintf("state-%v.db", jobName))

	db, err := sql.Open("sqlite", stateFilePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(createSchema)
	if err != nil {
		return nil, err
	}

	a := App{
		jobName:      jobName,
		checkCommand: checkCommand,
		doCommand:    doCommand,
		db:           db,
	}

	return &a, nil
}

func (a *App) Check() (bool, error) {
	timestamp := time.Now().Unix()

	log.Printf("check; command=%#+v", a.checkCommand)

	checkCommandProcess := exec.Command(
		"bash",
		"-c",
		a.checkCommand,
	)

	rawCheckCommandOutput, err := checkCommandProcess.CombinedOutput()
	if err != nil {
		return false, err
	}

	checkCommandOutput := strings.TrimSpace(string(rawCheckCommandOutput))

	log.Printf("check; output=%#+v", checkCommandOutput)

	row := a.db.QueryRow(getLastHistory, a.checkCommand)

	var lastRowId int64
	var lastJobName string
	var lastRawTimestamp int64
	var lastCheckCommand string
	var lastCheckCommandOutput string

	err = row.Scan(
		&lastRowId,
		&lastJobName,
		&lastRawTimestamp,
		&lastCheckCommand,
		&lastCheckCommandOutput,
	)

	history := History{
		RowId:              lastRowId,
		jobName:            lastJobName,
		Timestamp:          time.Unix(lastRawTimestamp, 0),
		CheckCommand:       lastCheckCommand,
		CheckCommandOutput: lastCheckCommandOutput,
	}

	noRows := err == sql.ErrNoRows

	if err != nil && !noRows {
		return false, err
	}

	if !noRows {
		if history.CheckCommandOutput == checkCommandOutput {
			log.Printf("check; lastOutput=%#+v", history.CheckCommandOutput)
			log.Printf("check; changed=false")
			return false, nil
		} else {
			log.Printf("check; lastOutput=%#+v", history.CheckCommandOutput)
		}
	}

	if a.doCommand != "" {
		_, err = a.db.Exec(addToHistory, a.jobName, timestamp, a.checkCommand, checkCommandOutput)
		if err != nil {
			return false, err
		}

		_, err = a.db.Exec(pruneHistory)
		if err != nil {
			return false, err
		}
	}

	log.Printf("check; changed=true")
	return true, nil
}

func (a *App) Do() error {
	log.Printf("do; command=%#+v", a.doCommand)

	doCommandProcess := exec.Command(
		"bash",
		"-c",
		a.doCommand,
	)

	rawCheckCommandOutput, err := doCommandProcess.CombinedOutput()
	if err != nil {
		return err
	}

	doCommandOutput := strings.TrimSpace(string(rawCheckCommandOutput))

	log.Printf("do; output=%#+v", doCommandOutput)

	return nil
}

func (a *App) Run() error {
	log.Printf("check-n-do; jobName=%#+v", a.jobName)

	changed, err := a.Check()
	if err != nil {
		return err
	}

	if !changed {
		return nil
	}

	if a.doCommand == "" {
		log.Printf("check-n-do; done=false (dry-run check)")
		return nil
	}

	err = a.Do()
	if err != nil {
		return err
	}

	log.Printf("check-n-do; done=true")

	return nil
}
