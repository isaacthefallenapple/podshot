// Package feedops provides types and operations to deal with the .json file storing subscribed-to feeds.
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

// feed defines the structure for the .json file.
type feed map[string][]podcast

// podcast defines the structure of a feed.
type podcast map[string]string //TODO refactor to a struct{ name, rss string }

func (f *feed) String() string {
	b, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return fmt.Sprintln(err)
	}
	return string(b)
}

// Feeds returns the global variable holding the feeds.
func Feeds() *feed {
	return &feeds
}

// SetPath sets the global path variable. path represents the path to the .json file.
func SetPath(p string) {
	path = filepath.Clean(p)
}

// Contains checks whether name is already in feeds.
func Contains(name string) bool {
	for _, p := range podcasts {
		if p["name"] == name {
			return true
		}
	}
	return false
}

// Add adds a new RSS feed. name specifies the name and feed the URL to the RSS feed.
// Returns false if name is already in feeds. Returns true and Commit's error otherwise.
func Add(name, feed string) (b bool, err error) {
	if Contains(name) {
		return
	}
	feeds["podcasts"] = append(feeds["podcasts"], podcast{"name": name, "rss-feed": feed})
	return true, Commit()
}

// Remove removes feed with name.
// Returns false if name is not in feeds. Returns true and Commit's error otherwise.
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


// Commit marshals feeds to json and writes it to a file. Reloads feeds with a call to LoadJson.
// Returns any errors that arise during marshaling, writing, or reloading.
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

// LoadJson reads the file specified in the global variable path and unmarshals it into feeds.
// Returns any errors that arise while reading or unmarshaling the file.
func LoadJson() (err error) { //TODO rename to LoadFeed or similar
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
