package boc

// 用以标识一切"通信"/"存储"的TL序列化对象的类型
type Magic uint32

// 用以标识同一个结构体类型中，不同的结构体实现
// Notice: 因为这些不同结构体的方法都是同一的；因此不将这些本质上就相同的结构体分开设计
type SumType string
