(function () {

    // 获取 wangEditor 构造函数和 jquery
    var E = window.wangEditor;
    var $ = window.jQuery;

    let AttachMenu = (function (_BtnMenu) {
      c2b_inherits(AttachMenu, _BtnMenu);

      var _super = c2b_createSuper(AttachMenu);

      function AttachMenu(editor) {
        c2b_classCallCheck(this, AttachMenu);

        const $elem = E.$(`<div class="w-e-menu" data-title="附件">
                    <i class="fa fa-paperclip" aria-hidden="true"></i>
                </div>`);
        return _super.call(this, $elem, editor);
      }

      c2b_createClass(AttachMenu, [
        {
          key: "clickHandler",
          value: function clickHandler() {
            $("#uploadAttachModal").modal("show");
          }
        },
        {
          key: "tryChangeActive",
          value: function tryChangeActive() {
            // this.active();
          }
        }
      ]);

      return AttachMenu;
    })(E.BtnMenu);


    var menuKey = 'attach';
    E.registerMenu(menuKey, AttachMenu);
})();