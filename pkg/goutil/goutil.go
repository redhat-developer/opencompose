package goutil

// Returns string address. Go won't allow you to take address of non-composite type
func StringAddr(s string) *string {
	return &s
}

func StringOrEmpty(p *string) string {
	if p != nil {
		return *p
	}
	return ""
}
