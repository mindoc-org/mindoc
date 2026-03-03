https://github.com/liu11hao11/pinyin_js

# pinyin_js
中文转拼音
##安装
```
npm install
```


##汉字转化成带音节的拼音
```javascript
var pinyin=require("pinyin_js");
console.log(pinyin.pinyin("你好"," "));
//输出结果是nǐ hǎo 
```

##汉字转化成不带音节的拼音
```javascript
var pinyin=require("pinyin_js");
console.log(pinyin.pinyinWithOutYin("你好"," "));
//输出结果是ni hao 
```

##判断是否是汉字
```javascript
var pinyin=require("pinyin_js");
console.log(pinyin.isChineseWord("你好"));//true
console.log(pinyin.isChineseWord("!你好",false));//true
console.log(pinyin.isChineseWord("!你好",true));//第二个参数：true是严格模式，默认为严格模式
//false
```

##首字母排序

```javascript
var pinyin=require("pinyin_js");
var users = [
    { 'user': '张三丰',   'age': 40 },
    { 'user': '123',   'age': 48 },
    { 'user': '张三',   'age': 48 },
    { 'user': '李四', 'age': 36 },  
    { 'user': '张三炮', 'age': 34 }
];
var sortResult = pinyin.sort(users, "user");
console.log(sortResult)
/*[ { user: '123', age: 48 },
    { user: '李四', age: 36 },
    { user: '张三', age: 48 },
    { user: '张三丰', age: 40 },
    { user: '张三炮', age: 34 } ]*/
    
```