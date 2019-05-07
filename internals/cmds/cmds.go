package cmds

import (
	"flag"
)

var (
	nestedCmd = flag.NewFlagSet("select", flag.ExitOnError)
	nestedNamePtr = nestedCmd.String("name", "", "Name of the Podcast (required)")
)

var (
	addCmd = flag.NewFlagSet("add", flag.ExitOnError)
	addNamePtr = addCmd.String("name", "", "Name of the Podcast (required)")
	addFeedPtr = addCmd.String("feed", "", "URL to RSS feed (required)")
)

var (
	rmCmd = flag.NewFlagSet("rm", flag.ExitOnError)
	rmNamePtr = rmCmd.String("name", "", "Name of the Podcast to be removed (required)")
)

var (
	getCmd = flag.NewFlagSet("get", flag.ExitOnError)
	getNamePtr = getCmd.String("name", "", "Name of the Podcast to be downloaded (required)")
	getPathPtr = getCmd.String("path", "%USERPROFILE%/Downloads", "Directory to download to (required)")
)
