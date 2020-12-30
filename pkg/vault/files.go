package vault

import (
	"io/ioutil"
	"net/url"
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
	if strings.HasPrefix(s, "file:/") {
		uri, err := url.ParseRequestURI(s)
		if err != nil {
			return "", err
		}
		file, err := ioutil.ReadFile(uri.Path)
		if err != nil {
			return "", err
		}
		return string(file), nil
	}
	return s, nil
}
