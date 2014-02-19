// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
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

	endmeta  = []byte("$endmeta")
	reg_name = regexp.MustCompile(`[^a-zA-Z0-9-_]`)
)

// PostIndex represents the post index page.
type PostIndex struct {
	Keywords    []string
	Title       string
	Description template.HTMLAttr
	Lang        template.HTMLAttr
	Dir         template.HTMLAttr
	Posts       []*Post
}

// Post represents a single document/post.
type Post struct {
	Keywords    []string
	Tags        []string
	Content     []byte
	RelPath     string
	Title       string
	Author      string
	Description template.HTMLAttr
	Lang        template.HTMLAttr
	Dir         template.HTMLAttr
	Date        time.Time
}

// NewPost creates a new, empty post with default settings.
func NewPost() *Post {
	p := new(Post)
	p.Date = time.Now()
	p.Tags = DefaultTags
	p.Keywords = DefaultKeywords
	p.Lang = template.HTMLAttr(DefaultLang)
	p.Dir = template.HTMLAttr(DefaultDir)
	return p
}

// GeneratePosts generates all posts.
// It returns a mapping of all unique tags found in the posts.
func GeneratePosts(path string, templ *template.Template) TagList {
	tags := NewTagList()
	posts := readPosts(path)

	dst := filepath.Join(path, DeployDir)
	dst = filepath.Join(dst, PostsDir)

	err := os.MkdirAll(dst, DirPermission)
	test(err, "Generate posts")

	page := &PostIndex{
		Title:       "Listing of all posts",
		Description: template.HTMLAttr("Listing of all posts"),
		Lang:        template.HTMLAttr(DefaultLang),
		Dir:         template.HTMLAttr(DefaultDir),
	}

	for post := range posts {
		post.save(dst, templ)

		// Update tag mapping.
		for _, tag := range post.Tags {
			tags.Add(tag, post)
		}

		// Update post index.
		page.Posts = append(page.Posts, post)
	}

	// Generate post index.
	dst = filepath.Join(dst, "index.html")
	fd, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, FilePermission)
	test(err, "Generate post index")

	defer fd.Close()

	PostsByDate(page.Posts).Sort()

	err = templ.ExecuteTemplate(fd, "postindex.html", page)
	test(err, "Generate post index")
	return tags
}

// save writes templated post data to the appropriate file.
func (p *Post) save(path string, templ *template.Template) {
	// Determine output file name.
	out := uniqueFilename(p.Title, p.Date)
	p.RelPath = strings.Join([]string{"", PostsDir, out}, "/")
	dst := filepath.Join(path, out)

	// Ensure output directory exists.
	err := os.MkdirAll(dst, DirPermission)
	test(err, "Generate post")

	// Prepare output buffer.
	dst = filepath.Join(dst, "index.html")
	fd, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, FilePermission)
	test(err, "Generate post")

	defer fd.Close()

	// Run post thruogh template.
	err = templ.ExecuteTemplate(fd, "post.html", p)
	test(err, "Generate post")
}

// uniqueFilename creates a unique file name from the
// given name and post date. The date component will contain
// the year, month, day, hour and minute:
//
//    <title>-<yyyymmddhhmm>.html
//
// The uniqueness is therefore only guaranteed up to the point
// where no two posts with the exact same title are made in the same
// minute.
func uniqueFilename(name string, date time.Time) string {
	name = strings.ToLower(name)
	name = strings.TrimSpace(name)
	name = strings.Replace(name, " ", "-", -1)
	name = reg_name.ReplaceAllString(name, "")
	name = fmt.Sprintf("%s-%s", name, date.Format("200601021504"))

	// remove duplicate --
	for strings.Index(name, "--") > -1 {
		name = strings.Replace(name, "--", "-", -1)
	}

	return name
}

// ReadPosts traverses the given directory and its children
// and returns each post as it is found through the returned channel.
func readPosts(root string) <-chan *Post {
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
	return p.Date.UTC().Format(TimeFormat)
}

// RenderContent renders the content in a readable format.
func (p *Post) RenderContent() template.HTML {
	return template.HTML(string(p.Content))
}

// RenderTags renders tags in a readable format.
func (p *Post) RenderTags() template.HTML {
	if len(p.Tags) == 0 {
		return ""
	}

	out := make([]string, len(p.Tags))

	for i, v := range p.Tags {
		out[i] = fmt.Sprintf(`<a href="/tag/%s">%s</a>`,
			template.HTMLEscapeString(v), template.HTMLEscapeString(v))
	}

	return template.HTML(strings.Join(out, ", "))
}

// RenderKeywords renders tags in a readable format.
func (p *Post) RenderKeywords() string {
	return strings.Join(p.Keywords, ", ")
}

// Load loads post data from the given file.
func (p *Post) Load(file string) error {
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

	if len(p.Title) == 0 {
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
	p.Author = section.S("author", p.Author)
	p.Description = template.HTMLAttr(section.S("description", string(p.Description)))
	p.Lang = template.HTMLAttr(section.S("lang", string(p.Lang)))
	p.Dir = template.HTMLAttr(section.S("dir", string(p.Dir)))

	p.Tags = toList(section.S("tags", strings.Join(p.Tags, ", ")))
	p.Keywords = toList(section.S("keywords", strings.Join(p.Keywords, ", ")))

	p.Date, err = time.Parse(TimeFormat, section.S("postdate", ""))
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
