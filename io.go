// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/jteeuwen/blackfriday"
)

// WriteIndex writes the front page.
// This is a special version of a normal Post.
func WriteIndex(site *Site) error {
	path := filepath.Join(site.Root, "index.md")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	post := NewPost()

	// Check if we have meta data.
	data, _, err = post.ReadMetadata(data)
	if err != nil {
		return err
	}

	// Parse content as markdown.
	post.Content = blackfriday.MarkdownCommon(data)

	// Generate output.
	path = filepath.Join(site.Root, "deploy")
	path = filepath.Join(path, "index.html")

	page := NewPostPage(post)

	fd, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, FilePermission)
	if err != nil {
		return err
	}

	defer fd.Close()

	return site.Render(fd, "index.html", page)
}

// WriteTags generates tag documents.
// These contain listings for all posts referencing a given tag.
func WriteTags(site *Site) error {
	dst := filepath.Join(site.Root, "deploy")
	dst = filepath.Join(dst, "tags")

	// Write individual posts.
	for _, tag := range site.Tags {
		err := writeTag(dst, site, tag)
		if err != nil {
			return err
		}
	}

	// Write tag index.
	return writeTagIndex(dst, site)
}

// writeTagIndex renders the tag index page.
func writeTagIndex(path string, site *Site) error {
	path = filepath.Join(path, "index.html")
	page := NewTagIndexPage(site)

	fd, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, FilePermission)
	if err != nil {
		return err
	}

	defer fd.Close()

	return site.Render(fd, "tagindex.html", page)
}

// writeTag renders the given tag.
func writeTag(path string, site *Site, tag Tag) error {
	path = filepath.Join(path, string(tag))

	// Create directories where necessary.
	err := os.MkdirAll(path, DirPermission)
	if err != nil {
		return err
	}

	path = filepath.Join(path, "index.html")
	page := NewTagPage(tag, site)

	fd, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, FilePermission)
	if err != nil {
		return err
	}

	defer fd.Close()

	return site.Render(fd, "tag.html", page)
}

// WritePosts generates post documents.
func WritePosts(site *Site) error {
	dst := filepath.Join(site.Root, "deploy")
	dst = filepath.Join(dst, "posts")

	// Write individual posts.
	for _, post := range site.Posts {
		tags := site.FindTags(post)
		err := writePost(dst, site, post, tags)
		if err != nil {
			return err
		}
	}

	// Write post index.
	return writePostIndex(dst, site)
}

// writePostIndex renders the posts index page.
func writePostIndex(path string, site *Site) error {
	path = filepath.Join(path, "index.html")
	page := NewPostIndexPage(site)

	fd, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, FilePermission)
	if err != nil {
		return err
	}

	defer fd.Close()

	return site.Render(fd, "postindex.html", page)
}

// writePost renders the given post.
func writePost(path string, site *Site, post *Post, tags []Tag) error {
	dir, file := post.SafePath()
	path = filepath.Join(path, dir)

	// Create directories where necessary.
	err := os.MkdirAll(path, DirPermission)
	if err != nil {
		return err
	}

	path = filepath.Join(path, file)
	page := NewPostPage(post, tags...)

	fd, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, FilePermission)
	if err != nil {
		return err
	}

	defer fd.Close()

	return site.Render(fd, "post.html", page)
}

// CopyStatic copies static content from source to target directories.
func CopyStatic(site *Site) error {
	src := filepath.Join(site.Root, "static")
	dst := filepath.Join(site.Root, "deploy")

	return filepath.Walk(src, func(file string, stat os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		dst := strings.Replace(file, src, dst, 1)

		if stat.IsDir() {
			return os.MkdirAll(dst, DirPermission)
		}

		fs, err := os.Open(file)
		if err != nil {
			return err
		}

		defer fs.Close()

		fd, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, FilePermission)
		if err != nil {
			return err
		}

		defer fd.Close()

		_, err = io.Copy(fd, fs)
		return err
	})
}
