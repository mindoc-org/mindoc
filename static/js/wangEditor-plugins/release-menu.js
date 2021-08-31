(function () {

    // 获取 wangEditor 构造函数和 jquery
    var E = window.wangEditor;
    var $ = window.jQuery;

    let ReleaseMenu = (function (_BtnMenu) {
      c2b_inherits(ReleaseMenu, _BtnMenu);

      var _super = c2b_createSuper(ReleaseMenu);

      function ReleaseMenu(editor) {
        c2b_classCallCheck(this, ReleaseMenu);

        const $elem = E.$(`<div class="w-e-menu" data-title="发布">
                    <i class="fa fa-cloud-upload" aria-hidden="true"></i>
                </div>`);
        return _super.call(this, $elem, editor);
      }

      c2b_createClass(ReleaseMenu, [
        {
          key: "clickHandler",
          value: function clickHandler() {
            window.releaseBook();
          }
        },
        {
          key: "tryChangeActive",
          value: function tryChangeActive() {
            // this.active();
          }
        }
      ]);

      return ReleaseMenu;
    })(E.BtnMenu);


    var menuKey = 'release';
    E.registerMenu(menuKey, ReleaseMenu);
})();