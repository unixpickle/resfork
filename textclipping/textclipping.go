// Package textclipping can extract data from the
// resource forks of Mac OS X textClipping files.
package textclipping

import (
	"encoding/binary"
	"errors"
	"io"
)

// ContentType is an identifier indicating the kind of
// data stored in a section of a text clipping.
type ContentType [4]byte

// These are common content types found in a typical
// textClipping resource fork.
// For the most part, these are placed in ascending
// order of obscurity.
var (
	UTF8Text   = ContentType{'u', 't', 'f', '8'}
	RTF        = ContentType{'R', 'T', 'F', ' '}
	UTF16Text  = ContentType{'u', 't', '1', '6'}
	WebArchive = ContentType{'w', 'e', 'b', 'a'}
	UStyle     = ContentType{'u', 's', 't', 'l'}
	Style      = ContentType{'s', 't', 'y', 'l'}
)

// TextClipping stores the contents of a text clipping's
// resource fork.
type TextClipping struct {
	header header
	footer []byte
	blocks [][]byte
}

// ReadTextClipping reads and parses a text clipping's
// contents.
func ReadTextClipping(r io.Reader) (*TextClipping, error) {
	var t TextClipping

	head, err := readHeader(r)
	if err != nil {
		return nil, err
	}
	t.header = head

	var totalRead int
	for totalRead < head.bodySize() {
		var blockSize uint32
		err := binary.Read(r, binary.BigEndian, &blockSize)
		if err != nil {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			return nil, err
		}
		block := make([]byte, int(blockSize))
		n, err := io.ReadFull(r, block)
		t.blocks = append(t.blocks, block[:n])
		if err != nil {
			return nil, err
		}
		totalRead += n + 4
	}
	if totalRead > head.bodySize() {
		return nil, errors.New("body was longer than expected")
	}

	t.footer = make([]byte, head.footerSize())
	if _, err := io.ReadFull(r, t.footer); err != nil {
		return nil, err
	}

	return &t, nil
}

// Types returns the available content types in this
// clipping.
func (t *TextClipping) Types() []ContentType {
	if len(t.blocks) == 0 {
		return nil
	}
	typeList := t.blocks[len(t.blocks)-1]
	var res []ContentType
	for i := 1; i < len(typeList)/16; i++ {
		var t ContentType
		copy(t[:], typeList[i*16:])
		res = append(res, t)
	}
	if len(res) > len(t.blocks)-1 {
		return nil
	}
	return res
}

// Data returns the data for the given content type.
// It returns nil if the type is unavailable.
func (t *TextClipping) Data(ct ContentType) []byte {
	types := t.Types()
	for i, x := range types {
		if x == ct {
			return t.blocks[i]
		}
	}
	return nil
}
