package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"sort"
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
  -e, --excludes=<paths>  Comma separated list of paths to exclude [default: ""]
  -d, --depth=<i>         Set maximum file hierarchy depth [default: 0]
  -l, --symlinks          Follow symlinks
  -h, --help              Show this screen
  -v, --verbose           Show what's being done
  --no-parallel           Disable parallel processing
  --no-sample             Calculate checksum over the entire file
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
			Sample:      !arguments["--no-sample"].(bool),
			Verbose:     isVerbose,
		}
		fileHashes := fshash.ReadPath(rArgs)
		for k, v := range fileHashes {
			globalHashMap[k] = append(globalHashMap[k], v...)
		}
	}

	// sort results
	keyHashes := make([]string, 0, len(globalHashMap))
	for k := range globalHashMap {
		keyHashes = append(keyHashes, k)
	}
	sort.Slice(keyHashes, func(i, j int) bool {
		return keyHashes[i] < keyHashes[j]
	})

	if isVerbose {
		fmt.Printf("\nDuplicate files:\n\n")
	}
	for _, h := range keyHashes {
		group := globalHashMap[h]
		if len(group) > 1 {
			for _, p := range group {
				fmt.Println(p)
			}
			fmt.Println()
		}
	}
}
