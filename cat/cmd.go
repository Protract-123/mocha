package cat

import (
	"fmt"
	"os"
	"path/filepath"
)

type Cmd struct {
	Name   string `arg:"positional,required"`
	Bucket string `arg:"-b, --bucket"`
}

func (cmd Cmd) Run(mochaDir string) error {
	bucketsDir := filepath.Join(mochaDir, "buckets")

	if cmd.Bucket != "" {
		manifestPath := filepath.Join(bucketsDir, cmd.Bucket, "bucket", fmt.Sprintf("%s.json", cmd.Name))

		data, err := os.ReadFile(manifestPath)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("app %s not found in bucket %s", cmd.Name, cmd.Bucket)
			}
			return err
		}

		_, err = os.Stdout.Write(data)
		if err != nil {
			return err
		}

		return nil
	}

	buckets, err := os.ReadDir(bucketsDir)
	if err != nil {
		return err
	}

	for _, bucket := range buckets {
		if !bucket.IsDir() {
			continue
		}

		manifestPath := filepath.Join(bucketsDir, bucket.Name(), "bucket", fmt.Sprintf("%s.json", cmd.Name))

		data, err := os.ReadFile(manifestPath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}

		_, err = os.Stdout.Write(data)
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("app %s not found in buckets", cmd.Name)
}
