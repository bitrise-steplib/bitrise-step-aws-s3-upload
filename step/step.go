package step

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
)

type Uploader struct {
	inputParser   stepconf.InputParser
	envRepository env.Repository
	logger        log.Logger
}

func NewUploader(inputParser stepconf.InputParser, envRepository env.Repository, logger log.Logger) Uploader {
	return Uploader{
		inputParser:   inputParser,
		envRepository: envRepository,
		logger:        logger,
	}
}

func (u Uploader) ProcessConfig() (Config, error) {
	var input Input
	err := u.inputParser.Parse(&input)
	if err != nil {
		return Config{}, err
	}

	stepconf.Print(input)
	u.logger.Println()
	u.logger.EnableDebugLog(input.Verbose)

	return Config{
		Path:            input.Path,
		BucketName:      input.BucketName,
		BucketPrefix:    input.BucketPrefix,
		AccessControl:   input.AccessControl,
		Region:          input.Region,
		AccessKeyID:     input.AccessKeyID,
		SecretAccessKey: input.SecretAccessKey,
		SessionToken:    input.SessionToken,
	}, nil
}

func (u Uploader) Run(config Config) error {
	items, err := u.collectItems(config.Path)
	if err != nil {
		return err
	}

	ctx := context.Background()
	defaultConfig, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion(config.Region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				string(config.AccessKeyID),
				string(config.SecretAccessKey),
				string(config.SessionToken),
			),
		),
	)
	if err != nil {
		return fmt.Errorf("load AWS config: %w", err)
	}

	client := s3.NewFromConfig(defaultConfig)
	uploader := manager.NewUploader(client)
	prefix := strings.Trim(config.BucketPrefix, "/")
	
	for _, item := range items {
		u.logger.Infof("Uploading file: %s", item.Path)

		f, err := os.Open(item.Path)
		if err != nil {
			return fmt.Errorf("open %s: %w", item.Path, err)
		}

		var key string
		if prefix == "" {
			key = item.Key
		} else {
			key = path.Join(prefix, item.Key)
		}

		_, err = uploader.Upload(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(config.BucketName),
			Key:         aws.String(key),
			Body:        f,
			ContentType: aws.String(item.ContentType),
		})
		if err != nil {
			return fmt.Errorf("upload %s: %w", item.Path, err)
		}

		if err := f.Close(); err != nil {
			u.logger.Warnf("Failed to close file %s: %s", item.Path, err)
		}
	}

	return nil
}
