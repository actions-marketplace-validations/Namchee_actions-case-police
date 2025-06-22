package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Namchee/actions-case-police/internal"
	"github.com/Namchee/actions-case-police/internal/entity"

	"github.com/fatih/color"
)

var (
	originalRepositoryMeta = &entity.Meta{
		Owner: "antfu",
		Name:  "case-police",
	}
)

// GetDictionary returns word-to-word mapping dictionary from the original case-police repository
func GetDictionary(ctx context.Context, client internal.GithubClient, files []string) map[string]string {
	allDictionary := make(map[string]string)

	wg := &sync.WaitGroup{}
	wg.Add(len(files))
	
	dictCh := make(chan map[string]string, len(files))
	errCh := make(chan error, len(files))

	for _, filename := range files {
		go func(filename string) {
			defer wg.Done()
			fileUrl := fmt.Sprintf("packages/case-police/dict/%s.json", filename)
			errStr := fmt.Errorf("failed to get contents of %s.json", filename)

			rawContent, err := client.GetRepositoryContents(ctx, originalRepositoryMeta, fileUrl)
			if (err != nil) {
				errCh <- errStr
				return
			}

			var dictionaryContents map[string]string
			err = json.Unmarshal([]byte(rawContent), &dictionaryContents)

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

	// Failing to get a dictionyar is OK.
	for failure := range errCh {
		red := color.New(color.FgRed)
		red.Println(failure.Error())
	}

	return allDictionary
}
