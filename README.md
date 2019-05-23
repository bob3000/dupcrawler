# dupcrawler

dupcrawler crawls the file system on a given path and outputs groups of
duplicate files.

## installation

    go get github.com/bob3000/dupcrawler

## usage

    dupcrawler - looking for duplicate files

    Usage:
      dupcrawler [options] <inputdir>
    
    Options:
      -e, --excludes=<paths>  Comma separated list of path to exclude [default: ""]
      -d, --depth=<i>         Set maximum file hierarchy depth [default: 0]
      -l, --symlinks          Follow symlinks
      -h --help               Show this screen
      -v, --verbose           Show whats being done
      --version               Show version
