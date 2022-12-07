package adnl

// ADNL中基于UDP的P2P协议的最简实现
type UDPDatagram struct {
	Receiver [32]byte // receiver的256位抽象地址
	Data     struct {
		Payload  []byte   // 真正需要发送的数据
		Sign     []byte   // sender的签名
		Sender   [32]byte // sender的256位抽象地址
		Preimage []byte   // （当receiver尚不知道sender的抽象地址时，需发送Preimage; 抽象地址是Preimage的256位哈希）
	}
}

// 将UDPDatagram对象封装成一个对应的UDPDatagram数据帧
func (ud *UDPDatagram) Packet() {}

// 临时用
type Destination struct {
	IP   string
	Port string
}

// 临时用
type UDPacket []byte

// 发送UDPDatagram给receiver
func Send(dest Destination, packet UDPacket) {}

// 接收到某个sender发来的UDPDatagram
func Receive(datagram UDPDatagram) {
	// 1.取出前256位抽象地址
	// 2.根据该地址取出自身维护的对应的公私钥对
	// 3.利用对应私钥解析出数据对象
	// 4.校验该数据对象（主要是验证签名）
	// 5.后续处理（如更新preimage；以及其它数据更新）
}
