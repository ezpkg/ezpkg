package examples

func AppendStr(s, part string) string {
	return s + part
}

func WillPanic() {
	panic("let's panic! 💥")
}

func CallFunc(fn func()) {
	fn()
}
