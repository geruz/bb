package resource

type InProvider interface {
	In(interface{}) error
}
