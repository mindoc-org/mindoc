(function () {

    // 获取 wangEditor 构造函数和 jquery
    var E = window.wangEditor;
    var $ = window.jQuery;

    let SaveMenu = (function (_BtnMenu) {
      c2b_inherits(SaveMenu, _BtnMenu);

      var _super = c2b_createSuper(SaveMenu);

      function SaveMenu(editor) {
        c2b_classCallCheck(this, SaveMenu);

        const $elem = E.$(`<div class="w-e-menu" data-title="保存">
                    <i class="fa fa-floppy-o" aria-hidden="true"></i>
                </div>`);
        return _super.call(this, $elem, editor);
      }

      c2b_createClass(SaveMenu, [
        {
          key: "clickHandler",
          value: function clickHandler() {
            window.saveDocument();
          }
        },
        {
          key: "tryChangeActive",
          value: function tryChangeActive() {
            // this.active();
          }
        }
      ]);

      return SaveMenu;
    })(E.BtnMenu);


    var menuKey = 'save';
    E.registerMenu(menuKey, SaveMenu);
})();