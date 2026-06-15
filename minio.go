package minio

import (
	"fmt"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/omeid/conex"
)

var (
	// Image to use for the box.
	Image = "minio/minio:latest"
	// Port used for connecting to minio.
	Port = "9000"

	// MinioUpWaitTime dictates how long we should wait for Minio to accept connections.
	MinioUpWaitTime = 10 * time.Second
)

func init() {
	conex.Require(func() string { return Image })
}

// Config used to configure the Minio box.
type Config struct {
	AccessKeyID     string
	SecretAccessKey string
}

// Client aliases minio.Client to avoid requiring the user to import the minio-go SDK directly.
type Client = minio.Client

// MakeBucketOptions aliases minio.MakeBucketOptions.
type MakeBucketOptions = minio.MakeBucketOptions

// PutObjectOptions aliases minio.PutObjectOptions.
type PutObjectOptions = minio.PutObjectOptions

// GetObjectOptions aliases minio.GetObjectOptions.
type GetObjectOptions = minio.GetObjectOptions

// Box returns a Client and the container running the minio
// server. It calls t.Fatal on errors.
func Box(t testing.TB, config *Config) (*Client, conex.Container) {
	if config == nil {
		config = &Config{}
	}
	if config.AccessKeyID == "" {
		config.AccessKeyID = "minioadmin"
	}
	if config.SecretAccessKey == "" {
		config.SecretAccessKey = "minioadmin"
	}

	c := conex.Box(t, &conex.Config{
		Image:  Image,
		Cmd:    []string{"server", "/data"},
		Env: []string{
			"MINIO_ROOT_USER=" + config.AccessKeyID,
			"MINIO_ROOT_PASSWORD=" + config.SecretAccessKey,
		},
		Expose: []string{Port},
	})

	conex.Logf(t, "minio", "Waiting for Minio to accept connections")

	err := c.Wait(Port, MinioUpWaitTime)
	if err != nil {
		c.Drop()
		t.Fatal("Minio failed to start:", err)
	}

	conex.Logf(t, "minio", "Minio is now accepting connections")

	endpoint := fmt.Sprintf("%s:%s", c.Address(), Port)

	// Initialize minio client object.
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		c.Drop()
		t.Fatal("Failed to initialize Minio client:", err)
	}

	return client, c
}
