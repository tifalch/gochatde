# gochatde
[![Build Status](https://travis-ci.org/TibFalch/gochatde.svg?branch=master)](https://travis-ci.org/TibFalch/gochatde)

go version of chatde

## Usage
### Flags
```
gochatde [-color] [-gzip] [-check=false] [-debug] <IP>[:<port>]
    -color -colour
        adds coloration to the output

    -gzip
        gzip compresses all the data transfered over the wire

    -check=false
        disables checksum feature
        might be deprecated soon

    -debug
        enables debug output

    IP   = the ip to connect to
    port = port to connect to (default=15327)
```

### Editor Commands
```
 Δ §help
§bye, §quit
        end gochatde
§file <file>
        send encrypted file
§ls
        lists current directory
§cd <dir>
        changes directory
```

## Encryption
[https://github.com/LFalch/delta-l](delta-l)

## Rust Port
[https://github.com/TibFalch/chatde-rs](chatde-rs)
