package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
)

const VERSION = "0.1.0"

// only enqueue the paths that matched regular expression
var rxOk *regexp.Regexp = nil

// output directory
var output string

type FileExtender struct {
	// will use the default implementation of all but Visit and Filter
	gocrawl.DefaultExtender
}

// override visit
func (x *FileExtender) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	// prepare the directory and file
	path := ctx.URL().Path
	if path == "/" {
		// continue
		return nil, true
	}

	path = output + strings.Replace(path, "/", string(os.PathSeparator), -1)
	dname, _ := filepath.Split(path)
	os.MkdirAll(dname, 0755)

	// save res.Body to file
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	body := buf.Bytes()

	err := ioutil.WriteFile(path, body, 0644)
	if err != nil {
		log.Println(path, "write error:", err)
	}

	// return nil and true - let gocrawl find the links
	return nil, true
}

// override filter
// returns true to visit or false to ignore the URL
func (x *FileExtender) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	return !isVisited && (rxOk == nil || rxOk.MatchString(ctx.NormalizedURL().String()))
}

func main() {
	log.Println("version 0.1.0")

	var rx string
	maxVisit := flag.Int("m", 0, "max. visit(s), default 0, unlimited")
	delay := flag.Int("d", 500, "delay in ms between each request, default 500")
	flag.StringVar(&output, "o", ".", "output directory, default current directory")
	flag.StringVar(&rx, "rx", "", "only enqueue the paths that matched regular expression")
	flag.Parse()

	if len(rx) > 0 {
		rxOk = regexp.MustCompile(rx)
	}

	url := flag.Args()

	// set custom options
	opts := gocrawl.NewOptions(new(FileExtender))
	opts.RobotUserAgent = "Googlebot"
	opts.UserAgent = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
	opts.CrawlDelay = time.Duration(*delay) * time.Millisecond
	opts.LogFlags = gocrawl.LogAll
	opts.MaxVisits = *maxVisit

	// create crawler and start
	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run(url)
}
