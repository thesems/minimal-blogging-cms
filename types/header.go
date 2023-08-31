package types

type Header struct {
	Navigation []*Page
	User       string
}

func NewHeader(navigation []*Page, user string) *Header {
	return &Header{Navigation: navigation, User: user}
}
