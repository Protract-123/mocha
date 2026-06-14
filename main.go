package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/Protract-123/mocha/bucket"
	"github.com/alexflint/go-arg"
)

var args struct {
	MochaDir string `arg:"--,env:MOCHA_DIR" default:"$USERPROFILE/mocha"`

	Bucket *bucket.Cmd `arg:"subcommand:bucket"`
}

func main() {
	arg.MustParse(&args)

	mochaDir, err := filepath.Abs(os.ExpandEnv(args.MochaDir))
	if err != nil {
		log.Fatalf("failed to resolve mocha dir: %v", err)
	}

	if args.Bucket != nil {
		err := args.Bucket.Run(mochaDir)
		if err != nil {
			println(err.Error())
		}
	}
}
