package cmds

import "flag"

func isDefault(set *flag.FlagSet) func(string) bool  {
	return func(arg string) bool {
		f := set.Lookup(arg)
		return f.Value.String() == f.DefValue
	}
}
