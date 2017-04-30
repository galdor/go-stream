//
// Copyright (c) 2017 Nicolas Martyanoff <khaelin@gmail.com>
//
// Permission to use, copy, modify, and distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

package stream

import (
	"bytes"
	"io"
)

func isEOF(err error) bool {
	return err == io.EOF || err == io.ErrUnexpectedEOF
}

type Stream struct {
	Reader io.Reader
	Buf    []byte
}

func NewStream(r io.Reader) *Stream {
	return &Stream{
		Reader: r,
		Buf:    []byte{},
	}
}

func NewStreamBytes(data []byte) *Stream {
	return NewStream(bytes.NewReader(data))
}

func (s *Stream) IsEmpty() (bool, error) {
	if len(s.Buf) > 0 {
		return false, nil
	}

	_, err := s.Peek(1)
	if err != nil {
		if isEOF(err) {
			return true, nil
		}

		return false, err
	}

	return false, nil
}

func (s *Stream) Peek(n int) ([]byte, error) {
	if len(s.Buf) < n {
		blen := len(s.Buf)
		rest := n - blen
		s.Buf = append(s.Buf, make([]byte, rest)...)

		nbRead, err := io.ReadAtLeast(s.Reader, s.Buf[blen:], rest)
		if err != nil {
			s.Buf = s.Buf[0 : blen+nbRead]
			return nil, err
		}
	}

	return dupBytes(s.Buf[0:n]), nil
}

func (s *Stream) PeekUpTo(n int) ([]byte, error) {
	_, err := s.Peek(n)
	if err != nil && !isEOF(err) {
		return nil, err
	}

	if n >= len(s.Buf) {
		n = len(s.Buf)
	}

	return dupBytes(s.Buf[0:n]), nil
}

func (s *Stream) StartsWithBytes(bs []byte) (bool, error) {
	data, err := s.Peek(len(bs))
	if err != nil {
		return false, err
	}

	return bytes.Equal(data, bs), nil
}

func (s *Stream) StartsWithByte(b byte) (bool, error) {
	data, err := s.Peek(1)
	if err != nil {
		return false, err
	}

	return data[0] == b, nil
}

func (s *Stream) Skip(n int) error {
	_, err := s.Peek(n)
	if err != nil {
		return err
	}

	s.Buf = s.Buf[n:]
	return nil
}

func (s *Stream) SkipBytes(bs []byte) (bool, error) {
	data, err := s.Peek(len(bs))
	if err != nil {
		return false, err
	}

	if bytes.Equal(data, bs) {
		s.Buf = s.Buf[len(bs):]
		return true, nil
	}

	return false, nil
}

func (s *Stream) SkipByte(b byte) (bool, error) {
	return s.SkipBytes([]byte{b})
}

func (s *Stream) SkipWhile(fn func(byte) bool) error {
	_, err := s.ReadWhile(fn)
	return err
}

func (s *Stream) Read(n int) ([]byte, error) {
	data, err := s.Peek(n)
	if err != nil {
		return nil, err
	}

	s.Buf = s.Buf[n:]
	return data, nil
}

func (s *Stream) ReadAll() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})

loop:
	for {
		data, err := s.PeekUpTo(4096)
		if err != nil {
			return nil, err
		}

		buf.Write(data)
		s.Buf = s.Buf[len(data):]

		empty, err := s.IsEmpty()
		if err != nil {
			return nil, err
		} else if empty {
			break loop
		}
	}

	return buf.Bytes(), nil
}

func (s *Stream) ReadWhile(fn func(byte) bool) ([]byte, error) {
	const blockSize = 32

	data := []byte{}

loop:
	for {
		block, err := s.PeekUpTo(blockSize)
		if err != nil {
			return nil, err
		} else if len(block) == 0 {
			break loop
		}

		for i, b := range block {
			if !fn(b) {
				s.Buf = s.Buf[i:]
				break loop
			}

			data = append(data, b)
		}

		s.Buf = s.Buf[len(block):]
	}

	return data, nil
}

func (s *Stream) PeekUntil(delim []byte) ([]byte, error) {
	for {
		idx := bytes.Index(s.Buf, delim)
		if idx >= 0 {
			return dupBytes(s.Buf[0:idx]), nil
		}

		blen := len(s.Buf)
		s.Buf = append(s.Buf, make([]byte, 4096)...)

		nread, err := s.Reader.Read(s.Buf[blen:])
		s.Buf = s.Buf[0 : blen+nread]
		if err != nil {
			if isEOF(err) {
				err = nil
			}
			return nil, err
		}
	}
}

func (s *Stream) PeekUntilByte(delim byte) ([]byte, error) {
	return s.PeekUntil([]byte{delim})
}

func (s *Stream) ReadUntil(delim []byte) ([]byte, error) {
	data, err := s.PeekUntil(delim)
	if err != nil || data == nil {
		return nil, err
	}

	s.Buf = s.Buf[len(data):]
	return data, nil
}

func (s *Stream) ReadUntilAndSkip(delim []byte) ([]byte, error) {
	data, err := s.PeekUntil(delim)
	if err != nil || data == nil {
		return nil, err
	}

	s.Buf = s.Buf[len(data)+len(delim):]
	return data, nil
}

func (s *Stream) ReadUntilByte(delim byte) ([]byte, error) {
	return s.ReadUntil([]byte{delim})
}

func (s *Stream) ReadUntilByteAndSkip(delim byte) ([]byte, error) {
	return s.ReadUntilAndSkip([]byte{delim})
}

func dupBytes(data []byte) []byte {
	ndata := make([]byte, len(data))
	copy(ndata, data)
	return ndata
}
