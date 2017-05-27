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

func Int32Addr(i int32) *int32 {
	return &i
}

func BoolAddr(b bool) *bool {
	return &b
}

func IsElementInArray(e interface{}, a []interface{}) bool {
	for _, element := range a {
		if element == e {
			return true
		}
	}
	return false
}
