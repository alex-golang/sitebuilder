// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"html/template"
	"strings"
)

// Tag represents a tag (surprise!).
type Tag string

// RenderTags renders the given set of tags.
func RenderTags(tags []Tag) template.HTML {
	if len(tags) == 0 {
		return template.HTML("")
	}

	html := make([]string, 0, len(tags))

	for _, tag := range tags {
		html = append(html,
			fmt.Sprintf(`<a href="/tags/%s/" title="Other posts in tag: %s">%s</a>`,
				tag, tag, tag))
	}

	return template.HTML(strings.Join(html, ", "))
}
