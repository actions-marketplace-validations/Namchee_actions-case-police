package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Namchee/actions-case-police/internal"
	"github.com/Namchee/actions-case-police/internal/entity"
	"github.com/Namchee/actions-case-police/internal/repository"
	"github.com/Namchee/actions-case-police/internal/service"
	"github.com/Namchee/actions-case-police/internal/utils"
)

var (
	logger *log.Logger
)

func init() {
	logger = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lmsgprefix)
}

func main() {
	ctx := context.Background()

	cfg, err := entity.ReadConfiguration()
	if err != nil {
		logger.Fatalln(
			fmt.Errorf("failed to read action configuration: %w", err),
		)
	}

	meta, err := entity.CreateMeta(
		utils.ReadEnvString("GITHUB_REPOSITORY"),
	)
	if err != nil {
		logger.Fatalln(
			fmt.Errorf("failed to read metadata: %w", err),
		)
	}
	event, err := utils.GetEventNumber(os.DirFS("/"))
	if err != nil {
		logger.Fatalln(
			fmt.Errorf("failed to read repository event: %w", err),
		)
	}

	client := internal.NewGithubClient(ctx, cfg.Token)

	issue, err := client.GetIssue(ctx, meta, event)
	if err != nil {
		logger.Fatalln(
			fmt.Errorf("failed to get issue data: %w", err),
		)
	}

	dictionary := repository.GetDictionary(
		ctx,
		meta,
		client,
		cfg.Preset,
	)

	utils.MergeDictionary(&cfg.Dictionary, &dictionary)
	if len(cfg.Exclude) > 0 {
		utils.RemoveEntries(&cfg.Dictionary, cfg.Exclude)
	}

	result := service.PolicizeIssue(issue, cfg)

	if len(result.Changes) > 0 {
		err = client.EditIssue(ctx, meta, event, result)

		if err != nil {
			log.Fatalln(
				fmt.Errorf("failed to edit issue: %w", err),
			)
		}
	}

	service.LogResult(result, cfg)
}
