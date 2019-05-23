package fshash

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
)

const chunkSizeKB = 256

// ReadPathArgs contains arguments for the ReadPath function
type ReadPathArgs struct {
	CurDepth    int
	Excludes    []string
	FPath       string
	FollowLinks bool
	MaxDepth    int
	Parallel    bool
	Verbose     bool
}

// Map is mapping of file hashes to a list of files
type Map map[string]FileList

// Sort sorts the FileList functioning as map values
func (fl *Map) Sort() {
	for _, l := range *fl {
		sort.Sort(l)
	}
}

// FileList is a just a list of file path list
type FileList []string

// Len returns the length of the file list
func (fm FileList) Len() int {
	return len(fm)
}

// Swap exchanges two values in the referenced list
func (fm FileList) Swap(i, j int) {
	fm[i], fm[j] = fm[j], fm[i]
}

// Less checks if i < j
func (fm FileList) Less(i, j int) bool {
	return fm[i] < fm[j]
}

func readPath(a ReadPathArgs, fm *Map, wg *sync.WaitGroup,
	mtx *sync.Mutex) {
	// return early because of constraints?
	if a.MaxDepth > 0 && a.CurDepth >= a.MaxDepth {
		return
	}
	for _, e := range a.Excludes {
		if strings.HasSuffix(a.FPath, e) {
			return
		}
	}

	// begin file examination
	fstat, err := os.Stat(a.FPath)
	if err != nil {
		log.Fatal(err)
	}
	if fstat.Mode() == os.ModeSymlink {
		if !a.FollowLinks {
			return
		}
		args := a
		args.FPath, err = os.Readlink(a.FPath)
		if err != nil {
			log.Fatal(err)
		}
		if a.Parallel {
			wg.Add(1)
			go readPath(args, fm, wg, mtx)
		} else {
			readPath(args, fm, wg, mtx)
		}
		return
	}
	f, err := os.Open(a.FPath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		f.Close()
		if a.Parallel {
			wg.Done()
		}
	}()

	// if it's a directory, call func recursively
	if fstat.IsDir() {
		if a.Verbose {
			fmt.Printf("Examining %s ...\n", a.FPath)
		}
		files, err := f.Readdirnames(0)
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
			args := a
			args.FPath = path.Join(a.FPath, f)
			args.CurDepth = a.CurDepth + 1
			if a.Parallel {
				wg.Add(1)
				go readPath(args, fm, wg, mtx)
			} else {
				readPath(args, fm, wg, mtx)
			}
		}
		return
	}

	// calculate hash and save it to fileHashes
	chunkSize := chunkSizeKB * 1024
	fChunk := make([]byte, chunkSize)
	var hash [sha1.Size]byte
	hashComponents := make([]byte, len(hash)+len(fChunk))
	hashStr := ""
	for {
		_, err = f.Read(fChunk)
		if err == io.EOF {
			break
		}
		copy(hashComponents, hashStr)
		copy(hashComponents[len(hashStr):], fChunk)
		hash = sha1.Sum(hashComponents)
		hashStr = base64.StdEncoding.EncodeToString(hash[:])
	}
	if a.Parallel {
		mtx.Lock()
		(*fm)[hashStr] = append((*fm)[hashStr], a.FPath)
		mtx.Unlock()
	} else {
		(*fm)[hashStr] = append((*fm)[hashStr], a.FPath)
	}
}

// ReadPath crawls the file system from a specified path and creates a mapping
// SHA1 hashes to file paths
func ReadPath(args ReadPathArgs) Map {
	var fileHashes = make(Map)
	if args.Parallel {
		var wg sync.WaitGroup
		var mtx = sync.Mutex{}
		wg.Add(1)
		readPath(args, &fileHashes, &wg, &mtx)
		wg.Wait()
	} else {
		readPath(args, &fileHashes, nil, nil)
	}
	fileHashes.Sort()
	return fileHashes
}
