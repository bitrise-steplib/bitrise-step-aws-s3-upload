package step

import (
	"github.com/bitrise-io/go-steputils/v2/stepconf"
)

type Input struct {
	Path            string          `env:"path,required"`
	BucketName      string          `env:"bucket_name,required"`
	BucketPrefix    string          `env:"bucket_prefix"`
	AccessControl   string          `env:"acl_control,opt[private,public-read]"`
	Region          string          `env:"region,required"`
	AccessKeyID     stepconf.Secret `env:"access_key_id,required"`
	SecretAccessKey stepconf.Secret `env:"secret_access_key,required"`
	SessionToken    stepconf.Secret `env:"session_token"`
	Verbose         bool            `env:"verbose,opt[true,false]"`
}

type Config struct {
	Path            string
	BucketName      string
	BucketPrefix    string
	AccessControl   string
	Region          string
	AccessKeyID     stepconf.Secret
	SecretAccessKey stepconf.Secret
	SessionToken    stepconf.Secret
}
