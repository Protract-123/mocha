package bucket

import (
	"os"
	"strings"
)

type Bucket struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Cmd struct {
	Add   *addCmd   `arg:"subcommand:add"`
	Known *knownCmd `arg:"subcommand:known"`
	Rm    *rmCmd    `arg:"subcommand:rm"`
	List  *listCmd  `arg:"subcommand:list"`
}

func (cmd *Cmd) Run(mochaDir string) error {
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

func parseBucketList(file string) ([]Bucket, error) {
	bucketsJson, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var buckets []Bucket

	for _, line := range strings.Split(string(bucketsJson), "\n") {
		line = strings.TrimSpace(line)

		if line == "{" || line == "}" || line == "" {
			continue
		}

		parts := strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			continue
		}

		name := strings.Trim(parts[0], `"`)
		url := strings.Trim(parts[1], `",`)

		buckets = append(buckets, Bucket{Name: name, URL: url})
	}

	return buckets, nil
}
