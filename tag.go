// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

// TagList represents a set of unique tags and all posts
// referring to them.
type TagList map[string][]*Post

// TagData describes some metadata for a tag.
type TagData struct {
	Name      string
	PostCount int
}

// TagIndex represents the tag index page.
type TagIndex struct {
	Keywords    []string
	Tags        []TagData
	Title       string
	Description template.HTMLAttr
	Lang        template.HTMLAttr
	Dir         template.HTMLAttr
}

// TagPage represents a single tag page.
type TagPage struct {
	Keywords    []string
	Posts       []*Post
	Title       string
	Description template.HTMLAttr
	Lang        template.HTMLAttr
	Dir         template.HTMLAttr
	Tag         string
}

// NewTagList returns a new, empty tag list.
func NewTagList() TagList {
	return make(TagList)
}

// GenerateTags generates output directories for all tag files
// and fills them with documents referencing the appropriate posts.
func (t TagList) Generate(path string, tmp *template.Template) {
	dst := filepath.Join(path, DeployDir)
	dst = filepath.Join(dst, TagsDir)

	err := os.MkdirAll(dst, 0744)
	test(err, "Generate tags")

	page := &TagIndex{
		Title:       "Listing of all tags",
		Description: template.HTMLAttr("Listing of all tags"),
		Lang:        template.HTMLAttr(DefaultLang),
		Dir:         template.HTMLAttr(DefaultDir),
	}

	for tag, posts := range t {
		t.save(dst, tag, posts, tmp)

		// Update tag index
		page.Tags = append(page.Tags, TagData{
			Name:      tag,
			PostCount: len(posts),
		})
	}

	// Generate tag index.
	dst = filepath.Join(dst, "index.html")
	fd, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, FilePermission)
	test(err, "Generate tag index")

	defer fd.Close()

	TagsByName(page.Tags).Sort()

	err = tmp.ExecuteTemplate(fd, "tagindex.html", page)
	test(err, "Generate tag index")
}

// save generates a tag page, listing all posts which reference
// the given tag.
func (t TagList) save(path, tag string, posts []*Post, tmp *template.Template) {
	page := &TagPage{
		Title:       fmt.Sprintf("[%s]", tag),
		Description: template.HTMLAttr(fmt.Sprintf("Listing of posts for tag: %s", tag)),
		Lang:        template.HTMLAttr(DefaultLang),
		Dir:         template.HTMLAttr(DefaultDir),
		Tag:         tag,
		Posts:       posts,
	}

	PostsByDate(posts).Sort()

	// Ensure output directory exists.
	path = filepath.Join(path, tag)
	err := os.MkdirAll(path, DirPermission)
	test(err, "Generate tag")

	// Prepare output buffer.
	path = filepath.Join(path, "index.html")
	fd, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, FilePermission)
	test(err, "Generate tag")

	defer fd.Close()

	// Run page thruogh template.
	err = tmp.ExecuteTemplate(fd, "tag.html", page)
	test(err, "Generate tag")
}

// Add maps the given tag to the specified post.
func (t TagList) Add(tag string, post *Post) {
	t[tag] = append(t[tag], post)
}
