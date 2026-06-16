package commands

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"text/tabwriter"

	"github.com/Protract-123/mocha/bucket"
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

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)

	_, err = fmt.Fprintln(w, "Name\tSource\tUpdated\tManifests")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, "----\t------\t-------\t---------")
	if err != nil {
		return err
	}

	for _, entry := range bucketMetadata {
		_, err = fmt.Fprintf(w, "%s\t%s\t%s\t%d\n",
			entry.Name,
			entry.Source,
			entry.LastUpdated,
			entry.ManifestCount,
		)
		if err != nil {
			return err
		}
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
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
	knownBuckets, err := bucket.GetKnownBuckets(mochaDir)
	if err != nil {
		return err
	}

	identifiedBucket, ok := findBucket(knownBuckets, cmd.Name)

	if !ok && cmd.RepositoryURL.String() == "" {
		return fmt.Errorf("bucket %s is not known, please provide URL", cmd.Name)
	} else if !ok && IsValidURL(cmd.RepositoryURL) {
		identifiedBucket = bucket.Bucket{
			Name:   cmd.Name,
			Source: cmd.RepositoryURL.String(),
		}
	} else if !ok {
		return fmt.Errorf("please provide valid URL for bucket %s", cmd.Name)
	}

	_, err = exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("git was not found in PATH, please install git by running `mocha install git`")
	}

	return bucket.DownloadBucket(identifiedBucket, mochaDir)
}

func findBucket(buckets []bucket.Bucket, name string) (bucket.Bucket, bool) {
	for _, entry := range buckets {
		if entry.Name == name {
			return entry, true
		}
	}
	return bucket.Bucket{}, false
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
