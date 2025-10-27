package shell

type InitHook int

const (
	None InitHook = iota
)

func (h InitHook) String() string {
	switch h {
	case None:
		return "None"
	default:
		return "Unknown"
	}
}