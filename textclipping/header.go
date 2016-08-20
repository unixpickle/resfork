package textclipping

import (
	"encoding/binary"
	"errors"
	"io"
)

type header []byte

func readHeader(r io.Reader) (header, error) {
	var headerSize uint32
	if err := binary.Read(r, binary.BigEndian, &headerSize); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return nil, err
	}
	if headerSize < 0x10 {
		return nil, errors.New("header is too small")
	}
	header := make(header, int(headerSize)-4)
	_, err := io.ReadFull(r, header)
	if err != nil {
		return nil, err
	}
	return header, nil
}

func (h header) bodySize() int {
	return int(binary.BigEndian.Uint32(h[4:8]))
}

func (h header) footerSize() int {
	return int(binary.BigEndian.Uint32(h[8:12]))
}
