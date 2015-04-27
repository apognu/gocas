package authenticator

type Dummy struct{}

func (Dummy) Auth(u string, p string) bool {
	if u == p {
		return true
	}
	return false
}
