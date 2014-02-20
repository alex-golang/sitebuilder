// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"html/template"
)

// PostPage represents a page, displaying posts.
type PostPage struct {
	*Page
	tags    []Tag
	content []byte
}

// NewPostPage returns a new PostPage for the given post
// and tags.
func NewPostPage(post *Post, tags ...Tag) *PostPage {
	p := new(PostPage)
	p.Page = NewPage()
	p.Page.keywords = post.Keywords
	p.Page.title = post.Title
	p.Page.description = post.Description
	p.Page.lang = post.Lang
	p.Page.dir = post.Dir
	p.Page.date = post.Date
	p.content = post.Content
	p.tags = tags
	return p
}

func (p *PostPage) Content() template.HTML { return template.HTML(string(p.content)) }
func (p *PostPage) HasTags() bool          { return len(p.tags) > 0 }
func (p *PostPage) Tags() template.HTML    { return RenderTags(p.tags) }
