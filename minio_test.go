package minio_test

import (
	"context"
	"os"
	"testing"

	"github.com/conex/minio"
	"github.com/omeid/conex"
)

func TestMain(m *testing.M) {
	os.Exit(conex.Run(m))
}

func TestMinio(t *testing.T) {
	client, _ := minio.Box(t, nil)

	ctx := context.Background()
	bucketName := "testbucket"
	location := "us-east-1"

	err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if run multiple times)
		exists, errBucketExists := client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			t.Logf("We already own %s\n", bucketName)
		} else {
			t.Fatal(err)
		}
	} else {
		t.Logf("Successfully created %s\n", bucketName)
	}
}
