// Package textclipping can parse the resource forks of
// Mac OS X textClipping files.
package textclipping

import (
	"encoding/binary"
	"fmt"
	"io"
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

		// TODO: figure out why the last block's size
		// doesn't match its content length.
		if err == io.ErrUnexpectedEOF {
			break
		}
	}
	return &t, nil
}

// Blocks returns the underlying list of data blocks
// in the text clipping document.
func (t *TextClipping) Blocks() [][]byte {
	return t.blocks
}
