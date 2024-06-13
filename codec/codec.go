package codec

import "io"

// Header
// 一个典型的 RPC 调用如下：err = client.Call("Arith.Multiply", args, &reply)
// 客户端发送的请求包括服务名 Arith，方法名 Multiply，参数 args 三个，服务端的响应包括错误 error，返回值 reply 2 个。
// 将请求和响应中的参数和返回值抽象为 body，剩余的信息放在 header 中，那么就可以抽象出数据结构 Header
type Header struct {
	ServiceMethod string // "Service.Method" 服务名和方法名，通常与 Go 语言中的结构体和方法相映射
	Seq           uint64 // 请求的序号，也可以认为是某个请求的 ID，用来区分不同的请求
	Error         string // 错误信息，客户端置为空，服务端如果如果发生错误，将错误信息置于 Error 中
}

// Codec 对消息体进行编解码的接口
type Codec interface {
	io.Closer
	ReadHeader(*Header) error // 将conn中的数据解码到 Header 结构体中
	ReadBody(any) error       // 将conn中的数据解码到 body 中
	Write(*Header, any) error // 将 Header 和 body 编码为二进制格式，并写入conn
}

// NewCodecFunc 抽象出 Codec 的构造函数，客户端和服务端可以通过 Codec 的 Type 得到构造函数，从而创建 Codec 实例
type NewCodecFunc func(io.ReadWriteCloser) Codec

type Type string

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json" // not implemented
)

var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}
