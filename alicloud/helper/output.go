package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"sort"
	"strings"
)

// GetUserHomeDir returns the current user's home directory.
func GetUserHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("Get current user got an error: %#v.", err)
	}
	return usr.HomeDir, nil
}

// WriteToFile persists data to filePath: strings are written raw, other
// types are JSON-marshaled first. A leading "~" expands to the user's home
// directory. An existing target file is removed before writing.
func WriteToFile(filePath string, data interface{}) error {
	var out string
	switch data.(type) {
	case string:
		out = data.(string)
		break
	case nil:
		return nil
	default:
		bs, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			return fmt.Errorf("MarshalIndent data %#v got an error: %#v", data, err)
		}
		out = string(bs)
	}

	if strings.HasPrefix(filePath, "~") {
		home, err := GetUserHomeDir()
		if err != nil {
			return err
		}
		if home != "" {
			filePath = strings.Replace(filePath, "~", home, 1)
		}
	}

	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}

	return ioutil.WriteFile(filePath, []byte(out), 422)
}

// DataResourceIdHash generates a hash for the set hash function used by
// data source IDs.
func DataResourceIdHash(ids []string) string {
	var buf bytes.Buffer

	for _, id := range ids {
		buf.WriteString(fmt.Sprintf("%s-", id))
	}

	return fmt.Sprintf("%d", Hashcode(buf.String()))
}

// DataResourceIdHashFromMap derives a stable data source ID from a request map:
// keys are sorted so the same query condition always produces the same ID.
func DataResourceIdHashFromMap(request map[string]interface{}) string {
	keys := make([]string, 0, len(request))
	for key := range request {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s=%v", key, request[key]))
	}
	return DataResourceIdHash(parts)
}
