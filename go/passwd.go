package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type User struct {
	Username string
	// Passwd string
	Uid int
	Gid int
	// Comment string
	HomeDir string
	Shell   string
}

func LookupUserFromPasswd(id int) User {
	passwdBytes, err := os.ReadFile("/etc/passwd")
	check(
		err,
		"Failed reading /etc/passwd from disk",
	)
	lines := strings.Split(string(passwdBytes), "\n")
	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) != 7 {
			panic(fmt.Sprintf("Unexpected number of fields in %s", line))
		}
		uid, err := strconv.Atoi(fields[2])
		check(
			err,
			fmt.Sprintf("Invalid int value for uid: %s\n\n%v", fields[2], err),
		)
		if uid != id {
			continue
		}
		gid, err := strconv.Atoi(fields[3])
		check(
			err,
			fmt.Sprintf("Invalid int value for gid: %s\n\n%v", fields[3], err),
		)
		return User{
			Username: fields[0],
			// Passwd: fields[1]
			Uid: uid,
			Gid: gid,
			// Comment: fields[4]
			HomeDir: fields[5],
			Shell:   fields[6],
		}
	}
	panic(
		fmt.Sprintf("Did not find the user id %d in /etc/passwd:\n\n%s", id, string(passwdBytes)),
	)
}
