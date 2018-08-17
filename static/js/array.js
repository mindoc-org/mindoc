/**
 * 删除数组中的匹配值
 * @param $callback
 */
Array.prototype.remove = function ($callback) {
    var $isFunction = typeof $callback === "function";

    var arr = [];
    for(var $i = 0,$len = this.length; $i < $len;$i ++){
        if($isFunction){
            if($callback(this[$i])){
                arr.push($i);
            }
        }else if(this[$i] == $callback){
            arr.push($i);
        }
    }
    for($i = 0,$len = arr.length; $i < $len;$i++){
        this.slice($i,1);
    }
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
