package utils

import "github.com/urfave/cli/v2"

func FlagsMerge(flags ...cli.Flag) []cli.Flag {
	return append(make([]cli.Flag, 0), flags...)
}
