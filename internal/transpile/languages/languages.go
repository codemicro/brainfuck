package languages

type Language interface {
	Header()
	MainStart()
	MainEnd()
	ValueDelta(delta int)
	PointerDelta(delta int)
	Output()
	Input()
	LoopStart()
	LoopEnd()
	Bytes() []byte
}

var Index = map[string]Language{
	// TODO: We really don't need to initialise all of these languages here
	"go": NewGo(),
	"python": NewPython(),
}
