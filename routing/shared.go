package routing

var Shared = NewSessionMux2("")

func HFilter(pattern string, h Handler) {
	Shared.HFilter(pattern, h)
}
func HFilterFunc(pattern string, h HandleFunc) {
	Shared.HFilterFunc(pattern, h)
}
func H(pattern string, h Handler) {
	Shared.H(pattern, h)
}
func HFunc(pattern string, h HandleFunc) {
	Shared.HFunc(pattern, h)
}