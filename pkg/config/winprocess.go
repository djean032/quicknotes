package winprocess

import (
	"errors"
	"github.com/shirou/gopsutil/v3/process"
	"log"
	"os"
)

func GetProcessName() (string, error) {
	// Example: Get current process
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		log.Fatalf("Error getting current process: %v", err)
	}

	// Get parent process
	parent, err := p.Parent()
	if err != nil {
		log.Fatalf("Error getting parent process: %v", err)
	}

	if parent == nil {
		return "", errors.New("No parent process found (this might be a system process).")
	}

	name, _ := parent.Name()

	return name, nil
}
