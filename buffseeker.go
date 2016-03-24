package main

import (
	"io"
)

// MessageBuffer for encrypting
type MessageBuffer struct {
	data   []byte
	offset int
}

// NewMessageBuffer creates new MessageBuffer (ByteBuffer with Seek)
func NewMessageBuffer(s string) *MessageBuffer {
	return &MessageBuffer{[]byte(s), 0}
}

func (b *MessageBuffer) Read(p []byte) (n int, err error) {
	delta := len(b.data) - b.offset
	if len(p) >= delta {
		for k, v := range b.data[b.offset:] {
			p[k] = v
		}
		b.offset = len(b.data)
		return delta, io.EOF
	}
	for k, v := range b.data[b.offset : b.offset+len(p)] {
		p[k] = v
	}
	b.offset += len(p)
	return len(p), nil
}

// Seek makes this a io.Seeker
func (b *MessageBuffer) Seek(offset int64, whence int) (int64, error) {
	var err error
	switch whence {
	case 0:
		b.offset = int(offset)
	case 2:
		b.offset = len(b.data)
		fallthrough
	case 1:
		b.offset += int(offset)
	}
	if int(b.offset) > len(b.data) {
		err = io.EOF
	}
	return int64(b.offset), err
}
