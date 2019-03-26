package main

type appe struct {
	body string
}

func (a *appe) process(src []byte) ([]byte, error) {
	return append(append(src[:0:0], src...), a.body...), nil
}

type inse struct {
	body string
}

func (i *inse) process(src []byte) ([]byte, error) {
	return append(append(src[:0:0], i.body...), src...), nil
}

type del struct{}

func (*del) process(src []byte) ([]byte, error) {
	return []byte(""), nil
}
