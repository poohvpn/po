package pooh

import (
	"fmt"
	"os/exec"
	"strings"
)

func Run(name string, args ...string) (string, error) {
	outBytes, err := exec.Command(name, args...).CombinedOutput()
	outString := strings.TrimSpace(string(outBytes))
	if err != nil {
		if outString == "" {
			outString = err.Error()
		}
		cmd := strings.Join(append([]string{name}, args...), " ")
		return "", fmt.Errorf("%s: %s", cmd, outString)
	}
	return outString, nil
}
