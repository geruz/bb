package discovery

import (
	"errors"
)

func UseSimpleProvider(addr *string, port *int) AdddressProvider {
	return func() (Address, error) {
		if addr == nil {
			return Address{}, errors.New("addres ponter is nil")
		}
		if port == nil {
			return Address{}, errors.New("port ponter is nil")
		}
		return Address{*addr, *port}, nil
	}
}
