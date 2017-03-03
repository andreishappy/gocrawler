# Crawler two ways

### Overview

This project is an exploration of go that culminated in the implementation of a web crawler in two ways:

1. Using synchronization
  - uses a common map protected with a lock to build the result (this is an example of "Communicating by sharing memory")
2. Using channels
  - uses channels to synchronize goroutines in line with the "Share memory by communicating" principle of go
  - through experimentation I discovered that using unbuffered channels causes goroutines to wait unnecessarily, so a buffer of 10 was added to all channels

### Running tests

To run all the tests

```
go test andrei/...
```

### Running the crawler

To run the crawler 

```
go run src/andrei/gocrawl/gocrawl.go <mode> <url>
```

`<mode>` options:
1. `sync` runs the crawler using synchronization
2. `channel` runs the crawler using channels
3. `both` runs both

The results of crawling the website are saved to json files. I thought it's a nicer output than printing to the console and can easily be ingested by a front-end application to show the site map

If you prefer you can install 
```
go install andrei/gocrawl
```
And then run the binary in the bin directory with the same arguments as above.
