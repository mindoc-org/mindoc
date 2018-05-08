package acl

import (
	"sync"
	"fmt"
)

var AclList = &sync.Map{}

func AddMemberPermission(account string,resource Resource)  {
	key := account + "!" + resource.Code

	AclList.Store(key,resource)
}

//判断指定的资源是否可以访问
func IsAllow(account ,controllerName,actionName,methodName string) bool {

	key := fmt.Sprintf("%s!%s!%s!%s",account,controllerName,actionName,methodName)
	fmt.Println(key)
	if _,ok := AclList.Load(key);ok {
		return true
	}
	key = fmt.Sprintf("%s!%s!%s!*",account,controllerName,actionName)
	fmt.Println(key)
	if _,ok := AclList.Load(key);ok {
		return true
	}

	return false
}
