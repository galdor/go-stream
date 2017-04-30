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
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertStreamEqual(t *testing.T, bytes []byte, stream *Stream) {
	data, err := ioutil.ReadAll(stream.Reader)
	if err != nil {
		t.Fatalf("cannot read stream reader: %v", err)
	}

	assert.Equal(t, bytes, append(stream.Buf, data...))
}

func TestStreamIsEmpty(t *testing.T) {
	assert := assert.New(t)

	var stream *Stream

	stream = NewStreamBytes([]byte{})
	assert.True(stream.IsEmpty())

	stream = NewStreamBytes([]byte{1, 2, 3})
	stream.Skip(3)
	assert.True(stream.IsEmpty())

	stream = NewStreamBytes([]byte{1, 2, 3})
	assert.False(stream.IsEmpty())
}

func TestStreamPeek(t *testing.T) {
	assert := assert.New(t)

	var stream *Stream
	var data []byte
	var err error

	stream = NewStreamBytes([]byte{})
	data, err = stream.Peek(3)
	assert.Error(err)

	stream = NewStreamBytes([]byte{})
	data, err = stream.Peek(0)
	assert.NoError(err)
	assert.Equal([]byte{}, data)

	stream = NewStreamBytes([]byte{1, 2, 3})
	data, err = stream.Peek(3)
	assert.NoError(err)
	assert.Equal([]byte{1, 2, 3}, data)

	stream = NewStreamBytes([]byte{1, 2, 3})
	data, err = stream.Peek(1)
	assert.NoError(err)
	assert.Equal([]byte{1}, data)
}

func TestStreamPeekUpTo(t *testing.T) {
	assert := assert.New(t)

	var stream *Stream
	var data []byte
	var err error

	stream = NewStreamBytes([]byte{})
	data, err = stream.PeekUpTo(3)
	assert.NoError(err)
	assert.Equal([]byte{}, data)

	stream = NewStreamBytes([]byte{})
	data, err = stream.PeekUpTo(0)
	assert.NoError(err)
	assert.Equal([]byte{}, data)

	stream = NewStreamBytes([]byte{1, 2, 3})
	data, err = stream.PeekUpTo(1)
	assert.NoError(err)
	assert.Equal([]byte{1}, data)

	stream = NewStreamBytes([]byte{1, 2, 3})
	data, err = stream.PeekUpTo(6)
	assert.NoError(err)
	assert.Equal([]byte{1, 2, 3}, data)
}

func TestStreamStartsWithBytes(t *testing.T) {
	assert := assert.New(t)

	var stream *Stream
	var res bool
	var err error

	stream = NewStreamBytes([]byte{})
	res, err = stream.StartsWithBytes([]byte{})
	assert.NoError(err)
	assert.True(res)

	stream = NewStreamBytes([]byte{})
	res, err = stream.StartsWithBytes([]byte{1})
	assert.Error(err)
	assert.False(res)

	stream = NewStreamBytes([]byte{1})
	res, err = stream.StartsWithBytes([]byte{1})
	assert.NoError(err)
	assert.True(res)

	stream = NewStreamBytes([]byte{1, 2, 3})
	res, err = stream.StartsWithBytes([]byte{1})
	assert.NoError(err)
	assert.True(res)

	stream = NewStreamBytes([]byte{1, 2, 3})
	res, err = stream.StartsWithBytes([]byte{1, 2, 3})
	assert.NoError(err)
	assert.True(res)

	stream = NewStreamBytes([]byte{1, 2, 3})
	res, err = stream.StartsWithBytes([]byte{4, 5})
	assert.NoError(err)
	assert.False(res)
}

func TestStreamStartsWithByte(t *testing.T) {
	assert := assert.New(t)

	var stream *Stream
	var res bool
	var err error

	stream = NewStreamBytes([]byte{})
	res, err = stream.StartsWithByte(1)
	assert.Error(err)
	assert.False(res)

	stream = NewStreamBytes([]byte{1})
	res, err = stream.StartsWithByte(1)
	assert.NoError(err)
	assert.True(res)

	stream = NewStreamBytes([]byte{1, 2, 3})
	res, err = stream.StartsWithByte(1)
	assert.NoError(err)
	assert.True(res)

	stream = NewStreamBytes([]byte{1, 2, 3})
	res, err = stream.StartsWithByte(4)
	assert.NoError(err)
	assert.False(res)
}

func TestStreamSkip(t *testing.T) {
	assert := assert.New(t)

	var stream *Stream

	stream = NewStreamBytes([]byte{})
	assert.NoError(stream.Skip(0))
	assertStreamEqual(t, []byte{}, stream)

	stream = NewStreamBytes([]byte{})
	assert.Error(stream.Skip(1))
	assertStreamEqual(t, []byte{}, stream)

	stream = NewStreamBytes([]byte{1, 2, 3})
	assert.NoError(stream.Skip(0))
	assertStreamEqual(t, []byte{1, 2, 3}, stream)

	stream = NewStreamBytes([]byte{1, 2, 3})
	assert.NoError(stream.Skip(1))
	assertStreamEqual(t, []byte{2, 3}, stream)

	stream = NewStreamBytes([]byte{1, 2, 3})
	assert.NoError(stream.Skip(3))
	assertStreamEqual(t, []byte{}, stream)

	stream = NewStreamBytes([]byte{1, 2, 3})
	assert.Error(stream.Skip(6))
	assertStreamEqual(t, []byte{1, 2, 3}, stream)
}

func TestStreamSkipBytes(t *testing.T) {
	assert := assert.New(t)

	var stream *Stream
	var res bool
	var err error

	stream = NewStreamBytes([]byte{})
	res, err = stream.SkipBytes([]byte{})
	assert.NoError(err)
	assert.True(res)
	assertStreamEqual(t, []byte{}, stream)

	stream = NewStreamBytes([]byte{1, 2, 3})
	res, err = stream.SkipBytes([]byte{})
	assert.NoError(err)
	assert.True(res)
	assertStreamEqual(t, []byte{1, 2, 3}, stream)

	stream = NewStreamBytes([]byte{1, 2, 3})
	res, err = stream.SkipBytes([]byte{1, 2})
	assert.NoError(err)
	assert.True(res)
	assertStreamEqual(t, []byte{3}, stream)

	stream = NewStreamBytes([]byte{1, 2, 3})
	res, err = stream.SkipBytes([]byte{1, 2, 3})
	assert.NoError(err)
	assert.True(res)
	assertStreamEqual(t, []byte{}, stream)

	stream = NewStreamBytes([]byte{1, 2, 3})
	res, err = stream.SkipBytes([]byte{4, 5})
	assert.NoError(err)
	assert.False(res)
	assertStreamEqual(t, []byte{1, 2, 3}, stream)

	stream = NewStreamBytes([]byte{1, 2, 3})
	res, err = stream.SkipBytes([]byte{1, 2, 3, 4})
	assert.Error(err)
	assert.False(res)
	assertStreamEqual(t, []byte{1, 2, 3}, stream)
}

func TestStreamSkipByte(t *testing.T) {
	assert := assert.New(t)

	var stream *Stream
	var res bool
	var err error

	stream = NewStreamBytes([]byte{})
	res, err = stream.SkipByte(1)
	assert.Error(err)
	assert.False(res)
	assertStreamEqual(t, []byte{}, stream)

	stream = NewStreamBytes([]byte{1, 2, 3})
	res, err = stream.SkipByte(1)
	assert.NoError(err)
	assert.True(res)
	assertStreamEqual(t, []byte{2, 3}, stream)

	stream = NewStreamBytes([]byte{1, 2, 3})
	res, err = stream.SkipByte(4)
	assert.NoError(err)
	assert.False(res)
	assertStreamEqual(t, []byte{1, 2, 3}, stream)
}

func TestStreamRead(t *testing.T) {
	assert := assert.New(t)

	var stream *Stream
	var data []byte
	var err error

	stream = NewStreamBytes([]byte{})
	data, err = stream.Read(0)
	assert.NoError(err)
	assert.Equal([]byte{}, data)
	assertStreamEqual(t, []byte{}, stream)

	stream = NewStreamBytes([]byte{})
	data, err = stream.Read(3)
	assert.Error(err)
	assertStreamEqual(t, []byte{}, stream)

	stream = NewStreamBytes([]byte{1, 2, 3})
	data, err = stream.Read(2)
	assert.NoError(err)
	assert.Equal([]byte{1, 2}, data)
	assertStreamEqual(t, []byte{3}, stream)

	stream = NewStreamBytes([]byte{1, 2, 3})
	data, err = stream.Read(3)
	assert.NoError(err)
	assert.Equal([]byte{1, 2, 3}, data)
	assertStreamEqual(t, []byte{}, stream)

	stream = NewStreamBytes([]byte{1, 2})
	data, err = stream.Read(3)
	assert.Error(err)
	assertStreamEqual(t, []byte{1, 2}, stream)
}

func TestStreamReadWhile(t *testing.T) {
	assert := assert.New(t)

	var stream *Stream
	var data []byte
	var err error

	even := func(b byte) bool {
		return b%2 == 0
	}

	stream = NewStreamBytes([]byte{})
	data, err = stream.ReadWhile(even)
	assert.NoError(err)
	assert.Equal([]byte{}, data)
	assertStreamEqual(t, []byte{}, stream)

	stream = NewStreamBytes([]byte{2, 4, 5})
	data, err = stream.ReadWhile(even)
	assert.NoError(err)
	assert.Equal([]byte{2, 4}, data)
	assertStreamEqual(t, []byte{5}, stream)

	stream = NewStreamBytes([]byte{2, 4, 8})
	data, err = stream.ReadWhile(even)
	assert.NoError(err)
	assert.Equal([]byte{2, 4, 8}, data)
	assertStreamEqual(t, []byte{}, stream)

	stream = NewStreamBytes([]byte{1, 2})
	data, err = stream.ReadWhile(even)
	assert.NoError(err)
	assert.Equal([]byte{}, data)
	assertStreamEqual(t, []byte{1, 2}, stream)
}

func TestStreamPeekUntil(t *testing.T) {
	assert := assert.New(t)

	var stream *Stream
	var data []byte
	var err error

	stream = NewStreamBytes([]byte{})
	data, err = stream.PeekUntil([]byte{})
	assert.NoError(err)
	assert.Equal([]byte{}, data)

	stream = NewStreamBytes([]byte{})
	data, err = stream.PeekUntil([]byte{2, 3})
	assert.NoError(err)
	assert.Nil(data)

	stream = NewStreamBytes([]byte{1, 2, 4})
	data, err = stream.PeekUntil([]byte{2, 3})
	assert.NoError(err)
	assert.Nil(data)

	stream = NewStreamBytes([]byte{1, 2, 4})
	data, err = stream.PeekUntil([]byte{1, 2})
	assert.NoError(err)
	assert.Equal([]byte{}, data)

	stream = NewStreamBytes([]byte{1, 2, 4})
	data, err = stream.PeekUntil([]byte{2, 4})
	assert.NoError(err)
	assert.Equal([]byte{1}, data)
}

func TestStreamReadUntil(t *testing.T) {
	assert := assert.New(t)

	var stream *Stream
	var data []byte
	var err error

	stream = NewStreamBytes([]byte{})
	data, err = stream.ReadUntil([]byte{})
	assert.NoError(err)
	assert.Equal([]byte{}, data)
	assertStreamEqual(t, []byte{}, stream)

	stream = NewStreamBytes([]byte{})
	data, err = stream.ReadUntil([]byte{2, 3})
	assert.NoError(err)
	assert.Nil(data)
	assertStreamEqual(t, []byte{}, stream)

	stream = NewStreamBytes([]byte{1, 2, 4})
	data, err = stream.ReadUntil([]byte{2, 3})
	assert.NoError(err)
	assert.Nil(data)
	assertStreamEqual(t, []byte{1, 2, 4}, stream)

	stream = NewStreamBytes([]byte{1, 2, 4})
	data, err = stream.ReadUntil([]byte{1, 2})
	assert.NoError(err)
	assert.Equal([]byte{}, data)
	assertStreamEqual(t, []byte{1, 2, 4}, stream)

	stream = NewStreamBytes([]byte{1, 2, 4})
	data, err = stream.ReadUntil([]byte{2, 4})
	assert.NoError(err)
	assert.Equal([]byte{1}, data)
	assertStreamEqual(t, []byte{2, 4}, stream)
}

func TestStreamReadUntilAndSkip(t *testing.T) {
	assert := assert.New(t)

	var stream *Stream
	var data []byte
	var err error

	stream = NewStreamBytes([]byte{})
	data, err = stream.ReadUntilAndSkip([]byte{})
	assert.NoError(err)
	assert.Equal([]byte{}, data)
	assertStreamEqual(t, []byte{}, stream)

	stream = NewStreamBytes([]byte{})
	data, err = stream.ReadUntilAndSkip([]byte{2, 3})
	assert.NoError(err)
	assert.Nil(data)
	assertStreamEqual(t, []byte{}, stream)

	stream = NewStreamBytes([]byte{1, 2, 4})
	data, err = stream.ReadUntilAndSkip([]byte{2, 3})
	assert.NoError(err)
	assert.Nil(data)
	assertStreamEqual(t, []byte{1, 2, 4}, stream)

	stream = NewStreamBytes([]byte{1, 2, 4})
	data, err = stream.ReadUntilAndSkip([]byte{1, 2})
	assert.NoError(err)
	assert.Equal([]byte{}, data)
	assertStreamEqual(t, []byte{4}, stream)

	stream = NewStreamBytes([]byte{1, 2, 4})
	data, err = stream.ReadUntilAndSkip([]byte{2, 4})
	assert.NoError(err)
	assert.Equal([]byte{1}, data)
	assertStreamEqual(t, []byte{}, stream)
}
