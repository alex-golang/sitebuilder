// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"sort"
)

type PostsByDate []*Post

func (p PostsByDate) Len() int           { return len(p) }
func (p PostsByDate) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PostsByDate) Less(i, j int) bool { return p[j].Date.After(p[i].Date) }
func (p PostsByDate) Sort()              { sort.Sort(p) }

type PostsByTitle []*Post

func (p PostsByTitle) Len() int           { return len(p) }
func (p PostsByTitle) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PostsByTitle) Less(i, j int) bool { return p[i].Title < p[j].Title }
func (p PostsByTitle) Sort()              { sort.Sort(p) }

type TagsByName []TagData

func (p TagsByName) Len() int           { return len(p) }
func (p TagsByName) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p TagsByName) Less(i, j int) bool { return p[i].Name < p[j].Name }
func (p TagsByName) Sort()              { sort.Sort(p) }
