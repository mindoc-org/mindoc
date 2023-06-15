/**
 * 自定义一个语法
 */
var CustomHookA = Cherry.createSyntaxHook('codeBlock', Cherry.constants.HOOKS_TYPE_LIST.PAR, {
  makeHtml(str) {
    console.warn('custom hook', 'hello');
    return str;
  },
  rule(str) {
    const regex = {
      begin: '',
      content: '',
      end: '',
    };
    regex.reg = new RegExp(regex.begin + regex.content + regex.end, 'g');
    return regex;
  },
});
/**
 * 自定义一个自定义菜单
 * 点第一次时，把选中的文字变成同时加粗和斜体
 * 保持光标选区不变，点第二次时，把加粗斜体的文字变成普通文本
 */
var customMenuA = Cherry.createMenuHook('加粗斜体', {
  iconName: 'font',
  onClick: function (selection) {
    // 获取用户选中的文字，调用getSelection方法后，如果用户没有选中任何文字，会尝试获取光标所在位置的单词或句子
    let $selection = this.getSelection(selection) || '同时加粗斜体';
    // 如果是单选，并且选中内容的开始结束内没有加粗语法，则扩大选中范围
    if (!this.isSelections && !/^\s*(\*\*\*)[\s\S]+(\1)/.test($selection)) {
      this.getMoreSelection('***', '***', () => {
        const newSelection = this.editor.editor.getSelection();
        const isBoldItalic = /^\s*(\*\*\*)[\s\S]+(\1)/.test(newSelection);
        if (isBoldItalic) {
          $selection = newSelection;
        }
        return isBoldItalic;
      });
    }
    // 如果选中的文本中已经有加粗语法了，则去掉加粗语法
    if (/^\s*(\*\*\*)[\s\S]+(\1)/.test($selection)) {
      return $selection.replace(/(^)(\s*)(\*\*\*)([^\n]+)(\3)(\s*)($)/gm, '$1$4$7');
    }
    /**
     * 注册缩小选区的规则
     *    注册后，插入“***TEXT***”，选中状态会变成“***【TEXT】***”
     *    如果不注册，插入后效果为：“【***TEXT***】”
     */
    this.registerAfterClickCb(() => {
      this.setLessSelection('***', '***');
    });
    return $selection.replace(/(^)([^\n]+)($)/gm, '$1***$2***$3');
  }
});
/**
 * 定义一个空壳，用于自行规划cherry已有工具栏的层级结构
 */
var customMenuB = Cherry.createMenuHook('实验室', {
  iconName: '',
});
/**
 * 定义一个自带二级菜单的工具栏
 */
var customMenuC = Cherry.createMenuHook('帮助中心', {
  iconName: 'question',
  onClick: (selection, type) => {
    switch (type) {
      case 'shortKey':
        return `${selection}快捷键看这里：https://codemirror.net/5/demo/sublime.html`;
      case 'github':
        return `${selection}我们在这里：https://github.com/Tencent/cherry-markdown`;
      case 'release':
        return `${selection}我们在这里：https://github.com/Tencent/cherry-markdown/releases`;
      default:
        return selection;
    }
  },
  subMenuConfig: [
    { noIcon: true, name: '快捷键', onclick: (event) => { cherry.toolbar.menus.hooks.customMenuCName.fire(null, 'shortKey') } },
    { noIcon: true, name: '联系我们', onclick: (event) => { cherry.toolbar.menus.hooks.customMenuCName.fire(null, 'github') } },
    { noIcon: true, name: '更新日志', onclick: (event) => { cherry.toolbar.menus.hooks.customMenuCName.fire(null, 'release') } },
  ]
});

var basicConfig = {
  id: 'markdown',
  externals: {
    echarts: window.echarts,
    katex: window.katex,
    MathJax: window.MathJax,
  },
  isPreviewOnly: false,
  engine: {
    global: {
      urlProcessor(url, srcType) {
        console.log(`url-processor`, url, srcType);
        return url;
      },
    },
    syntax: {
      codeBlock: {
        theme: 'twilight',
      },
      table: {
        enableChart: false,
        // chartEngine: Engine Class
      },
      fontEmphasis: {
        allowWhitespace: false, // 是否允许首尾空格
      },
      strikethrough: {
        needWhitespace: false, // 是否必须有前后空格
      },
      mathBlock: {
        engine: 'MathJax', // katex或MathJax
        src: 'https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-svg.js', // 如果使用MathJax plugins，则需要使用该url通过script标签引入
      },
      inlineMath: {
        engine: 'MathJax', // katex或MathJax
      },
      emoji: {
        useUnicode: false,
        customResourceURL: 'https://github.githubassets.com/images/icons/emoji/unicode/${code}.png?v8',
        upperCase: true,
      },
      // toc: {
      //     tocStyle: 'nested'
      // }
      // 'header': {
      //   strict: false
      // }
    },
    customSyntax: {
      // SyntaxHookClass
      CustomHook: {
        syntaxClass: CustomHookA,
        force: false,
        after: 'br',
      },
    },
  },
  toolbars: {
    toolbar: [
      'bold',
      'italic',
      {
        strikethrough: ['strikethrough', 'underline', 'sub', 'sup', 'ruby', 'customMenuAName'],
      },
      'size',
      '|',
      'color',
      'header',
      '|',
      'drawIo',
      '|',
      'ol',
      'ul',
      'checklist',
      'panel',
      'detail',
      '|',
      'formula',
      {
        insert: ['image', 'audio', 'video', 'link', 'hr', 'br', 'code', 'formula', 'toc', 'table', 'pdf', 'word', 'ruby'],
      },
      'graph',
      'togglePreview',
      'settings',
      'switchModel',
      'codeTheme',
      'export',
      {
        customMenuBName: ['ruby', 'audio', 'video', 'customMenuAName'],
      },
      'customMenuCName',
      'theme'
    ],
    bubble: ['bold', 'italic', 'underline', 'strikethrough', 'sub', 'sup', 'quote', 'ruby', '|', 'size', 'color'], // array or false
    sidebar: ['mobilePreview', 'copy', 'theme'],
    customMenu: {
      customMenuAName: customMenuA,
      customMenuBName: customMenuB,
      customMenuCName: customMenuC,
    },
  },
  drawioIframeUrl: '/static/cherry/drawio_demo.html',
  editor: {
    defaultModel: 'edit&preview',
    height: 800,
  },
  previewer: {
    // 自定义markdown预览区域class
    // className: 'markdown'
  },
  keydown: [],
  //extensions: [],
  //callback: {
  //changeString2Pinyin: pinyin,
  //}
};

fetch('').then((response) => response.text()).then((value) => {
  var markdownarea = document.getElementById("markdown_area").value
  var config = Object.assign({}, basicConfig, { value: markdownarea });// { value: value });不显示获取的初始化值
  window.cherry = new Cherry(config);
});
