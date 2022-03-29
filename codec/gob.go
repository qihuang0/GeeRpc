package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

/*
	定义 GobCodec 结构体，这个结构体由四部分构成，conn 是由构建函数传入，
	通常是通过 TCP 或者 Unix 建立 socket 时得到的链接实例，
	dec 和 enc 对应 gob 的 Decoder 和 Encoder，
	buf 是为了防止阻塞而创建的带缓冲的 Writer，一般这么做能提升性能。
*/

type GobCodec struct {
	conn   io.ReadWriteCloser
	buf    *bufio.Writer
	encode *gob.Encoder
	decode *gob.Decoder
}

var _ CodeC = (*GobCodec)(nil)

func NewGobCodec(conn io.ReadWriteCloser) CodeC {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn:   conn,
		buf:    buf,
		encode: gob.NewEncoder(conn),
		decode: gob.NewDecoder(buf),
	}
}

func (g *GobCodec) Close() error {
	return g.conn.Close()
}

func (g *GobCodec) ReadHeader(header *Header) error {
	return g.decode.Decode(header)
}

func (g *GobCodec) ReadBody(i interface{}) error {
	return g.encode.Encode(i)
}

func (g *GobCodec) Write(header *Header, i interface{}) (err error) {
	defer func() {
		_ = g.buf.Flush()
		if err != nil {
			_ = g.Close()
		}
	}()
	if err := g.encode.Encode(header); err != nil {
		log.Println("rpc codec: gob error encoding header:", err)
		return err
	}
	if err := g.encode.Encode(i); err != nil {
		log.Println("rpc codec: gob error encoding body:", err)
		return err
	}
	return nil
}
