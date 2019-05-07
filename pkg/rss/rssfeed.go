package rss

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Item struct {
	XMLName xml.Name `xml:"item"`
	Title   string   `xml:"title"`
	PubDate string   `xml:"pubDate"`
	Guid    string   `xml:"guid"`
	Link    string   `xml:"link"`
	Image   struct {
		Href string `xml:"href,attr"`
	} `xml:"image"`
	Description string     `xml:"description"`
	Enclosure   *Enclosure `xml:"enclosure"`
	Content     string     `xml:"encoded"`
	Duration    string     `xml:"duration"`
	Explicit    string     `xml:"explicit"`
	Subtitle    string     `xml:"subtitle"`
	EpisodeType string     `xml:"episodeType"`
}

type Enclosure struct {
	Length int    `xml:"length,attr"`
	Type   string `xml:"type,attr"`
	Url    string `xml:"url,attr"`
}

func (item *Item) Download(path string, wg *sync.WaitGroup) (int64, error) {
	defer wg.Done()

	var url string
	enclosure := item.Enclosure
	if enclosure != nil {
		url = enclosure.Url
	} else {
		url = item.Link
	}

	urlSegments := strings.Split(url, "/")
	urlFileName := urlSegments[len(urlSegments)-1]

	LastIndexEnd := func(s, substr string) int {
		return strings.LastIndex(s, substr) + len(substr)
	}

	fileName := urlFileName[:LastIndexEnd(urlFileName, ".mp3")]

	file, err := os.Create(filepath.Join(filepath.Clean(path), fileName))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	return io.Copy(file, resp.Body)
}

func (item *Item) Date() time.Time {
	pubdate := item.PubDate
	t, err := time.Parse("Mon, _2 Jan 2006 15:04:05 -0700", pubdate)
	if err != nil {
		panic(err)
	}
	return t
}

func (item *Item) ToJson() (data []byte, err error) {
	data, err = json.MarshalIndent(item, "", "    ")
	return
}

func (item *Item) String() string {
	s, err := item.ToJson()
	if err != nil {
		return string(err.Error())
	}
	return string(s)
}

type feed struct {
	*itemScanner
	next chan *Item
}

func NewFeedFromURL(url string) (fd *feed) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	scanner := newItemScanner(resp.Body)
	fd = &feed{scanner, nil}
	fd.init()
	return fd
}

func NewFeedFromFile(path string) (fd *feed) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	scanner := newItemScanner(file)
	fd = &feed{scanner, nil}
	fd.init()
	return
}

func (feed *feed) init() {
	nextChan := make(chan *Item)
	feed.next = nextChan
	go feed.iterator()
}

func (feed *feed) iterator() {
	for feed.Scan() {
		var item Item
		err := xml.Unmarshal(feed.Bytes(), &item)
		if err != nil {
			break
		}
		feed.next <- &item
	}
	close(feed.next)
}

func (feed *feed) Next() chan *Item {
	if feed.next == nil {
		panic(errors.New("feed hasn't been initialised yet"))
	}
	return feed.next
}

func (feed *feed) DownloadN(path string, n int) {
	if feed.next == nil {
		panic(errors.New("feed hasn't been initialised yet"))
	}

	var downloaders sync.WaitGroup

	for n > 0 || n < 0 {
		item, ok := <-feed.next
		if !ok {
			break
		}
		downloaders.Add(1)
		go item.Download(path, &downloaders)
		n--
	}
	downloaders.Wait()
}

func (feed *feed) Download(path string) {
	feed.DownloadN(path, -1)
}

func (feed *feed) DownloadLatest(path string) {
	feed.DownloadN(path, 1)
}
