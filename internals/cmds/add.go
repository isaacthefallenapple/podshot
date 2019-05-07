package cmds

import (
	"fmt"
	"github.com/isaacthefallenapple/podshot/pkg/feedops"
	"os"
)

var addDefault = isDefault(addCmd)

func Add(args []string) {
	addCmd.Parse(args)
	if addDefault("name") || addDefault("feed") {
		addCmd.Usage()
		os.Exit(1)
	}
	succ, err := feedops.Add(*addNamePtr, *addFeedPtr)
	fmt.Print("\n\t")
	if succ {
		fmt.Printf("Successfully added a feed for %q", *addNamePtr)
	} else {
		if err != nil {
			fmt.Printf("Encountered an error while adding feed for %q:\n%s", *addNamePtr, err)
		} else {
			fmt.Printf("%q is already registered as a podcast", *addNamePtr)
		}
	}
	fmt.Println()
	os.Exit(1)
}
