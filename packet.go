/*
MIT License

Copyright (c) 2017 Jumbo
Author: Jumbo
Email: Jumbo.Wu@hotmail.com

File: packet.go

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
  Process buffer data
*/


package packet

import (
	"errors"

	//log "github.com/sirupsen/logrus"
)

const (
	PACKET_LIMIT = 1024 * 32 //32k
)

type Packet struct {
	datalen uint16 //header 2
	encrypt byte   //header 1
	data    []byte //payload protobuf数据 PkgHead | PkgBody(IMessage) size:PACKET_LIMIT - 3
}

func (p *Packet) DataLen() uint16 {
	return p.datalen
}

func (p *Packet) Encrypt() bool {

	if p.encrypt == byte(1) {
		return true
	} else {
		return false
	}

}

func (p *Packet) Data() []byte {
	return p.data
}

func (p *Packet) SetEncrypt(b bool) {
	if b == true {
		p.encrypt = byte(1)
	} else {
		p.encrypt = byte(0)
	}
}
func Pkt(data []byte) *Packet {
	return &Packet{data: data, datalen: uint16(len(data))}
}

//输入参数:protobuf 序列化后的二进制数据，返回pack封装后的二进制
func Pack(p *Packet, b *Buffer) (err error) {

	//log.Debugf("pack packet:%v, Buffer pos:%v, Buffer:%v", p, b.pos, b)
	if p.datalen < 0 {
		err = errors.New("pack minus refer")
		return
	}

	if p.datalen > PACKET_LIMIT {
		err = errors.New("pack surpass count")
		return
	}

	if b.Left() < 2+1+int(p.datalen) {
		err = errors.New("pack short buf for write")
		return
	}

	err = b.WriteU16(p.datalen)
	if err != nil {
		err = errors.New("write datalen failed")
		return
	}
	//log.Debugf("after write datalen pack packet:%v, Buffer pos:%v, Buffer:%v", p, b.pos, b)

	err = b.WriteByte(p.encrypt)
	if err != nil {
		err = errors.New("write encrypt failed")
		return
	}

	//log.Debugf("after write encrypt pack packet:%v, Buffer pos:%v, Buffer:%v", p, b.pos, b)
	for i := 0; i < int(p.datalen); i++ {
		err = b.WriteByte(p.data[i])
		//log.Debugf("pack packet data:%v i:%", b, i)
		if err != nil {
			return
		}
	}

	//log.Debugf("after pack packet:%v", b)
	return
}

func Unpack(b *Buffer) (p *Packet, err error) {

	p = Pkt(make([]byte, 0))
	//log.Debugf("unpack read buf :%v", b)
	p.datalen, err = b.ReadU16()
	//log.Debugf("unpack read U16 for datalen:%v", p.datalen)
	if err != nil {
		err = errors.New("unpack read dataLen failed")
		//log.Debugf("unpack read dataLen failed:%v", p.DataLen())
		return
	}

	if b.Left() < 1+int(p.datalen) {
		err = errors.New("unpack short buf for read")
		return
	}

	p = Pkt(make([]byte, p.datalen))

	p.encrypt, err = b.ReadByte()
	if err != nil {
		return
	}

	if p.datalen < 0 {
		err = errors.New("unpack minus refer")
		return
	}

	if p.datalen > PACKET_LIMIT {
		err = errors.New("unpack surpass count")
		return
	}

	for i := 0; i < int(p.datalen); i++ {
		p.data[i], err = b.ReadByte()
		if err != nil {
			return
		}
	}

	//log.Debugf("p:%v", p.DataLen())

	return
}
