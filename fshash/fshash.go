package fshash

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"
	"sync"

	"gopkg.in/cheggaaa/pb.v1"
)

const defaultChunkSize = 256 * 1024
const sampleRatio int64 = 16
const cpuMultiplicator = 8

// ReadPathArgs contains arguments for the ReadPath function
type ReadPathArgs struct {
	CurDepth    int
	Excludes    []string
	FPath       string
	FollowLinks bool
	MaxDepth    int
	Parallel    bool
	Sample      bool
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

func calcHash(fpath string, wg *sync.WaitGroup) string {
	f, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		f.Close()
		if wg != nil {
			wg.Done()
		}
	}()
	fstat, err := os.Stat(fpath)
	if err != nil {
		log.Fatal(err)
	}
	fSize := fstat.Size()
	var sampleChunk int64
	if sampleRatio > 0 && fSize > defaultChunkSize {
		sampleChunk = fSize / sampleRatio
	} else {
		sampleChunk = defaultChunkSize
	}
	fChunk := make([]byte, sampleChunk)
	var hash [sha1.Size]byte
	hashComponents := make([]byte, len(hash)+len(fChunk))
	hashStr := ""
	for {
		_, err := f.Read(fChunk)
		if err == io.EOF {
			break
		}
		copy(hashComponents, hashStr)
		copy(hashComponents[len(hashStr):], fChunk)
		hash = sha1.Sum(hashComponents)
		hashStr = base64.StdEncoding.EncodeToString(hash[:])
		if sampleRatio > 0 {
			_, err := f.Seek(sampleChunk*(sampleRatio-1), 1)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	return hashStr
}

func readPath(a ReadPathArgs, li *[]string, wg *sync.WaitGroup,
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
			go readPath(args, li, wg, mtx)
		} else {
			readPath(args, li, wg, mtx)
		}
		return
	}

	// if it's a directory, call func recursively
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
				go readPath(args, li, wg, mtx)
			} else {
				readPath(args, li, wg, mtx)
			}
		}
		return
	}
	if a.Parallel {
		mtx.Lock()
		*li = append(*li, a.FPath)
		mtx.Unlock()
	} else {
		*li = append(*li, a.FPath)
	}
}

// ReadPath crawls the file system from a specified path and creates a mapping
// SHA1 hashes to file paths
func ReadPath(args ReadPathArgs) Map {
	if args.Verbose {
		fmt.Print("Reading file tree ...\n\n")
	}
	fileList := make([]string, 0, 100)
	if args.Parallel {
		var wg sync.WaitGroup
		var mtx = sync.Mutex{}
		wg.Add(1)
		readPath(args, &fileList, &wg, &mtx)
		wg.Wait()
	} else {
		readPath(args, &fileList, nil, nil)
	}

	var fileHashes = make(Map)
	var wg sync.WaitGroup
	var mtx = sync.Mutex{}
	progressBar := pb.New(len(fileList))
	progressBar.SetMaxWidth(80)
	if args.Verbose {
		fmt.Println("\nCalculating file hashes ...")
		progressBar.Start()
	}

	maxGoRoutines := runtime.NumCPU() * cpuMultiplicator
	numGoRoutines := 0
	for _, s := range fileList {
		if args.Parallel {
			wg.Add(1)
			numGoRoutines++
			go func(fpath string) {
				hashStr := calcHash(fpath, &wg)
				mtx.Lock()
				fileHashes[hashStr] = append(fileHashes[hashStr], fpath)
				mtx.Unlock()
				if args.Verbose {
					progressBar.Increment()
				}
			}(s)
			if numGoRoutines >= maxGoRoutines {
				wg.Wait()
			}
		} else {
			hashStr := calcHash(s, nil)
			fileHashes[hashStr] = append(fileHashes[hashStr], s)
			if args.Verbose {
				progressBar.Increment()
			}
		}
	}
	if args.Parallel {
		wg.Wait()
	}
	progressBar.FinishPrint("Done")
	fileHashes.Sort()
	return fileHashes
}
