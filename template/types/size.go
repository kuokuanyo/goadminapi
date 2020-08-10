package types

import "strconv"

type S map[string]string

func SizeMD(md int) S {
	var s = make(S)
	if md > 0 && md < 13 {
		s["md"] = strconv.Itoa(md)
	}
	return s
}
