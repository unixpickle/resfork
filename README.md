# resfork

This is a Go package to read and write files' [resource forks](https://en.wikipedia.org/wiki/Resource_fork). Resource forks are used for various things on OS X, such as metadata, text clippings, thumbnails, etc. It will also include sub-packages for common types of data stored in resource forks (e.g. text clippings).

# Contents

The root repository contains a Go package for opening and reading resource forks.

The [textclipping](textclipping) directory contains a Go package for processing Mac OS X textClippings.

The [textclipping/cmd](textclipping/cmd) directory contains a command-line utility for reading textClippings.
