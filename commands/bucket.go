package commands

import (
	"fmt"
	"net/url"
	"os/exec"
	"strconv"

	"github.com/Protract-123/mocha/bucket"
	"github.com/Protract-123/mocha/output"
)

type BucketCommand struct {
	Add   *addBucketCommand    `arg:"subcommand:add"`
	Known *knownBucketsCommand `arg:"subcommand:known"`
	Rm    *removeBucketCommand `arg:"subcommand:rm"`
	List  *listBucketsCommand  `arg:"subcommand:list"`
}

type listBucketsCommand struct{}
type knownBucketsCommand struct{}
type removeBucketCommand struct {
	Name string `arg:"positional,required"`
}
type addBucketCommand struct {
	Name          string  `arg:"positional,required"`
	RepositoryURL url.URL `arg:"positional"`
}

func (cmd *BucketCommand) Run(mochaDir string) error {
	if cmd.Known != nil {
		return cmd.Known.Run(mochaDir)
	}
	if cmd.Add != nil {
		return cmd.Add.Run(mochaDir)
	}
	if cmd.Rm != nil {
		return cmd.Rm.Run(mochaDir)
	}
	if cmd.List != nil {
		return cmd.List.Run(mochaDir)
	}

	return nil
}

func (cmd *listBucketsCommand) Run(mochaDir string) error {
	bucketMetadata, err := bucket.GetAllBucketMetadata(mochaDir)
	if err != nil {
		return fmt.Errorf("failed to get bucket metadata: %w", err)
	}

	headers := []string{"Name", "Source", "Updated", "Manifests"}
	rows := make([][]string, len(bucketMetadata))

	for index, entry := range bucketMetadata {
		rows[index] = []string{
			entry.Name,
			entry.Source,
			entry.LastUpdated.Format("02-01-2006 15:04:05"),
			strconv.Itoa(entry.ManifestCount),
		}
	}

	err = output.PrintTable(headers, rows)
	if err != nil {
		return fmt.Errorf("failed to display bucket metadata: %w", err)
	}

	return nil
}

func (cmd *knownBucketsCommand) Run(mochaDir string) error {
	knownBuckets, err := bucket.GetKnownBuckets(mochaDir)
	if err != nil {
		return fmt.Errorf("failed to get known buckets: %w", err)
	}

	for _, entry := range knownBuckets {
		fmt.Print(output.AnsiBoldMagenta, entry.Name, output.AnsiReset, ": ", entry.Source, "\n")
	}

	return nil
}

func (cmd *removeBucketCommand) Run(mochaDir string) error {
	err := bucket.DeleteBucket(cmd.Name, mochaDir)
	if err != nil {
		return fmt.Errorf("failed to delete bucket %q: %w", cmd.Name, err)
	}

	return nil
}

func (cmd *addBucketCommand) Run(mochaDir string) error {
	if cmd.Name == "" {
		return fmt.Errorf("no bucket name specified")
	}

	_, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("git is required to add buckets, please install git by running `mocha install git`")
	}

	identifiedBucket, err := bucket.GetKnownBucket(cmd.Name, mochaDir)
	if err != nil {
		if cmd.RepositoryURL.String() == "" || !IsValidURL(cmd.RepositoryURL) {
			return fmt.Errorf("bucket %s is not known, please provide a valid URL", cmd.Name)
		}

		identifiedBucket = bucket.Bucket{
			Name:   cmd.Name,
			Source: cmd.RepositoryURL.String(),
		}
	}

	err = bucket.DownloadBucket(identifiedBucket, mochaDir)
	if err != nil {
		return fmt.Errorf("failed to download bucket %q: %w", identifiedBucket.Name, err)
	}

	return nil
}

func IsValidURL(url url.URL) bool {
	if url.Scheme != "http" && url.Scheme != "https" {
		return false
	}
	if url.Host == "" {
		return false
	}

	return true
}
