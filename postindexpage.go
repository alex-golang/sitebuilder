// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"html/template"
)

type PostIndexEntry struct {
	Title       template.HTML
	Description template.HTMLAttr
	Path        template.HTMLAttr
}

type PostIndex struct {
	Year  int
	Posts []*PostIndexEntry
}

// PostIndexPage represents a page, displaying all posts.
type PostIndexPage struct {
	*Page
	Years []*PostIndex
}

// NewPostIndexPage returns a new PostIndexPage for the given site.
func NewPostIndexPage(site *Site) *PostIndexPage {
	p := new(PostIndexPage)
	p.Page = NewPage()
	p.Page.title = "Listing of posts"
	p.Page.description = p.Page.title
	p.Page.keywords = "posts, archive, history, index"

	if len(site.Posts) == 0 {
		return p
	}

	p.Years = make([]*PostIndex, 0, len(site.Posts))

	PostsByDate(site.Posts).Sort()

	for _, post := range site.Posts {
		index := p.getIndex(post.Date.Year())

		index.Posts = append(index.Posts, &PostIndexEntry{
			Title:       template.HTML(post.Title),
			Description: template.HTMLAttr(post.Description),
			Path:        template.HTMLAttr(post.Path),
		})
	}

	return p
}

func (p *PostIndexPage) getIndex(year int) *PostIndex {
	for _, index := range p.Years {
		if index.Year == year {
			return index
		}
	}

	index := &PostIndex{Year: year}
	p.Years = append(p.Years, index)
	return index
}
