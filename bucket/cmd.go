package bucket

type Cmd struct {
	Known *knownCmd `arg:"subcommand:known"`
}

func (cmd *Cmd) Run() error {
	if cmd.Known != nil {
		return cmd.Known.Run()
	}

	return nil
}
