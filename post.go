// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jteeuwen/blackfriday"
	"github.com/jteeuwen/ini"
)

// TimeFormat represents the timestamp format in post metadata.
const TimeFormat = "2006-01-02 15:04 MST"

var (
	DefaultLang     = "en,en-GB"
	DefaultDir      = "ltr"
	DefaultTags     = []string{}
	DefaultKeywords = []string{}

	endmeta = []byte("$endmeta")
)

// Post represents a single document/post.
type Post struct {
	Keywords    []string
	Tags        []string
	Content     []byte
	InFile      string
	OutFile     string
	Title       string
	Description string
	Author      string
	Lang        string
	Dir         string
	Date        time.Time
}

// NewPost creates a new, empty post with default settings.
func NewPost() *Post {
	p := new(Post)
	p.Date = time.Now()
	p.Tags = DefaultTags
	p.Keywords = DefaultKeywords
	p.Lang = DefaultLang
	p.Dir = DefaultDir
	return p
}

// ReadPosts traverses the given directory and its children
// and returns each post as it is found through the returned channel.
func ReadPosts(root string) <-chan *Post {
	c := make(chan *Post)

	go func() {
		defer close(c)

		root = filepath.Join(root, "posts")
		err := filepath.Walk(root, func(file string, stat os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if stat.IsDir() {
				return nil
			}

			p := NewPost()
			err = p.Load(file)
			if err != nil {
				return err
			}

			c <- p
			return nil
		})

		test(err, "Read posts")
	}()

	return c
}

// RenderDate renders the date in a readable format.
func (p *Post) RenderDate() string {
	return p.Date.Format(TimeFormat)
}

// RenderContent renders the content in a readable format.
func (p *Post) RenderContent() template.HTML {
	return template.HTML(string(p.Content))
}

// RenderTags renders tags in a readable format.
func (p *Post) RenderTags() string {
	return strings.Join(p.Tags, ", ")
}

// RenderKeywords renders tags in a readable format.
func (p *Post) RenderKeywords() string {
	return strings.Join(p.Keywords, ", ")
}

// Load loads post data from the given file.
func (p *Post) Load(file string) error {
	p.InFile = file

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	// Extract metadata if applicable.
	index := bytes.Index(data, endmeta)
	if index > -1 {
		meta := bytes.TrimSpace(data[:index])
		data = data[index+len(endmeta):]

		err = p.parseMeta(meta)
		if err != nil {
			return err
		}
	}

	if len(p.Title) != 0 {
		_, p.Title = filepath.Split(file)
	}

	// Read post data.
	p.Content = blackfriday.MarkdownCommon(data)
	return nil
}

// parseMeta parses metadata and fills the post fields with it.
func (p *Post) parseMeta(data []byte) error {
	ini := ini.New()
	err := ini.LoadBytes(data)
	if err != nil {
		return err
	}

	section := ini.Section("")

	p.Title = section.S("title", p.Title)
	p.Description = section.S("description", p.Description)
	p.Author = section.S("author", p.Author)
	p.Lang = section.S("lang", p.Lang)
	p.Dir = section.S("dir", p.Dir)

	p.Tags = toList(section.S("tags", strings.Join(p.Tags, ", ")))
	p.Keywords = toList(section.S("keywords", strings.Join(p.Keywords, ", ")))

	p.Date, err = time.Parse(TimeFormat, section.S("postdates", ""))
	if err != nil {
		p.Date = time.Now()
	}

	return nil
}

// toList splits the input string and removes empty entries.
func toList(value string) []string {
	split := strings.Split(value, ",")
	list := make([]string, 0, len(split))

	for _, value = range split {
		value = strings.TrimSpace(value)
		if len(value) > 0 {
			list = append(list, value)
		}
	}

	return list
}
