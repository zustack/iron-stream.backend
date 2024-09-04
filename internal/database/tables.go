package database

import (
	"fmt"
	"os/exec"
)

func ExecuteSQLFile(dbPath, sqlFilePath string) error {
	cmd := exec.Command("sqlite3", dbPath, fmt.Sprintf(".read %s", sqlFilePath))
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error exec : %v\nExit: %s\n", err, output)
		return err
	}
	return nil
}
