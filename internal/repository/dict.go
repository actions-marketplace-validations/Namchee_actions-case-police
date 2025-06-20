package repository

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"time"
)

// GetDictionary returns word-to-word mapping dictionary
func GetDictionary(fsys fs.FS, files []string) map[string]string {
	dict := make(map[string]string)

	for _, filename := range files {
		fileDict := parseDictionaryFile(fsys, filename)

		for k, v := range fileDict {
			dict[k] = v
		}
	}

	return dict
}

func GetDictionaryFromBaseProject(files []string) map[string]string {
	baseUrl := "https://raw.githubusercontent.com/antfu/case-police/refs/heads/main/packages/case-police/dict"
	client := http.Client{
		Timeout: 3 * time.Second,
	}

	for _, filename := range files {
		fileUrl := fmt.Sprintf("%s/%s.json", baseUrl, filename)
		req, err := http.NewRequest(http.MethodGet, fileUrl, nil)
		res, httpErr := client.Do(req)

		defer res.Body.Close()
		var data map[string]string
		err = json.NewDecoder(res.Body).Decode(&data)
	}
}

func parseDictionaryFile(fsys fs.FS, filename string) map[string]string {
	file, err := fsys.Open(fmt.Sprintf("%s.json", filename))

	if err != nil {
		return map[string]string{}
	}

	var data map[string]string
	err = json.NewDecoder(file).Decode(&data)

	if err != nil {
		return map[string]string{}
	}

	return data
}
