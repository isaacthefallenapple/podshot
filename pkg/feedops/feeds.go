package feedops

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

var feeds feed
var podcasts = feeds["podcasts"]
var path = `..\..\myfeeds.json`

type feed map[string][]podcast
type podcast map[string]string

func (f *feed) String() string {
	b, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return fmt.Sprintln(err)
	}
	return string(b)
}

func init() {
	if err := LoadJson(); err != nil {
		panic(err)
	}
}

func Feeds() *feed {
	return &feeds
}

func SetPath(p string) {
	path = filepath.Clean(p)
}

func Contains(name string) bool {
	for _, p := range podcasts {
		if p["name"] == name {
			return true
		}
	}
	return false
}

func Add(name, feed string) (b bool, err error) {
	if Contains(name) {
		return
	}
	feeds["podcasts"] = append(feeds["podcasts"], podcast{"name": name, "rss-feed": feed})
	return true, Commit()
}

func Remove(name string) (b bool, err error) {
	if !Contains(name) {
		return
	}
	for i, p := range podcasts {
		if p["name"] == name {
			feeds["podcasts"] = append(podcasts[:i], podcasts[i+1:]...)
			b = true
			break
		}
	}
	err = Commit()
	return
}

func Commit() (err error) {
	jsonData, err := json.Marshal(feeds)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = ioutil.WriteFile(path, jsonData, 0222); err != nil {
		fmt.Println(err)
		return
	}
	return LoadJson()
}

func LoadJson() (err error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(file, &feeds)
	if err != nil {
		fmt.Println(err)
		return
	}
	podcasts = feeds["podcasts"]
	return
}
