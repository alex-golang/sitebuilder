// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jteeuwen/ini"
)

const (
	// TimeFormat represents the timestamp format in post metadata.
	TimeFormat = "2006-01-02 15:04 MST"

	// DateFormat represents a rendered date
	DateFormat = "Jan _2, 2006"
)

var (
	// DefaultLang defines the default ISO language code which is to
	// be used for each generated page. This can be overridden by a
	// command line option and for each post individually through the
	// use of metadata information.
	DefaultLang = "en,en-GB"

	// DefaultDir defines the default text direction which is to
	// be used for each generated page. This can be overridden by a
	// command line option and for each post individually through the
	// use of metadata information.
	DefaultDir = "ltr"

	// DefaultTags defines the default set of tags to be used for each
	// generated page. This can be overridden by a command line option and for
	// each post individually through the use of metadata information.
	DefaultTags = ""

	// DefaultKeywords defines the default set of keywords to be used for each
	// generated page. This can be overridden by a command line option and for
	// each post individually through the use of metadata information.
	DefaultKeywords = ""

	endMeta = []byte("$endmeta")
	regName = regexp.MustCompile(`[^a-zA-Z0-9-_]`)
)

// Post represents a single document/post.
type Post struct {
	Content     []byte
	Keywords    string
	Title       string
	Path        string
	Description string
	Lang        string
	Dir         string
	Date        time.Time
}

// NewPost creates a new, empty post with default settings.
func NewPost() *Post {
	p := new(Post)
	p.Keywords = DefaultKeywords
	p.Lang = DefaultLang
	p.Dir = DefaultDir
	return p
}

// SafePath creates a relative file path using the post's title and post date.
//
//    yyyy/mm/dd/title.html
//
// The value is returned as the directory path and the file name.
func (p *Post) SafePath() (string, string) {
	dir := p.Date.Format("2006/01/02")
	file := strings.ToLower(p.Title)
	file = strings.Replace(file, " ", "-", -1)
	file = regName.ReplaceAllString(file, "")
	file = fmt.Sprintf("%s.html", file)

	// Trim duplicate -
	for strings.Index(file, "--") > -1 {
		file = strings.Replace(file, "--", "-", -1)
	}

	p.Path = strings.Join([]string{"", "posts", dir, file}, "/")
	return dir, file
}

// ReadMetadata reads post meta data from the given slice.
// It returns any remaining data and tags specified in the document.
func (p *Post) ReadMetadata(data []byte) ([]byte, string, error) {
	index := bytes.Index(data, endMeta)
	if index == -1 {
		return data, "", nil
	}

	ini := ini.New()
	err := ini.LoadBytes(data[:index])
	if err != nil {
		return data, "", err
	}

	section := ini.Section("") // Global section

	p.Title = section.S("title", p.Title)
	p.Description = section.S("description", p.Description)
	p.Keywords = section.S("keywords", p.Keywords)
	p.Lang = section.S("lang", p.Lang)
	p.Dir = section.S("dir", p.Dir)
	p.Date, err = time.Parse(TimeFormat,
		section.S("postdate", p.Date.Format(TimeFormat)))

	return data[index+len(endMeta):], section.S("tags", DefaultTags), err
}
