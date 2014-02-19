// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	path, debug := parseArgs()

	templates := LoadTemplates(path)

	tags := GeneratePosts(path, templates)
	tags.Generate(path, templates)

	CopyStatic(path)

	if debug {
		return
	}
}

// parseArgs processes command line options and
// returns the ones we are interested in.
func parseArgs() (string, bool) {
	flag.Usage = usage

	version := flag.Bool("version", false, "")
	debug := flag.Bool("debug", false, "")
	tags := flag.String("tags", "", "")
	keywords := flag.String("keywords", "", "")

	flag.StringVar(&DefaultLang, "lang", DefaultLang, "")
	flag.StringVar(&DefaultDir, "dir", DefaultDir, "")
	flag.Parse()

	if *version {
		fmt.Printf("%s\n", Version())
		os.Exit(0)
	}

	var path string

	if flag.NArg() == 0 {
		path, _ = os.Getwd()
	} else {
		path = flag.Arg(0)
	}

	path = ValidatePath(path)

	DefaultTags = toList(*tags)
	DefaultKeywords = toList(*keywords)

	return path, *debug
}

// usage prints usage information.
func usage() {
	fmt.Printf(`usage: %v [options] [<path>]

[output options]
  -lang=%s
    Default ISO language code to use. This can be overridden on a per-document
    basis with the 'lang' metadata key.

  -dir=%s
    Default text direction to use. This can be overridden on a per-document
    basis with the 'dir' metadata key.

  -tags=tag1,tag2,...,tagN
    Default, comma-separated list of tags to use. This can be overridden on a
    per-document basis with the 'tags' metadata key.

  -keywords=word1,word2,...,wordN
    Default, comma-separated list of keywords to use. This can be overridden on
    a per-document basis with the 'keywords' metadata key.

  -debug
    Generates output in debug mode. This means that the entire site will
    be regenerated, without compression of HTML, JS, CSS and PNG images.

[misc options]
  -version
    Displays version information.
`,
		os.Args[0], DefaultLang, DefaultDir)
}
