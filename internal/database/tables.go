package database

import (
	"os/exec"
)

func ExecuteSQLFile(dbPath, sqlFilePath string) error {
	cmd := exec.Command("sqlite3", dbPath, ".read", sqlFilePath)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}
