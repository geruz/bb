package transport

type Extension interface {
	Configure(*HttpTransport)
}
