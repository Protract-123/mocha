package bucket

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"
)

type listCmd struct{}

func (cmd *listCmd) Run(mochaDir string) error {
	bucketsDir := filepath.Join(mochaDir, "buckets")

	buckets, err := os.ReadDir(bucketsDir)
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

	for _, entry := range buckets {
		if !entry.IsDir() {
			continue
		}

		bucketName := entry.Name()
		bucketPath := filepath.Join(bucketsDir, bucketName)

		sourceCmd := exec.Command("git", "config", "remote.origin.url")
		sourceCmd.Dir = bucketPath
		sourceOut, err := sourceCmd.Output()
		if err != nil {
			continue
		}
		bucketSource := strings.TrimSpace(string(sourceOut))

		updatedCmd := exec.Command("git", "log", "--format=%aD", "-n", "1")
		updatedCmd.Dir = bucketPath
		updatedOut, err := updatedCmd.Output()
		if err != nil {
			continue
		}
		bucketUpdated, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", strings.TrimSpace(string(updatedOut)))
		if err != nil {
			continue
		}

		manifests, err := os.ReadDir(filepath.Join(bucketPath, "bucket"))
		if err != nil {
			continue
		}
		manifestCount := len(manifests)

		_, err = fmt.Fprintf(w, "%s\t%s\t%s\t%d\n",
			bucketName,
			bucketSource,
			bucketUpdated.Format("02-01-2006 15:04:05"),
			manifestCount,
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
