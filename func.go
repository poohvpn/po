package pooh

func TryOk(fs ...func() bool) bool {
	for _, f := range fs {
		if f() {
			return true
		}
	}
	return false
}
