// Package textclipping can parse the resource forks of
// Mac OS X textClipping files.
package textclipping

import (
	"encoding/binary"
	"fmt"
	"io"
)

// ContentType is an identifier indicating the kind of
// data stored in a section of a text clipping.
type ContentType [4]byte

var (
	UTF8Text = ContentType{'u', 't', 'f', '8'}
	RichText = ContentType{'R', 'T', 'F', ' '}
)

// TextClipping stores the contents of a text clipping's
// resource fork.
type TextClipping struct {
	blocks [][]byte
}

// ReadTextClipping reads and parses a text clipping's
// contents.
func ReadTextClipping(r io.Reader) (*TextClipping, error) {
	var t TextClipping
	for {
		var blockSize uint32
		err := binary.Read(r, binary.BigEndian, &blockSize)
		if err == io.EOF {
			break
		}
		if len(t.blocks) == 0 {
			if blockSize < 4 {
				return nil, fmt.Errorf("unexpected first size: %v", blockSize)
			}
			blockSize -= 4
		}
		block := make([]byte, int(blockSize))
		n, err := io.ReadFull(r, block)
		t.blocks = append(t.blocks, block[:n])

		if err == io.ErrUnexpectedEOF && blockSize == 0x100 {
			// TODO: figure out why the last block's size
			// doesn't match its content length.
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return &t, nil
}

// Blocks returns the underlying list of data blocks
// in the text clipping document.
func (t *TextClipping) Blocks() [][]byte {
	return t.blocks
}

// Types returns the available content types in this
// clipping.
func (t *TextClipping) Types() []ContentType {
	if len(t.blocks) < 2 {
		return nil
	}
	typeList := t.blocks[len(t.blocks)-2]
	var res []ContentType
	for i := 1; i < len(typeList)/16; i++ {
		var t ContentType
		copy(t[:], typeList[i*16:])
		res = append(res, t)
	}
	if len(res) > len(t.blocks)-2 {
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
			return t.blocks[i+1]
		}
	}
	return nil
}
