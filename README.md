# dupcrawler

dupcrawler crawls the file system and outputs groups of duplicate files.

## installation

    go get github.com/bob3000/dupcrawler

## usage

    > dupcrawler -h
    dupcrawler - looking for duplicate files

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
      --version               Show version

    > dupcrawler -v photos/
    Reading file tree ...
    
    Examining photos ...
    Examining photos/2000 ...
    Examining photos/2018 ...
    Examining photos/2019 ...
    Examining photos/2017 ...
    Examining photos/2016 ...
    
    Calculating file hashes ...
     2626 / 2626 [===================================================] 100.00% 3m21s
    Done
    
    Duplicate files:
    
    photos/2017/VID_20171214_201502898.mp4
    photos/2017/VID_20171214_201502898_1.mp4