// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

// TagIndexPage represents a page, displaying all tags.
type TagIndexPage struct {
	*Page
	site *Site
}

// NewTagIndexPage returns a new TagIndexPage for the given site.
func NewTagIndexPage(site *Site) *TagIndexPage {
	p := new(TagIndexPage)
	p.Page = NewPage()
	p.Page.title = "Listing of tags"
	p.Page.description = p.Page.title
	p.Page.keywords = "tags, posts, archive, history, index"
	p.site = site
	return p
}

func (p *TagIndexPage) Tags() []Tag {
	TagsByName(p.site.Tags).Sort()
	return p.site.Tags
}

func (p *TagIndexPage) PostCount(tag Tag) int {
	return p.site.PostCount(tag)
}
