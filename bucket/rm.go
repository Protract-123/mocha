package bucket

import (
	"os"
	"path/filepath"
)

type rmCmd struct {
	Name string `arg:"positional,required"`
}

func (cmd *rmCmd) Run(mochaDir string) error {
	bucketsDir := filepath.Join(mochaDir, "buckets")
	destDir := filepath.Join(bucketsDir, cmd.Name)

	err := os.RemoveAll(destDir)
	if err != nil {
		return err
	}
	return nil
}
