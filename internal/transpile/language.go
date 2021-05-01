package transpile

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
