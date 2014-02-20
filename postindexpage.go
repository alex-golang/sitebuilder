// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"html/template"
)

// PostIndexPage represents a page, displaying all posts.
type PostIndexPage struct {
	*Page
	site *Site
}

// NewPostIndexPage returns a new PostIndexPage for the given site.
func NewPostIndexPage(site *Site) *PostIndexPage {
	p := new(PostIndexPage)
	p.Page = NewPage()
	p.Page.title = "Listing of posts"
	p.Page.description = p.Page.title
	p.Page.keywords = "posts, archive, history, index"
	p.site = site
	return p
}

func (p *PostIndexPage) Posts() []*Post {
	PostsByDate(p.site.Posts).Sort()
	return p.site.Posts
}

func (p *PostIndexPage) PostDate(post *Post) string {
	return post.Date.Format(DateFormat)
}

func (p *PostIndexPage) PostDescription(post *Post) template.HTMLAttr {
	return template.HTMLAttr(post.Description)
}

func (p *PostIndexPage) PostTitle(post *Post) template.HTML {
	return template.HTML(post.Title)
}

func (p *PostIndexPage) PostPath(post *Post) template.HTMLAttr {
	return template.HTMLAttr(post.Path)
}

func (p *PostIndexPage) PostHasTags(post *Post) bool {
	return len(p.site.FindTags(post)) > 0
}

func (p *PostIndexPage) PostTags(post *Post) template.HTML {
	return RenderTags(p.site.FindTags(post))
}
