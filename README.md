# CT Web Crawler 

### Introduction
A polite, slim and concurrent web crawler written in Go. 

### Dependencies
* [Google Go (golang)] (https://golang.org)
* [gocrawl] (https://github.com/PuerkitoBio/gocrawl)

### Build
```bash
$ go get github.com/PuerkitoBio/gocrawl
$ go build ctcrawler.go
```

### Usuage
```bash
$ ./ctcrawler http://example.com/index.html
```

### Know Issues / Future Development
* Authentication is not implemented and no plan for this
* URL query is not handled yet 

Enjoy! Welcome for forking and pull request.
