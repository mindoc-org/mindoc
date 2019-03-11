/**
 * 删除数组中的匹配值
 * @param $callback
 */
Array.prototype.remove = function ($callback) {
    var $isFunction = typeof $callback === "function";

    var arr = [];
    for (var $i = 0, $len = this.length; $i < $len; $i++) {
        if ($isFunction) {
            if ($callback(this[$i])) {
                arr.push($i);
            }
        } else if (this[$i] == $callback) {
            arr.push($i);
        }
    }
    for ($i = 0, $len = arr.length; $i < $len; $i++) {
        this.slice($i, 1);
    }
};
String.prototype.endWith = function (endStr) {
    var d = this.length - endStr.length;

    return (d >= 0 && this.lastIndexOf(endStr) === d)
};

//格式化文件大小
function formatBytes($size) {
    if (typeof $size === "number") {
        var $units = [" B", " KB", " MB", " GB", " TB"];

        for ($i = 0; $size >= 1024 && $i < 4; $i++) $size /= 1024;

        return $size.toFixed(2) + $units[$i];
    }
    return $size;
}

/**
 * 将多维的json转换为一维的json
 * @param $json
 * @param $parentKey
 */
function foreachJson($json, $parentKey) {
    var data = {};

    $.each($json, function (key, item) {
        var cKey = $parentKey;

        if (Array.isArray($json)) {
            key = "[";
        }

        if ($parentKey !== undefined && $parentKey !== "" && key !== "") {
            if($parentKey.endsWith("[")) {
                cKey = $parentKey + key + "]";
            } else if (key === "[") {
                cKey = $parentKey + key;
            } else {
                cKey = $parentKey + "." + key;
            }
        } else {
            cKey = key;
        }


        var node = {};
        node["key"] = cKey;
        node["type"] = Array.isArray(item) ? "array" : typeof item;
        node["value"] = item;
        if (typeof key === "string" && key !== "[") {
            data[cKey] = node;
        }

        if (typeof item === "object") {
            var items = foreachJson(item, cKey);
            $.each(items,function (k,v) {
                data[k] = v;
            });
        }
    });
    return data;

}