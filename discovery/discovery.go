package discovery

type Address struct {
	Host string
	Port int
}
type AdddressProvider func() (Address, error)
