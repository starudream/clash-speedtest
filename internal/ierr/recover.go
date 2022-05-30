package ierr

func Recover(fs ...func()) {
	if ev := recover(); ev != nil {
		CheckErr(ev)
	}
	for i := 0; i < len(fs); i++ {
		fs[i]()
	}
}
