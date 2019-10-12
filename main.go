package main

import (
	"fmt"
	"github.com/snark/darling/internal/cmd/darling"
	flag "github.com/spf13/pflag"
	"strings"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i *arrayFlags) Type() string {
	return "string"
}

func main() {
	var blacklistWords arrayFlags
	var whitelistWords arrayFlags
	flag.VarP(&blacklistWords, "blacklist", "b", "blacklist term")
	flag.VarP(&whitelistWords, "whitelist", "w", "whitelist term")
	var number = flag.IntP("limit", "n", 0, "restrict to n matching items per feed")
	var after = flag.String("since", "", "restrict to items after a given time")
	flag.Usage = func() {
		fmt.Printf("Usage: darling [options] <feed_url>...\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	tail := flag.Args()

	if len(tail) > 0 && *number >= 0 {
		darling.FilterFeeds(blacklistWords, whitelistWords, after, number, tail)
	} else {
		flag.Usage()
	}
}
