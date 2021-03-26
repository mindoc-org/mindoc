!function (window, document) {
    function d(a) {
        var e, c = document.createElement("iframe"),
            d = "https://login.dingtalk.com/login/qrcode.htm?goto=" + a.goto ;
        d += a.style ? "&style=" + encodeURIComponent(a.style) : "",
            d += a.href ? "&href=" + a.href : "",
            c.src = d,
            c.frameBorder = "0",
            c.allowTransparency = "true",
            c.scrolling = "no",
            c.width =  a.width ? a.width + 'px' : "365px",
            c.height = a.height ? a.height + 'px' : "400px",
            e = document.getElementById(a.id),
            e.innerHTML = "",
            e.appendChild(c)
    }
    window.DDLogin = d
}(window, document);