/*
MIT License

Copyright (c) 2017 Jumbo
Author: Jumbo
Email: Jumbo.Wu@hotmail.com

File: buffer.go

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

/*
  Buffer data struct , Writer and Reader
*/

// buffer
package packet

import (
	"errors"
	"math"

	//log "github.com/sirupsen/logrus"
)

type Buffer struct {
	pos  int
	data []byte
}

func (b *Buffer) Data() []byte {
	return b.data
}

func (b *Buffer) Length() int {
	return len(b.data)
}

func (b *Buffer) Used() int {
	return b.pos
}

func (b *Buffer) Left() int {
	return b.Length() - b.Used()
}

func (b *Buffer) SetPos(p int) {
	b.pos = p
}

//输入参数：读取的buffer
func Reader(data []byte) *Buffer {
	return &Buffer{data: data}
}

//输入参数：写入的buffer
func Writer(data []byte) *Buffer {
	return &Buffer{data: data}
}

//writer
func (b *Buffer) WriteZeros(n int) (err error) {

	if b.pos+n > len(b.data) {
		err = errors.New("write zeros failed")
		return
	}

	for i := 0; i < n; i++ {
		b.data[b.pos] = byte(0)
		b.pos++
	}

	return
}

func (b *Buffer) WriteBool(v bool) (err error) {

	if b.pos+1 > len(b.data) {
		err = errors.New("write bool failed")
		return
	}

	if v {
		b.data[b.pos] = byte(1)
	} else {
		b.data[b.pos] = byte(0)
	}

	b.pos++

	return
}

func (b *Buffer) WriteByte(v byte) (err error) {

	if b.pos+1 > len(b.data) {
		err = errors.New("write byte failed")
		return
	}

	b.data[b.pos] = v
	b.pos++

	return
}

func (b *Buffer) WriteBytes(v []byte) (err error) {

	size := len(v)
	err = b.WriteU16(uint16(size))
	if err != nil {
		return
	}

	if b.pos+size > len(b.data) {
		err = errors.New("write bytes failed")
		return
	}

	for i := 0; i < size; i++ {
		b.data[b.pos] = v[i]
		b.pos++
	}

	return
}

func (b *Buffer) WriteString(v string) (err error) {
	bytes := []byte(v)
	err = b.WriteBytes(bytes)

	if err != nil {
		err = errors.New("write string failed")
		return
	}

	return
}

func (b *Buffer) WriteS8(v int8) (err error) {
	err = b.WriteByte(byte(v))
	if err != nil {
		err = errors.New("write S8 failed")
		return
	}

	return
}

func (b *Buffer) WriteU16(v uint16) (err error) {
	if b.pos+2 > len(b.data) {
		err = errors.New("write U16 failed")
		return
	}

	for i := 0; i < 2; i++ {
		b.data[b.pos] = byte(v >> ((1 - uint(i)) * 8))
		b.pos++
	}

	return
}

func (b *Buffer) WriteS16(v int16) (err error) {
	err = b.WriteU16(uint16(v))
	if err != nil {
		err = errors.New("write S16 failed")
		return
	}

	return
}

func (b *Buffer) WriteU32(v uint32) (err error) {

	if b.pos+4 > len(b.data) {
		err = errors.New("write U32 failed")
		return
	}

	for i := 0; i < 4; i++ {
		b.data[b.pos] = byte(v >> (3 - uint(i)) * 8)
		b.pos++
	}

	return
}

func (b *Buffer) WriteS32(v int32) (err error) {
	err = b.WriteU32(uint32(v))
	if err != nil {
		err = errors.New("write S32 failed")
		return
	}

	return
}

func (b *Buffer) WriteU64(v uint64) (err error) {

	if b.pos+8 > len(b.data) {
		err = errors.New("write U64 failed")
		return
	}

	for i := 0; i < 8; i++ {
		b.data[b.pos] = byte(v >> (7 - uint(i)) * 8)
		b.pos++
	}

	return
}

func (b *Buffer) WriteS64(v int64) (err error) {
	err = b.WriteU64(uint64(v))
	if err != nil {
		err = errors.New("write S64 failed")
		return
	}

	return
}

func (b *Buffer) WriteFloat32(f float32) (err error) {
	v := math.Float32bits(f)
	err = b.WriteU32(v)
	if err != nil {
		err = errors.New("write float32 failed")
		return
	}

	return
}

func (b *Buffer) WriteFloat64(f float64) (err error) {
	v := math.Float64bits(f)
	err = b.WriteU64(v)
	if err != nil {
		err = errors.New("write float64 failed")
		return
	}

	return
}

//reader
func (b *Buffer) ReadBool() (ret bool, err error) {
	r, _err := b.ReadByte()

	if r != byte(1) {
		return false, _err
	}

	return true, _err
}

func (b *Buffer) ReadByte() (ret byte, err error) {
	if b.pos >= len(b.data) {
		err = errors.New("read byte failed")
		return
	}

	ret = b.data[b.pos]
	b.pos++
	return

}

func (b *Buffer) ReadBytes() (ret []byte, err error) {
	if b.pos+2 > len(b.data) {
		err = errors.New("read bytes header failed")
		return
	}

	size, _ := b.ReadU16()
	if b.pos+int(size) > len(b.data) {
		err = errors.New("read bytes data failed")
		return
	}

	ret = b.data[b.pos : b.pos+int(size)]
	b.pos += int(size)

	return
}

func (b *Buffer) ReadString() (ret string, err error) {
	if b.pos+2 > len(b.data) {
		err = errors.New("read string header failed")
		return
	}

	size, _ := b.ReadU16()
	if b.pos+int(size) > len(b.data) {
		err = errors.New("read string data failed")
		return
	}

	bytes := b.data[b.pos : b.pos+int(size)]
	b.pos += int(size)
	ret = string(bytes)

	return
}

func (b *Buffer) ReadS8() (ret int8, err error) {
	_ret, _err := b.ReadByte()
	ret = int8(_ret)
	err = _err
	return
}

func (b *Buffer) ReadU16() (ret uint16, err error) {
	if b.pos+2 > len(b.data) {
		err = errors.New("read uint16 failed")
		return
	}

	buf := b.data[b.pos : b.pos+2]

	//log.Debugf("Buffer ReadU16 buf:%v", buf)
	ret = uint16(buf[0])<<8 | uint16(buf[1])
	//log.Debugf("Buffer ReadU16 ret:%v", ret)
	b.pos += 2

	return
}

func (b *Buffer) ReadS16() (ret int16, err error) {
	_ret, _err := b.ReadU16()
	ret = int16(_ret)
	err = _err
	return
}

func (b *Buffer) ReadU32() (ret uint32, err error) {
	if b.pos+4 > len(b.data) {
		err = errors.New("read uint32 failed")
		return
	}

	buf := b.data[b.pos : b.pos+4]
	ret = uint32(buf[0]<<24) | uint32(buf[1]<<16) | uint32(buf[2]<<8) | uint32(buf[3])
	b.pos += 4

	return
}

func (b *Buffer) ReadS32() (ret int32, err error) {
	_ret, _err := b.ReadU32()
	ret = int32(_ret)
	err = _err

	return
}

func (b *Buffer) ReadU64() (ret uint64, err error) {
	if b.pos+8 > len(b.data) {
		err = errors.New("read uint64 failed")
		return
	}

	ret = 0
	buf := b.data[b.pos : b.pos+8]
	for i, v := range buf {
		ret |= uint64(v) << uint((7-i)*8)
	}

	b.pos += 8
	return
}

func (b *Buffer) ReadS64() (ret int64, err error) {
	_ret, _err := b.ReadU64()
	ret = int64(_ret)
	err = _err
	return
}

func (b *Buffer) ReadFloat32() (ret float32, err error) {
	_ret, _err := b.ReadU32()
	if _err != nil {
		return float32(0), _err
	}

	ret = math.Float32frombits(_ret)
	if math.IsNaN(float64(ret)) || math.IsInf(float64(ret), 0) {
		return 0, nil
	}

	return ret, nil
}

func (b *Buffer) ReadFloat64() (ret float64, err error) {
	_ret, _err := b.ReadU64()
	if _err != nil {
		return float64(0), _err
	}

	ret = math.Float64frombits(_ret)
	if math.IsNaN(ret) || math.IsInf(ret, 0) {
		return 0, nil
	}

	return ret, nil
}
