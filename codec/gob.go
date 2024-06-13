package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

// GobCodec 实现了 Codec 接口
type GobCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *gob.Decoder
	enc  *gob.Encoder
}

// ReadHeader 将conn中的数据解码到 Header 结构体中
func (c *GobCodec) ReadHeader(h *Header) error {
	return c.dec.Decode(h)
}

// ReadBody 将conn中的数据解码到 body 中
func (c *GobCodec) ReadBody(body any) error {
	return c.dec.Decode(body)
}

// Write 将 Header 和 body 编码为二进制格式，并写入conn
func (c *GobCodec) Write(h *Header, body any) (err error) {
	defer func() {
		c.buf.Flush()
		if err != nil {
			c.Close()
		}
	}()
	if err := c.enc.Encode(h); err != nil {
		log.Println("rpc codec: gob error encoding header:", err)
		return err
	}
	if err := c.enc.Encode(body); err != nil {
		log.Println("rpc codec: gob error encoding body:", err)
		return err
	}
	return nil
}

func (c *GobCodec) Close() error {
	return c.conn.Close()
}

// NewGobCodec NewCodecFunc函数类型的变量，创建一个 GobCodec 实例
func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}
