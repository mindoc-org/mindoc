**基本信息**

|ID|函数|名称|版本|实现|关键字|
|---|---|---|---|---|---|----|
|Z0001|public.TestOnly|echo测试|1.0.0|2.0.1|echo|

**详细描述**

echo测试

**请求参数**

这里仅仅是业务请求参数，不包含平台参数
这里演示了一些参数使用

|字段|必需 |类型|名称 |描述 |示例 |
|:---|:---:|:---:|:---|:---|:---|
|str2|true|string(10)|带长度字符串 ||A001|
|int3|false|int(min,100)|区间限制整数 ||15|
|double1|false|double|任意双精度浮点数| |1|
|dt1|false|datetime|日期时间 ||2017-02-03 12:13:14|
|b1|false|boolean|布尔值| |true|
|enum1|false|enum(BJ,SH,SZ,CQ,TJ)|枚举对象| |SH|
|obj1|false|object|对象| ||
|obj1.str1|false|string|对象的字符串属性| |xy|
|obj1.int1|false|int(1,10)|对象的整数属性| |1|
|lst1|false|list|列表对象|| |
|lst1.str1|false|string|列表对象的字符串属性| |xy|
|lst1.int1|false|int(1,10)|列表对象的整数属性| |1|
|strlst1|false|list&lt;string&gt;|字符串列表| ||
|map1|false|map&lt;string,string&gt;|字符串MAP| ||

**响应结果**

这里仅仅是业务响应结果，不包含平台响应结果， 平台调用失败时候只有平台的响应结果
这是echo测试，返回和请求相同
	
|字段|必需 |类型|名称 |描述 |示例 |
|:---|:---:|:---:|:---|:---|:---|
|str2|true|string(10)|带长度字符串 ||A001|
|int3|false|int(min,100)|区间限制整数 ||15|
|double1|false|double|任意双精度浮点数| |1|
|dt1|false|datetime|日期时间 ||2017-02-03 12:13:14|
|b1|false|boolean|布尔值| |true|
|enum1|false|enum(BJ,SH,SZ,CQ,TJ)|枚举对象| |SH|
|obj1|false|object|对象| ||
|obj1.str1|false|string|对象的字符串属性| |xy|
|obj1.int1|false|int(1,10)|对象的整数属性| |1|
|lst1|false|list|列表对象|| |
|lst1.str1|false|string|列表对象的字符串属性| |xy|
|lst1.int1|false|int(1,10)|列表对象的整数属性| |1|
|strlst1|false|list&lt;string&gt;|字符串列表| ||
|map1|false|map&lt;string,string&gt;|字符串MAP| ||

**响应错误定义**

这里仅仅是业务处理的错误，错误定义为4未整数，需要大于1000
	
|错误代码   | 错误描述 |
|:----------|:----:|
|1001 | 枚举对象没有传SH |
