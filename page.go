// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"html/template"
	"time"
)

// Page represents a single page to be rendered.
type Page struct {
	title       string
	description string
	keywords    string
	lang        string
	dir         string
	date        time.Time
}

// NewPage creates a new page with default settings.
func NewPage() *Page {
	p := new(Page)
	p.keywords = DefaultKeywords
	p.lang = DefaultLang
	p.dir = DefaultDir
	return p
}

// HasDate returns true if the post has a date defined.
func (p *Page) HasDate() bool { return !p.date.IsZero() }

// Date returns a rendered form of the post date.
func (p *Page) Date() string {
	return p.date.UTC().Format(TimeFormat)
}

// HasKeywords returns true if the post has keywords defined.
func (p *Page) HasKeywords() bool { return len(p.keywords) > 0 }

// Keywords returns a rendered form of the page keywords.
func (p *Page) Keywords() template.HTMLAttr {
	return template.HTMLAttr(p.keywords)
}

// Lang returns a rendered form of the page language.
func (p *Page) Lang() template.HTMLAttr {
	return template.HTMLAttr(p.lang)
}

// Dir returns a rendered form of the page direction.
func (p *Page) Dir() template.HTMLAttr {
	return template.HTMLAttr(p.dir)
}

// HasTitle returns true if the post has a title defined.
func (p *Page) HasTitle() bool { return len(p.title) > 0 }

// Title returns a rendered form of the page title.
func (p *Page) Title() string {
	return p.title
}

// HasDescription returns true if the post has a description defined.
func (p *Page) HasDescription() bool { return len(p.description) > 0 }

// Description returns a rendered form of the page description.
func (p *Page) Description() template.HTMLAttr {
	return template.HTMLAttr(p.description)
}
