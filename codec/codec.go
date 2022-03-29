package codec

import "io"

type Header struct {
	ServiceMethod string // 服务名和方法名
	Seq           uint64 // 请求的序号
	Error         string
}

type CodeC interface { // 编解码的接口Codec
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

type NewCodecFunc func(io.ReadWriteCloser) CodeC // 抽象出构造函数

type Type string

const GobType Type = "application/gob" // 通过Type得到构造函数

var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}
