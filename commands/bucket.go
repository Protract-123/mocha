package commands

import (
	"fmt"
	"net/url"
	"os/exec"
	"strconv"

	"github.com/Protract-123/mocha/bucket"
	"github.com/Protract-123/mocha/output"
	"github.com/alexflint/go-arg"
)

type BucketCommand struct {
	Add    *addBucketCommand    `arg:"subcommand:add" help:"add a bucket by name or git repository URL"`
	Known  *knownBucketsCommand `arg:"subcommand:known" help:"list known buckets available to add by name"`
	Remove *removeBucketCommand `arg:"subcommand:remove" help:"remove an installed bucket"`
	List   *listBucketsCommand  `arg:"subcommand:list" help:"list installed buckets"`
}

type listBucketsCommand struct{}
type knownBucketsCommand struct{}
type removeBucketCommand struct {
	Name string `arg:"positional,required" help:"bucket name (e.g. main)"`
}
type addBucketCommand struct {
	Name string   `arg:"positional,required" help:"bucket name (e.g. main)"`
	URL  *url.URL `arg:"positional" help:"git repository URL for the bucket; omit to use a known bucket"`
}

func (cmd *BucketCommand) Run(mochaDir string) error {
	switch {
	case cmd.Add != nil:
		return cmd.Add.Run(mochaDir)
	case cmd.Known != nil:
		return cmd.Known.Run(mochaDir)
	case cmd.Remove != nil:
		return cmd.Remove.Run(mochaDir)
	case cmd.List != nil:
		return cmd.List.Run(mochaDir)
	default:
		return arg.ErrHelp
	}
}

func (cmd *listBucketsCommand) Run(mochaDir string) error {
	bucketMetadata, err := bucket.GetAllBucketMetadata(mochaDir)
	if err != nil {
		return fmt.Errorf("failed to get bucket metadata: %w", err)
	}

	if len(bucketMetadata) == 0 {
		output.LogOutput("no buckets to show")
		return nil
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

	if err := output.PrintTable(headers, rows); err != nil {
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
		output.LogOutput(fmt.Sprint(entry.Name, ": ", entry.Source))
	}

	return nil
}

func (cmd *removeBucketCommand) Run(mochaDir string) error {
	if err := bucket.DeleteBucket(cmd.Name, mochaDir); err != nil {
		return fmt.Errorf("failed to delete bucket %q: %w", cmd.Name, err)
	}

	return nil
}

func (cmd *addBucketCommand) Run(mochaDir string) error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git is required to add buckets")
	}

	var identifiedBucket bucket.Bucket

	if cmd.URL == nil {
		knownBucket, err := bucket.GetKnownBucket(cmd.Name, mochaDir)
		if err != nil {
			return fmt.Errorf("failed to get known bucket: %w", err)
		}
		identifiedBucket = knownBucket
	} else {
		if (cmd.URL.Scheme != "http" && cmd.URL.Scheme != "https") || cmd.URL.Host == "" {
			return fmt.Errorf("invalid repository URL %q", cmd.URL.String())
		}

		identifiedBucket = bucket.Bucket{
			Name:   cmd.Name,
			Source: cmd.URL.String(),
		}
	}

	if err := bucket.DownloadBucket(identifiedBucket, mochaDir); err != nil {
		return fmt.Errorf("failed to download bucket %q: %w", identifiedBucket.Name, err)
	}

	return nil
}
