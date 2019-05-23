package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/bob3000/dupcrawler/fshash"
	"github.com/docopt/docopt-go"
)

// VERSION referees to the programm's version
const VERSION = "1"

func main() {
	usage := `dupcrawler - looking for duplicate files

Usage:
  dupcrawler [options] <inputdir>

Options:
  -e, --excludes=<paths>  Comma separated list of path to exclude [default: ""]
  -d, --depth=<i>         Set maximum file hierarchy depth [default: 0]
  -l, --symlinks          Follow symlinks
  -h --help               Show this screen
  -v, --verbose           Show whats being done
  --version               Show version`

	arguments, _ := docopt.ParseArgs(usage, os.Args[1:], VERSION)

	maxDepth, err := strconv.Atoi(arguments["--depth"].(string))
	if err != nil {
		log.Fatal(err)
	}
	excludes := strings.Split(arguments["--excludes"].(string), ",")
	isVerbose := arguments["--verbose"].(bool)

	rArgs := fshash.ReadPathArgs{
		CurDepth:    0,
		Excludes:    excludes,
		FPath:       path.Clean(arguments["<inputdir>"].(string)),
		FollowLinks: arguments["--symlinks"].(bool),
		MaxDepth:    maxDepth,
		Parallel:    true,
		Verbose:     isVerbose,
	}

	fileHashes := fshash.ReadPath(rArgs)
	if isVerbose {
		fmt.Printf("\nDuplicate files:\n\n")
	}
	for _, group := range fileHashes {
		if len(group) > 1 {
			for _, p := range group {
				fmt.Println(p)
			}
			fmt.Println()
		}
	}
}
