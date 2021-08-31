(function () {

    // 获取 wangEditor 构造函数和 jquery
    var E = window.wangEditor;
    var $ = window.jQuery;

    let HistoryMenu = (function (_BtnMenu) {
      c2b_inherits(HistoryMenu, _BtnMenu);

      var _super = c2b_createSuper(HistoryMenu);

      function HistoryMenu(editor) {
        c2b_classCallCheck(this, HistoryMenu);

        const $elem = E.$(`<div class="w-e-menu" data-title="history">
                    <i class="fa fa-history" aria-hidden="true"></i>
                </div>`);
        return _super.call(this, $elem, editor);
      }

      c2b_createClass(HistoryMenu, [
        {
          key: "clickHandler",
          value: function clickHandler() {
            window.documentHistory();
          }
        },
        {
          key: "tryChangeActive",
          value: function tryChangeActive() {
            // this.active();
          }
        }
      ]);

      return HistoryMenu;
    })(E.BtnMenu);


    var menuKey = '历史';
    E.registerMenu(menuKey, HistoryMenu);
})();