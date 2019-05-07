package main

import (
	"flag"
	"github.com/isaacthefallenapple/podshot/internals/cmds"
	"os"
)

func main() {

	switch os.Args[1] {
	case "add":
		cmds.Add(os.Args[2:])
	case "remove":
		cmds.Remove(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
	/*
		if addCommand.Parsed() {

			if *feedPtr == "" {
				addCommand.Usage()
				os.Exit(1)
			}
			if *namePtr == "" {
				addCommand.Usage()
				os.Exit(1)
			}

			f := feedops.Feeds()
			feedops.Add(*namePtr, *feedPtr)
			fmt.Println(f)
		}*/
}
