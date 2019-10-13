package roc

import "fmt"

func convertErr(err int32, errstr string) error {
	if err == 0 {
		return nil
	}
	return fmt.Errorf("%s: %d", errstr, err)
}

// Init parses the `ip` and `port` and initializes the Address object
func NewAddress(family Family, ip string, port int) (*Address, error) {
	a := new(Address)
	err := addressInit(a, family, ip, int32(port))
	return a, convertErr(err, "Error when initializing an address")
}

// Family get Address family
func (a *Address) Family() Family {
	return addressFamily(a)
}

// Ip get Address ip
func (a *Address) IP() string {
	const buflen = 50
	buf := [buflen]byte{}
	return addressIp(a, buf[:])
}

// Port get Address port
func (a *Address) Port() int {
	return int(addressPort(a))
}

func OpenContext(config *ContextConfig) *Context {
	return contextOpen(config)
}

func (c *Context) Close() error {
	return convertErr(contextClose(c), "Error when closing context")
}

func OpenReceiver(ctx *Context, config *ReceiverConfig) *Receiver {
	return receiverOpen(ctx, config)
}

func (r *Receiver) Bind(portType PortType, proto Protocol, address *Address) error {
	err := receiverBind(r, portType, proto, address)
	return convertErr(err, "Error while binding a receiver")
}

func (r *Receiver) Read(frame *Frame) error {
	err := receiverRead(r, frame)
	return convertErr(err, "Error while reading from receiver")
}

func (r *Receiver) Close() error {
	err := receiverClose(r)
	return convertErr(err, "Error while closing a receiver")
}

func OpenSender(ctx *Context, config *SenderConfig) *Sender {
	return senderOpen(ctx, config)
}

func (s *Sender) Bind(address *Address) error {
	err := senderBind(s, address)
	return convertErr(err, "Error while binding address")
}

func (s *Sender) Connect(address *Address, portType PortType, proto Protocol) error {
	err := senderConnect(s, portType, proto, address)
	return convertErr(err, "Error while conecting sender")
}

func (s *Sender) Write(frame *Frame) error {
	err := senderWrite(s, frame)
	return convertErr(err, "Error while writing with sender")
}

func (s *Sender) Close() error {
	err := senderClose(s)
	return convertErr(err, "Error while closing sender")
}
