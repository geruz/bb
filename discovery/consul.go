package discovery

import (
	"fmt"
	"gitlab.corp.24au/golib/resolver"
)

func UseConsul(name string, major int, minor int) AdddressProvider {
	consulVersion := fmt.Sprintf("v%d-%d", major, minor)

	return func() (Address, error) {
		address, err := resolver.AddressFor(name, consulVersion)
		if err != nil {
			return Address{}, err
		}
		return Address{address.Host, address.Port}, nil
	}
}
