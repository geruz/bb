package configuration

type Control struct {
}

func (this Control) Stop() {}
func (this Control) Pause() {}

func NewControl() Control {
	return Control{}
}
