// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"html/template"
)

// TagPage represents a page, displaying posts.
type TagPage struct {
	*Page
	site *Site
	tag  Tag
}

// NewTagPage returns a new TagPage for the given post
// and tags.
func NewTagPage(tag Tag, site *Site) *TagPage {
	p := new(TagPage)
	p.Page = NewPage()
	p.Page.keywords = fmt.Sprintf("%s, tags, archive, posts, history", tag)
	p.Page.title = fmt.Sprintf("Posts in tag: %s", tag)
	p.Page.description = fmt.Sprintf("Listing of posts in tag: %s", tag)
	p.tag = tag
	p.site = site
	return p
}

func (p *TagPage) Tag() Tag { return p.tag }

func (p *TagPage) Posts(tag Tag) []*Post {
	posts := p.site.FindPosts(tag)
	PostsByDate(posts).Sort()
	return posts
}

func (p *TagPage) PostDate(post *Post) string {
	return post.Date.Format(DateFormat)
}

func (p *TagPage) PostDescription(post *Post) template.HTMLAttr {
	return template.HTMLAttr(post.Description)
}

func (p *TagPage) PostTitle(post *Post) template.HTML {
	return template.HTML(post.Title)
}

func (p *TagPage) PostPath(post *Post) template.HTMLAttr {
	return template.HTMLAttr(post.Path)
}
