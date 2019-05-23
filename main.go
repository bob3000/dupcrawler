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
  dupcrawler [options] <inputdir>...

Options:
  -e, --excludes=<paths>  Comma separated list of path to exclude [default: ""]
  -d, --depth=<i>         Set maximum file hierarchy depth [default: 0]
  -l, --symlinks          Follow symlinks
  -h --help               Show this screen
  -v, --verbose           Show whats being done
  --no-parallel			  Disable parallel directory processing
  --version               Show version`

	arguments, _ := docopt.ParseArgs(usage, os.Args[1:], VERSION)

	maxDepth, err := strconv.Atoi(arguments["--depth"].(string))
	if err != nil {
		log.Fatal(err)
	}
	excludes := strings.Split(arguments["--excludes"].(string), ",")
	isVerbose := arguments["--verbose"].(bool)
	inputPaths := arguments["<inputdir>"].([]string)
	globalHashMap := make(fshash.Map)
	for _, p := range inputPaths {
		cleanPath := path.Clean(p)
		rArgs := fshash.ReadPathArgs{
			CurDepth:    0,
			Excludes:    excludes,
			FPath:       cleanPath,
			FollowLinks: arguments["--symlinks"].(bool),
			MaxDepth:    maxDepth,
			Parallel:    !arguments["--no-parallel"].(bool),
			Verbose:     isVerbose,
		}
		fileHashes := fshash.ReadPath(rArgs)
		for k, v := range fileHashes {
			globalHashMap[k] = append(globalHashMap[k], v...)
		}
	}

	if isVerbose {
		fmt.Printf("\nDuplicate files:\n\n")
	}
	for _, group := range globalHashMap {
		if len(group) > 1 {
			for _, p := range group {
				fmt.Println(p)
			}
			fmt.Println()
		}
	}
}
