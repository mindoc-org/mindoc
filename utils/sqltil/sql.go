package sqltil

import "strings"

//转义like语法的%_符号
func EscapeLike(keyword string) string {
	return strings.Replace(strings.Replace(keyword,"_","\\_",-1),"%","\\%",-1)
}
