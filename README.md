## sitebuild

Sitebuild is a simple static website generator.
It exists mostly for my own use, as it is employed to generate
the contents for my personal website.

It comes as a single command line tool, which accepts
the path to a target directory as an argument.

This directory holds the site contents in a specific layout:

    [$path]
      |- [posts]
      |   |- a.md
      |   |- b.md
      |   |- ...
      |
      |- [static]
      |   |-[css]
      |   |-[js]
      |   |-[img]
      |   |-[...]
      |
      |- [templates]
          |- foo.html
          |- bar.html
          |- baz.html

* **posts**: contains the actual post contents as Markdown (`.md`) files.
  The directory structure inside this dir can be anything you want.
* **static**: This directory holds static content which shoul be included
  in the site as-is. This includes things like images, stylesheets,
  javascripts, etc. The contents of this directory (including sub directories)
  is copied over 1:1.
* **templates**: This directory holds a set of templates, with syntax compatible
  with Go's `html/template` package. These are used to generate the actual
  site pages.


### Usage

    go get github.com/jteeuwen/sitebuild


### Documentation

Documentation can be found at [godoc.org](http://godoc.org/github.com/jteeuwen/sitebuild).


### License

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.

