package main

import (
	"github.com/Protract-123/mocha/bucket"
	"github.com/alexflint/go-arg"
)

var args struct {
	MochaDir string `arg:"env" default:"%USERPROFILE%/mocha"`

	Bucket *bucket.Cmd `arg:"subcommand:bucket"`
}

func main() {
	arg.MustParse(&args)

	if args.Bucket != nil {
		err := args.Bucket.Run()
		if err != nil {
			println(err)
		}
	}
}
