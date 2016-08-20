// Command textclipping reads textClipping files on OS X.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/unixpickle/resfork"
	"github.com/unixpickle/resfork/textclipping"
)

var contentIDs = map[string]textclipping.ContentType{
	"utf8":   textclipping.UTF8Text,
	"rtf":    textclipping.RTF,
	"utf16":  textclipping.UTF16Text,
	"web":    textclipping.WebArchive,
	"ustyle": textclipping.UStyle,
	"style":  textclipping.Style,
}

func main() {
	var contentID string
	var noNewline bool
	flag.StringVar(&contentID, "type", "utf8",
		"content identifier (possibilities: utf8, rtf, utf16, web, ustyle, style)")
	flag.BoolVar(&noNewline, "nonewline", false, "don't print extra newline")
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "[flags] filepath")
		fmt.Fprintln(os.Stderr, "Available flags:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	contentType, ok := contentIDs[contentID]
	if !ok {
		fmt.Fprintln(os.Stderr, "Unknown content ID:", contentID)
		os.Exit(1)
	}

	r, err := resfork.Open(flag.Args()[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open resource fork:", err)
		os.Exit(1)
	}
	defer r.Close()
	clipping, err := textclipping.ReadTextClipping(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse resource fork:", err)
		os.Exit(1)
	}
	data := clipping.Data(contentType)
	if data != nil {
		io.Copy(os.Stdout, bytes.NewBuffer(data))
		if !noNewline {
			fmt.Println()
		}
	} else {
		fmt.Fprintln(os.Stderr, "No data with content ID:", contentID)
		os.Exit(1)
	}
}
