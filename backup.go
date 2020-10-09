package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	acl           = "private"
	dsPlaceholder = "{ds}"
	tmpDump       = "/tmp/backup.dump"
)

func main() {
	err := backup()
	if err != nil {
		log.Fatalf("Got Error Running Backup: %v", err)
	}
}

func backup() error {
	pgUrl := getEnv("POSTGRES_URL")
	dumpCmd := exec.Command("pg_dump", "-Fc", "-f"+tmpDump, pgUrl)
	out, err := dumpCmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		return err
	}
	defer os.Remove(tmpDump)

	f, err := os.Open(tmpDump)
	if err != nil {
		return err
	}
	defer f.Close()

	uploader, err := newUploader()
	if err != nil {
		return err
	}

	ds := time.Now().Format("2006-01-02")
	keys := strings.Split(getEnv("KEYS"), ",")
	for _, key := range keys {
		key = strings.ReplaceAll(key, dsPlaceholder, ds)
		upload, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(getEnv("S3_BUCKET")),
			Key:    aws.String(key),
			ACL:    aws.String(acl),
			Body:   f,
		})
		if err != nil {
			return err
		}
		fmt.Printf("> Uploaded %s\n", upload.Location)
	}

	return nil
}

func getEnv(key string) string {
	s := os.Getenv(key)
	if s == "" {
		log.Fatalf("Env Variable %s is missing.\n", key)
	}
	return s
}

func newUploader() (*s3manager.Uploader, error) {
	s, err := session.NewSession(&aws.Config{
		Region:   aws.String(getEnv("S3_REGION")),
		Endpoint: aws.String(getEnv("S3_ENDPOINT")),
		Credentials: credentials.NewStaticCredentials(
			getEnv("S3_KEY"),
			getEnv("S3_SECRET"),
			""),
	})
	if err != nil {
		return nil, err
	}
	return s3manager.NewUploader(s), nil
}
