// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/jteeuwen/blackfriday"
)

// Connection represents a connection between a tag and a post.
type Connection struct {
	Tag  Tag
	Post *Post
}

// Site holds a site's posts, tags and templates.
type Site struct {
	Posts       []*Post            // List of site posts.
	Tags        []Tag              // List of unique tags referenced by posts.
	Connections []Connection       // Bindings, connecting a post to a given tag.
	templates   *template.Template // Tree of all site templates.
	Root        string             // Root path for the site.
}

// LoadSite loads a new set for the given root path.
func LoadSite(root string) (*Site, error) {
	s := new(Site)
	s.Root = root

	// Load templates.
	err := s.loadTemplates()
	if err != nil {
		return nil, err
	}

	// Load posts.
	err = s.loadPosts()
	if err != nil {
		return nil, err
	}

	return s, nil
}

// PostCount counts the number of posts associated with the given tag.
func (s *Site) PostCount(tag Tag) int {
	var count int

	name := string(tag)
	for _, c := range s.Connections {
		if strings.EqualFold(string(c.Tag), name) {
			count++
		}
	}

	return count
}

// FindPosts finds all posts associated with the given tag.
func (s *Site) FindPosts(tag Tag) []*Post {
	list := make([]*Post, 0, 2)
	name := string(tag)

	for _, c := range s.Connections {
		if strings.EqualFold(string(c.Tag), name) {
			list = append(list, c.Post)
		}
	}

	return list
}

// FindTags finds all tags associated with the given post index.
func (s *Site) FindTags(post *Post) []Tag {
	list := make([]Tag, 0, 2)

	for _, c := range s.Connections {
		if c.Post == post && !containsTag(list, c.Tag) {
			list = append(list, c.Tag)
		}
	}

	return list
}

func containsTag(list []Tag, tag Tag) bool {
	name := string(tag)
	for _, t := range list {
		if strings.EqualFold(string(t), name) {
			return true
		}
	}
	return false
}

// PostIndex returns the index for the given post.
// Returns -1 if it was not found.
func (s *Site) PostIndex(p *Post) int {
	for i, v := range s.Posts {
		if v == p {
			return i
		}
	}
	return -1
}

// TagIndex returns the index for the given tag.
// Returns -1 if it was not found. This performs a case-insensitive compare.
func (s *Site) TagIndex(t Tag) int {
	name := string(t)
	for i, tag := range s.Tags {
		if strings.EqualFold(string(tag), name) {
			return i
		}
	}
	return -1
}

// Render renders a page using the specified template.
// Output is written to the given writer.
func (s *Site) Render(w io.Writer, name string, page interface{}) error {
	return s.templates.ExecuteTemplate(w, name, page)
}

// loadPosts loads all posts.
func (s *Site) loadPosts() error {
	path := filepath.Join(s.Root, "posts")
	return filepath.Walk(path, func(file string, stat os.FileInfo, err error) error {
		if err != nil || stat.IsDir() {
			return err
		}
		return s.loadPost(file)
	})
}

// Load loads post data from the given file.
// This also loads unique tags.
func (s *Site) loadPost(file string) error {
	// Read post data from file.
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	post := NewPost()

	// Check if we have meta data.
	data, tags, err := post.ReadMetadata(data)
	if err != nil {
		return err
	}

	// Parse content as markdown.
	post.Content = blackfriday.MarkdownCommon(data)

	// Add post to list.
	s.Posts = append(s.Posts, post)

	// Parse tags and create connection with current post.
	s.parseTags(tags, post)
	return nil
}

// parseTags reads tags from the given string.
// Tags are added to the set's tag list, provided they are unique.
//
// Additionally, it creates a binding between a tag and the given post.
func (s *Site) parseTags(value string, post *Post) {
	names := toList(value)
	if len(names) == 0 {
		return
	}

	for _, name := range names {
		tagIndex := s.TagIndex(Tag(name))

		if tagIndex == -1 {
			tagIndex = len(s.Tags)
			s.Tags = append(s.Tags, Tag(strings.ToLower(name)))
		}

		s.Connections = append(s.Connections, Connection{
			Tag:  s.Tags[tagIndex],
			Post: post,
		})
	}
}

// loadTemplates loads all templates.
func (s *Site) loadTemplates() error {
	path := filepath.Join(s.Root, "templates")
	fd, err := os.Open(path)
	if err != nil {
		return err
	}

	defer fd.Close()

	files, err := fd.Readdirnames(-1)
	if err != nil {
		return err
	}

	for i := range files {
		files[i] = filepath.Join(path, files[i])
	}

	s.templates, err = template.ParseFiles(files...)
	return err
}

// toList splits the given value and omits empty entries.
func toList(value string) []string {
	value = strings.TrimSpace(value)
	parts := strings.Split(value, ",")
	list := make([]string, 0, len(parts))

	for _, v := range parts {
		v = strings.TrimSpace(v)
		if len(v) > 0 {
			list = append(list, v)
		}
	}

	return list
}
