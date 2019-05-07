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

// Item defines the struct to marshal RSS into.
type Item struct {
	XMLName xml.Name `xml:"item"`
	// Tile stores the item's title.
	Title   string   `xml:"title"`

	// PubDate stores the item's date of publication (as a string). It can be accessed as a time.Time through the Date() method.
	PubDate string   `xml:"pubDate"`

	// Guid stores the item's guid.
	Guid    string   `xml:"guid"`

	// Link stores the item's online location.
	Link    string   `xml:"link"`

	// Image stores the item's thumbnail.
	Image   struct {
		Href string `xml:"href,attr"`
	} `xml:"image"`

	// Description stores the item's description.
	Description string     `xml:"description"`

	// Enclosure stores stores the item's enclosure, and its associated fields.
	Enclosure   *Enclosure `xml:"enclosure"`

	// Content stores the item's description in another format.
	Content     string     `xml:"encoded"`

	// Duration stores the item's length (as a string).
	Duration    string     `xml:"duration"`

	// Explicit stores whether the item has explicit content.
	Explicit    string     `xml:"explicit"`

	// Subtitle stores the item's subtitle.
	Subtitle    string     `xml:"subtitle"`

	// EpisodeType stores the item's type.
	EpisodeType string     `xml:"episodeType"`
}

// An Enclosure describes an item's metadata
type Enclosure struct {
	// Lenght stores the item's length in bytes.
	Length int    `xml:"length,attr"`

	// Type stores the item's audio type (e.g. audio/mpeg).
	Type   string `xml:"type,attr"`

	// Url stores the item's online location.
	Url    string `xml:"url,attr"`
}

// Download requests item's URL (either found in its Enclosure or, if that's nil, its Link) and
// download it into the given path.
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

// Date return the item's publication date as a time.Time.
func (item *Item) Date() time.Time {
	pubdate := item.PubDate
	t, err := time.Parse("Mon, _2 Jan 2006 15:04:05 -0700", pubdate)
	if err != nil {
		panic(err)
	}
	return t
}

// ToJson marshals the item to json with json.MarshalIndent and returns the resulting byte-slice and error.
func (item *Item) ToJson() ([]byte, error) {
	return json.MarshalIndent(item, "", "    ")
}

func (item *Item) String() string {
	s, err := item.ToJson()
	if err != nil {
		return string(err.Error())
	}
	return string(s)
}

// feed defines a type to iterate over the items of a RSS feed.
type feed struct {
	*itemScanner
	next chan *Item
}

// NewFeedFromURL returns a new feed read from the given URL.
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

// NewFeedFromFile returns a new feed read from the given file path.
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

// iterator will begin to read the RSS feed item by item and and put the unmarshalled xml into
// the feed's next channel.
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

// Next returns a channel which will yield consecutive items of the feed when iterated over.
func (feed *feed) Next() chan *Item {
	if feed.next == nil {
		panic(errors.New("feed hasn't been initialised yet"))
	}
	return feed.next
}

// DownloadN downloads n consecutive items into the path. Each file will be named after the file name
// from the item's URL. if n is less than 0, all items will be downloaded.
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

// Download downloads all of feed's items into path.
func (feed *feed) Download(path string) {
	feed.DownloadN(path, -1)
}

// DownloadLatest downloads one item of feed into path.
func (feed *feed) DownloadLatest(path string) {
	feed.DownloadN(path, 1)
}
