package step

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/bitrise-step-aws-s3-upload/step/mocks"
	"github.com/stretchr/testify/assert"
)

func TestConfigParsing(t *testing.T) {
	config := Config{
		Path:            "path",
		BucketName:      "bucket-name",
		BucketPrefix:    "bucket-prefix",
		AccessControl:   "private",
		Region:          "region",
		AccessKeyID:     "access-key-id",
		SecretAccessKey: "secret-access-key",
		SessionToken:    "session-token",
	}

	mockEnvRepository := mocks.NewRepository(t)
	mockEnvRepository.On("Get", "path").Return(config.Path)
	mockEnvRepository.On("Get", "bucket_name").Return(config.BucketName)
	mockEnvRepository.On("Get", "bucket_prefix").Return(config.BucketPrefix)
	mockEnvRepository.On("Get", "acl_control").Return(config.AccessControl)
	mockEnvRepository.On("Get", "region").Return(config.Region)
	mockEnvRepository.On("Get", "access_key_id").Return(string(config.AccessKeyID))
	mockEnvRepository.On("Get", "secret_access_key").Return(string(config.SecretAccessKey))
	mockEnvRepository.On("Get", "session_token").Return(string(config.SessionToken))
	mockEnvRepository.On("Get", "verbose").Return("false")

	inputParser := stepconf.NewInputParser(mockEnvRepository)
	sut := NewUploader(inputParser, mockEnvRepository, log.NewLogger())

	receivedConfig, err := sut.ProcessConfig()
	assert.NoError(t, err)
	assert.Equal(t, config, receivedConfig)

	mockEnvRepository.AssertExpectations(t)
}

func TestPutObjectInputACL(t *testing.T) {
	tests := []struct {
		name          string
		accessControl string
		expectedACL   types.ObjectCannedACL
	}{
		{name: "public-read is applied", accessControl: "public-read", expectedACL: types.ObjectCannedACLPublicRead},
		{name: "private leaves ACL unset", accessControl: "private", expectedACL: ""},
		{name: "empty leaves ACL unset", accessControl: "", expectedACL: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{BucketName: "bucket", AccessControl: tt.accessControl}
			input := newPutObjectInput(config, "key", "application/octet-stream", nil)

			assert.Equal(t, tt.expectedACL, input.ACL)
			assert.Equal(t, "bucket", *input.Bucket)
			assert.Equal(t, "key", *input.Key)
		})
	}
}
