package pooh

import (
	"fmt"
	"os/exec"
	"strings"
)

func Run(name string, args ...interface{}) (string, error) {
	argss := Strings(args...)
	outBytes, err := exec.Command(name, argss...).CombinedOutput()
	outString := strings.TrimSpace(string(outBytes))
	if err != nil {
		if outString == "" {
			outString = err.Error()
		}
		cmd := strings.Join(append([]string{name}, argss...), " ")
		return "", fmt.Errorf("%s: %s", cmd, outString)
	}
	return outString, nil
}
