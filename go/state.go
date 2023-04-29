package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const name = "dev-chroot"

type State struct {
	WorkspaceRoot string
}

func InitState(initialPath string, home string) map[string]State {
	return map[string]State{
		initialPath: {
			WorkspaceRoot: getShare(home),
		},
	}
}

func getShare(home string) string {
	// TODO first check $XDG_DATA_HOME
	// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
	share := filepath.Join(home, ".local", "share")
	// TODO check share exists
	if _, err := os.ReadDir(share); err != nil {
		panic(fmt.Errorf("Expected %s to exist but it did not", share))
	}
	devchrootShare := filepath.Join(share, "."+name)
	if _, err := os.Stat(devchrootShare); err != nil {
		// rw-r--r--
		check(
			os.Mkdir(devchrootShare, 0644),
			fmt.Sprintf("Failed to mkdir %s", devchrootShare),
		)
	}
	return devchrootShare
}

func GetHome() string {
	home, isSet := os.LookupEnv("HOME")
	if !isSet {
		panic("There is no \"HOME\" env var set!")
	}
	return home
}

func GetStateFilePath(home string) string {
	return filepath.Join(home, "."+name+".json")
}

// Read and return State from disk.
//
// If the state file did not exist, will return an error.
func HydrateState(home string) (map[string]State, error) {
	filePath := GetStateFilePath(home)
	stateFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(stateFile)
	states := map[string]State{}
	// linter seems to want me to pass the address of states?
	if err := decoder.Decode(&states); err != nil && err != io.EOF {
		panic(err)
	}
	return states, nil
}

func PersistStates(states map[string]State, home string) {
	// default mode 0666
	file, err := os.Create(GetStateFilePath(home))
	if err != nil {
		panic(err)
	}
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(states); err != nil {
		panic(err)
	}
}
