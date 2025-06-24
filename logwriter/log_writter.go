package logwriter

import (
	"encoding/json"
	"fmt"
	"os"
)

func (lw *LogWritter) Write() error {
	b, err := json.Marshal(lw.counterResult)
	if err != nil {
		return err
	}

	err = os.WriteFile(
		fmt.Sprintf("%s.json", lw.configuration.EnvironmentVariables.GetLogPath()),
		[]byte(b),
		os.FileMode(0644),
	)

	if err != nil {
		return err
	}

	return nil
}
