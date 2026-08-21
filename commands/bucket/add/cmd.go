package add

import (
	"fmt"
	"net/url"
	"os/exec"

	"github.com/Protract-123/mocha/bucket"
	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/output"
)

type Command struct {
	Name string   `arg:"positional,required" help:"bucket name (e.g. main)"`
	URL  *url.URL `arg:"positional" help:"git repository URL for the bucket; omit to use a known bucket"`
}

func (cmd *Command) Run() error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git is required to add buckets")
	}

	mochaDir := config.Current().MochaDirectory

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

	output.LogSuccess("successfully added bucket %q", identifiedBucket.Name)
	return nil
}
