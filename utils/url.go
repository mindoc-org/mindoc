package utils

import "strings"

func JoinURI(elem ...string) string {
	if len(elem) <= 0 {
		return ""
	}
	uri := ""

	for i,u := range elem {
		u = strings.Replace(u,"\\","/",-1)

		if i == 0 {
			if !strings.HasSuffix(u,"/") {
				u = u + "/"
			}
			uri = u
		}else{
			u = strings.Replace(u,"//","/",-1)
			if strings.HasPrefix(u,"/") {
				u = string(u[1:])
			}
			uri += u
		}
	}
	return uri
}
