package cmds

import (
	"fmt"
	"github.com/isaacthefallenapple/podshot/pkg/feedops"
	"os"
)

var rmDefault = isDefault(rmCmd)

func Remove(args []string) {
	rmCmd.Parse(args)
	if rmDefault("name") {
		rmCmd.Usage()
		os.Exit(1)
	}
	succ, err := feedops.Remove(*rmNamePtr)
	fmt.Print("\n\t")
	if succ {
		fmt.Printf("Successfully removed %q", *rmNamePtr)
	} else {
		if err != nil {
			fmt.Printf("Encountered and error while trying to remove %q\n%s", *rmNamePtr, err)
		} else {
			fmt.Printf("No podcast with name %q registered", *rmNamePtr)
		}
	}
	fmt.Println()
	os.Exit(1)
}
