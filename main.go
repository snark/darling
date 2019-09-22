package main

import (
	"flag"
	"fmt"
	"github.com/snark/darling/internal/cmd/darling"
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

func main() {
	var blacklistWords arrayFlags
	var whitelistWords arrayFlags
	flag.Var(&blacklistWords, "b", "blacklist term")
	flag.Var(&whitelistWords, "w", "whitelist term")
	flag.Usage = func() {
		fmt.Printf("Usage: darling [options] <feed_url>...\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	tail := flag.Args()

	if len(tail) > 0 {
		darling.FilterFeeds(blacklistWords, whitelistWords, tail)
	} else {
		flag.Usage()
	}
}
