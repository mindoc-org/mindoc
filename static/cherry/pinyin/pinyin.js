/**
 * Created by Alex on 2016/3/7.
 */
var hzpy = require("./hanziPinyin").hzpy;
var hzpyWithOutYin = require("./hanziPinyinWithoutYin").hzpy;
var _ = require("lodash");

function pinyin(word,splitStr) {
    splitStr = splitStr === undefined ? ' ' : splitStr;
    var str = '';
    var s;
    for (var i = 0; i < word.length; i++) {
        if (hzpy.indexOf(word.charAt(i)) != -1 && word.charCodeAt(i) > 200) {
            s = 1;
            while (hzpy.charAt(hzpy.indexOf(word.charAt(i)) + s) != ",") {
                str += hzpy.charAt(hzpy.indexOf(word.charAt(i)) + s);
                s++;
            }
            str += splitStr;
        }
        else {
            str += word.charAt(i);
        }
    }
    return str;
}

//无声调的拼音
function pinyinWithOutYin(word,splitStr) {
    splitStr = splitStr === undefined ? ' ' : splitStr;
    var str = '';
    var s;
    for (var i = 0; i < word.length; i++) {
        if (hzpyWithOutYin.indexOf(word.charAt(i)) != -1 && word.charCodeAt(i) > 200) {
            s = 1;
            while (hzpyWithOutYin.charAt(hzpyWithOutYin.indexOf(word.charAt(i)) + s) != ",") {
                str += hzpyWithOutYin.charAt(hzpyWithOutYin.indexOf(word.charAt(i)) + s);
                s++;
            }
            str +=splitStr;
        }
        else {
            str += word.charAt(i);
        }
    }

    return str;
}

function isChineseWord(word, modle) {
    if (!modle) {
        //modle为false是非严格中文！默认是严格中文
        modle = true;
    }
    var str = '';
    var isChinese = false;
    for (var i = 0; i < word.length; i++) {
        if (hzpy.indexOf(word.charAt(i)) != -1 && word.charCodeAt(i) > 200) {
            isChinese = true;
        }
        else {
            if (modle) {
                return false;
            }
        }
    }
    return isChinese;
}

function sort(array, key) {
    return _.sortBy(array, [function (o) {
        return pinyinWithOutYin(o[key],"");
    }]);
}

module.exports = {
    pinyin: pinyin,
    pinyinWithOutYin: pinyinWithOutYin,
    isChineseWord: isChineseWord,
    sort: sort
}
