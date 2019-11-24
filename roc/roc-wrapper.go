package roc

import "errors"

var (
	// ErrInvalidArguments indicates that one or more arguments passed to the function
	// are invalid
	ErrInvalidArguments = errors.New("One or more arguments are invalid")

	// ErrInvalidApi should never happen and indicates that the API don't follow the declared contract
	ErrInvalidApi = errors.New("Invalid API")
)

// NewAddress parses the `ip`, `port` and `family` and initializes the Address object
func NewAddress(family Family, ip string, port int) (*Address, error) {
	a := new(Address)
	errCode := addressInit(a, family, ip, int32(port))

	if errCode == 0 {
		return a, nil
	}
	if errCode < 0 {
		return nil, ErrInvalidArguments
	}
	return nil, ErrInvalidApi
}

// Family get Address family
func (a *Address) Family() Family {
	return addressFamily(a)
}

// Ip get Address ip
func (a *Address) IP() string {
	const buflen = 50
	buf := [buflen]byte{}
	return addressIp(a, buf[:], buflen)
}

// Port get Address port
func (a *Address) Port() int {
	return int(addressPort(a))
}

func OpenContext(config *ContextConfig) *Context {
	return contextOpen(config)
}

func (c *Context) Close() error {
	errCode := contextClose(c)
	if errCode == 0 {
		return nil
	}
	if errCode < 0 {
		return errors.New("Arguments are invalid or or there are objects attached to the context")
	}
	return ErrInvalidApi
}

func OpenReceiver(ctx *Context, config *ReceiverConfig) *Receiver {
	return receiverOpen(ctx, config)
}

func (r *Receiver) Bind(portType PortType, proto Protocol, address *Address) error {
	errCode := receiverBind(r, portType, proto, address)
	if errCode == 0 {
		return nil
	}
	if errCode < 0 {
		return errors.New("Arguments are invalid or Address can't be bound or There aren't enough resources")
	}
	return ErrInvalidApi
}

func (r *Receiver) Close() error {
	errCode := receiverClose(r)
	if errCode == 0 {
		return nil
	}
	if errCode < 0 {
		return ErrInvalidArguments
	}
	return ErrInvalidApi
}

func OpenSender(ctx *Context, config *SenderConfig) *Sender {
	return senderOpen(ctx, config)
}

func (s *Sender) Bind(address *Address) error {
	errCode := senderBind(s, address)
	if errCode == 0 {
		return nil
	}
	if errCode < 0 {
		return errors.New("Arguments are invalid or Sender is already bound or Address can't be bound or There aren't enough resources")
	}
	return ErrInvalidApi
}

func (s *Sender) Connect(address *Address, portType PortType, proto Protocol) error {
	errCode := senderConnect(s, portType, proto, address)
	if errCode == 0 {
		return nil
	}
	return errors.New("Arguments are invalid or roc_sender_write() was already called")
}

func (s *Sender) Close() error {
	errCode := senderClose(s)
	if errCode == 0 {
		return nil
	}
	if errCode < 0 {
		return ErrInvalidArguments
	}
	return ErrInvalidApi
}
