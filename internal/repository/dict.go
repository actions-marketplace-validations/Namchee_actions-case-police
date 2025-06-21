package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/fatih/color"
)

// GetDictionary returns word-to-word mapping dictionary from the original case-police repository
func GetDictionary(files []string) map[string]string {
	allDictionary := make(map[string]string)

	baseUrl := "https://raw.githubusercontent.com/antfu/case-police/refs/heads/main/packages/case-police/dict"
	client := http.Client{
		Timeout: 1 * time.Second,
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(files))
	
	dictCh := make(chan map[string]string, len(files))
	errCh := make(chan error, len(files))

	for _, filename := range files {
		go func(filename string) {
			defer wg.Done()
			fileUrl := fmt.Sprintf("%s/%s.json", baseUrl, filename)

			errStr := fmt.Errorf("failed to get dictionary file for %s.json", filename)

			req, err := http.NewRequest(http.MethodGet, fileUrl, nil)
			if (err != nil) {
				errCh <- errStr
				return
			}

			res, httpErr := client.Do(req)
			if (httpErr != nil) {
				errCh <- errStr
				return
			}
			defer res.Body.Close()


			if res.StatusCode != http.StatusOK {
				errCh <- errStr
				return
			}

			var dictionaryContents map[string]string
			err = json.NewDecoder(res.Body).Decode(&dictionaryContents)

			if (err != nil) {
				errCh <- errStr
				return
			}

			dictCh <- dictionaryContents
		}(filename)
	}

	go func() {
		wg.Wait()
		close(dictCh)
		close(errCh)
	}()

	for dictionary := range dictCh {
		for key, val := range dictionary {
			allDictionary[key] = val
		}
	}

	for failure := range errCh {
		red := color.New(color.FgRed)
		red.Println(failure.Error())
	}

	return allDictionary
}
