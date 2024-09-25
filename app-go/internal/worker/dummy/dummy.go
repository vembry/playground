package dummy

type dummy struct {
}

func New() *dummy {
	return &dummy{}
}

func (d *dummy) Name() string {
	return "dummy"
}
func (d *dummy) Start() {
}
func (d *dummy) Stop() {
}
