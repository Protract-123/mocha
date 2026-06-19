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
		return err
	}

	headers := []string{"Name", "Source", "Updated", "Manifests"}
	rows := make([][]string, len(bucketMetadata))

	for index, entry := range bucketMetadata {
		rows[index] = []string{
			entry.Name,
			entry.Source,
			entry.LastUpdated,
			strconv.Itoa(entry.ManifestCount),
		}
	}

	return output.PrintTable(headers, rows)
}

func (cmd *knownBucketsCommand) Run(mochaDir string) error {
	knownBuckets, err := bucket.GetKnownBuckets(mochaDir)
	if err != nil {
		return err
	}

	for _, entry := range knownBuckets {
		fmt.Print("\033[1;35m", entry.Name, "\033[0m: ", entry.Source, "\n")
	}

	return nil
}

func (cmd *removeBucketCommand) Run(mochaDir string) error {
	return bucket.DeleteBucket(cmd.Name, mochaDir)
}

func (cmd *addBucketCommand) Run(mochaDir string) error {
	_, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("git was not found in PATH, please install git by running `mocha install git`")
	}

	identifiedBucket, err := bucket.GetKnownBucket(cmd.Name, mochaDir)
	if err == nil {
		return bucket.DownloadBucket(identifiedBucket, mochaDir)
	}

	if cmd.RepositoryURL.String() == "" || !IsValidURL(cmd.RepositoryURL) {
		return fmt.Errorf("bucket %s is not known, please provide a valid URL", cmd.Name)
	}

	identifiedBucket = bucket.Bucket{
		Name:   cmd.Name,
		Source: cmd.RepositoryURL.String(),
	}

	return bucket.DownloadBucket(identifiedBucket, mochaDir)
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
