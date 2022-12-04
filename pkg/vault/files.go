package vault

import (
	"net/url"
	"os"
	"strings"
)

func maybeFromFiles(s1 string, s2 string) (v1 string, v2 string, err error) {
	v1, err = maybeFromFile(s1)
	if err != nil {
		return "", "", err
	}
	v2, err = maybeFromFile(s2)
	if err != nil {
		return "", "", err
	}
	return v1, v2, nil
}

func maybeFromFile(s string) (string, error) {
	if !strings.HasPrefix(s, "file:/") {
		return s, nil
	} else {
		path, err := absoluteFilePath(s)
		if err != nil {
			return "", err
		}
		file, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return string(file), nil
	}
}

func absoluteFilePath(s string) (path string, err error) {
	// support relative paths and user directory:
	if strings.HasPrefix(s, "file://.") {
		return s[len("file://"):], nil
	} else if strings.HasPrefix(s, "file://~") {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		} else {
			return userHomeDir + s[len("file://")+1:], nil
		}
	} else {
		// other paths:
		uri, err := url.ParseRequestURI(s)
		if err == nil {
			return uri.Path, nil
		} else {
			return "", err
		}
	}
}
