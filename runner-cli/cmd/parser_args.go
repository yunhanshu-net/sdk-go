package cmd

func InitArgs(args []string) {

	if args[0] == "run" {
		Run(args[1:])
	}

	if args[0] == "connect" {
		Connect(args[1:])
	}

	if args[0] == "send" {
		Send(args[1:])
	}

}
