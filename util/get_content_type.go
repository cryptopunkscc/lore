package util

import (
	"io"
	"os/exec"
	"strings"
)

func GetContentType(data io.Reader) (string, error) {
	cmd := exec.Command("file", "--mime-type", "-")
	cmd.Stdin = data

	res, err := cmd.Output()
	if err != nil {
		return "", err
	}

	sep := strings.SplitAfterN(string(res), ":", 2)
	typ := strings.TrimSpace(sep[1])

	return typ, nil
}
