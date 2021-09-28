/*! markmap-lib v0.7.8 | MIT License */
(function (exports, d3) {
'use strict';

function count(node) {
  var sum = 0,
      children = node.children,
      i = children && children.length;
  if (!i) sum = 1;
  else while (--i >= 0) sum += children[i].value;
  node.value = sum;
}

function node_count() {
  return this.eachAfter(count);
}

function node_each(callback) {
  var node = this, current, next = [node], children, i, n;
  do {
    current = next.reverse(), next = [];
    while (node = current.pop()) {
      callback(node), children = node.children;
      if (children) for (i = 0, n = children.length; i < n; ++i) {
        next.push(children[i]);
      }
    }
  } while (next.length);
  return this;
}

function node_eachBefore(callback) {
  var node = this, nodes = [node], children, i;
  while (node = nodes.pop()) {
    callback(node), children = node.children;
    if (children) for (i = children.length - 1; i >= 0; --i) {
      nodes.push(children[i]);
    }
  }
  return this;
}

function node_eachAfter(callback) {
  var node = this, nodes = [node], next = [], children, i, n;
  while (node = nodes.pop()) {
    next.push(node), children = node.children;
    if (children) for (i = 0, n = children.length; i < n; ++i) {
      nodes.push(children[i]);
    }
  }
  while (node = next.pop()) {
    callback(node);
  }
  return this;
}

function node_sum(value) {
  return this.eachAfter(function(node) {
    var sum = +value(node.data) || 0,
        children = node.children,
        i = children && children.length;
    while (--i >= 0) sum += children[i].value;
    node.value = sum;
  });
}

function node_sort(compare) {
  return this.eachBefore(function(node) {
    if (node.children) {
      node.children.sort(compare);
    }
  });
}

function node_path(end) {
  var start = this,
      ancestor = leastCommonAncestor(start, end),
      nodes = [start];
  while (start !== ancestor) {
    start = start.parent;
    nodes.push(start);
  }
  var k = nodes.length;
  while (end !== ancestor) {
    nodes.splice(k, 0, end);
    end = end.parent;
  }
  return nodes;
}

function leastCommonAncestor(a, b) {
  if (a === b) return a;
  var aNodes = a.ancestors(),
      bNodes = b.ancestors(),
      c = null;
  a = aNodes.pop();
  b = bNodes.pop();
  while (a === b) {
    c = a;
    a = aNodes.pop();
    b = bNodes.pop();
  }
  return c;
}

function node_ancestors() {
  var node = this, nodes = [node];
  while (node = node.parent) {
    nodes.push(node);
  }
  return nodes;
}

function node_descendants() {
  var nodes = [];
  this.each(function(node) {
    nodes.push(node);
  });
  return nodes;
}

function node_leaves() {
  var leaves = [];
  this.eachBefore(function(node) {
    if (!node.children) {
      leaves.push(node);
    }
  });
  return leaves;
}

function node_links() {
  var root = this, links = [];
  root.each(function(node) {
    if (node !== root) { // Don’t include the root’s parent, if any.
      links.push({source: node.parent, target: node});
    }
  });
  return links;
}

function hierarchy(data, children) {
  var root = new Node(data),
      valued = +data.value && (root.value = data.value),
      node,
      nodes = [root],
      child,
      childs,
      i,
      n;

  if (children == null) children = defaultChildren;

  while (node = nodes.pop()) {
    if (valued) node.value = +node.data.value;
    if ((childs = children(node.data)) && (n = childs.length)) {
      node.children = new Array(n);
      for (i = n - 1; i >= 0; --i) {
        nodes.push(child = node.children[i] = new Node(childs[i]));
        child.parent = node;
        child.depth = node.depth + 1;
      }
    }
  }

  return root.eachBefore(computeHeight);
}

function node_copy() {
  return hierarchy(this).eachBefore(copyData);
}

function defaultChildren(d) {
  return d.children;
}

function copyData(node) {
  node.data = node.data.data;
}

function computeHeight(node) {
  var height = 0;
  do node.height = height;
  while ((node = node.parent) && (node.height < ++height));
}

function Node(data) {
  this.data = data;
  this.depth =
  this.height = 0;
  this.parent = null;
}

Node.prototype = hierarchy.prototype = {
  constructor: Node,
  count: node_count,
  each: node_each,
  eachAfter: node_eachAfter,
  eachBefore: node_eachBefore,
  sum: node_sum,
  sort: node_sort,
  path: node_path,
  ancestors: node_ancestors,
  descendants: node_descendants,
  leaves: node_leaves,
  links: node_links,
  copy: node_copy
};

var version = "2.1.1";

const defaults = Object.freeze({
  children: data => data.children,
  nodeSize: node => node.data.size,
  spacing: 0,
});

// Create a layout function with customizable options. Per D3-style, the
// options can be set at any time using setter methods. The layout function
// will compute the tree node positions based on the options in effect at the
// time it is called.
function flextree(options) {
  const opts = Object.assign({}, defaults, options);
  function accessor(name) {
    const opt = opts[name];
    return typeof opt === 'function' ? opt : () => opt;
  }

  function layout(tree) {
    const wtree = wrap(getWrapper(), tree, node=>node.children);
    wtree.update();
    return wtree.data;
  }

  function getFlexNode() {
    const nodeSize = accessor('nodeSize');
    const spacing = accessor('spacing');
    return class FlexNode extends hierarchy.prototype.constructor {
      constructor(data) {
        super(data);
      }
      copy() {
        const c = wrap(this.constructor, this, node=>node.children);
        c.each(node => node.data = node.data.data);
        return c;
      }
      get size() { return nodeSize(this); }
      spacing(oNode) { return spacing(this, oNode); }
      get nodes() { return this.descendants(); }
      get xSize() { return this.size[0]; }
      get ySize() { return this.size[1]; }
      get top() { return this.y; }
      get bottom() { return this.y + this.ySize; }
      get left() { return this.x - this.xSize / 2; }
      get right() { return this.x + this.xSize / 2; }
      get root() {
        const ancs = this.ancestors();
        return ancs[ancs.length - 1];
      }
      get numChildren() {
        return this.hasChildren ? this.children.length : 0;
      }
      get hasChildren() { return !this.noChildren; }
      get noChildren() { return this.children === null; }
      get firstChild() {
        return this.hasChildren ? this.children[0] : null;
      }
      get lastChild() {
        return this.hasChildren ? this.children[this.numChildren - 1] : null;
      }
      get extents() {
        return (this.children || []).reduce(
          (acc, kid) => FlexNode.maxExtents(acc, kid.extents),
          this.nodeExtents);
      }
      get nodeExtents() {
        return {
          top: this.top,
          bottom: this.bottom,
          left: this.left,
          right: this.right,
        };
      }
      static maxExtents(e0, e1) {
        return {
          top: Math.min(e0.top, e1.top),
          bottom: Math.max(e0.bottom, e1.bottom),
          left: Math.min(e0.left, e1.left),
          right: Math.max(e0.right, e1.right),
        };
      }
    };
  }

  function getWrapper() {
    const FlexNode = getFlexNode();
    const nodeSize = accessor('nodeSize');
    const spacing = accessor('spacing');
    return class extends FlexNode {
      constructor(data) {
        super(data);
        Object.assign(this, {
          x: 0, y: 0,
          relX: 0, prelim: 0, shift: 0, change: 0,
          lExt: this, lExtRelX: 0, lThr: null,
          rExt: this, rExtRelX: 0, rThr: null,
        });
      }
      get size() { return nodeSize(this.data); }
      spacing(oNode) { return spacing(this.data, oNode.data); }
      get x() { return this.data.x; }
      set x(v) { this.data.x = v; }
      get y() { return this.data.y; }
      set y(v) { this.data.y = v; }
      update() {
        layoutChildren(this);
        resolveX(this);
        return this;
      }
    };
  }

  function wrap(FlexClass, treeData, children) {
    const _wrap = (data, parent) => {
      const node = new FlexClass(data);
      Object.assign(node, {
        parent,
        depth: parent === null ? 0 : parent.depth + 1,
        height: 0,
        length: 1,
      });
      const kidsData = children(data) || [];
      node.children = kidsData.length === 0 ? null
        : kidsData.map(kd => _wrap(kd, node));
      if (node.children) {
        Object.assign(node, node.children.reduce(
          (hl, kid) => ({
            height: Math.max(hl.height, kid.height + 1),
            length: hl.length + kid.length,
          }), node
        ));
      }
      return node;
    };
    return _wrap(treeData, null);
  }


  Object.assign(layout, {
    nodeSize(arg) {
      return arguments.length ? (opts.nodeSize = arg, layout) : opts.nodeSize;
    },
    spacing(arg) {
      return arguments.length ? (opts.spacing = arg, layout) : opts.spacing;
    },
    children(arg) {
      return arguments.length ? (opts.children = arg, layout) : opts.children;
    },
    hierarchy(treeData, children) {
      const kids = typeof children === 'undefined' ? opts.children : children;
      return wrap(getFlexNode(), treeData, kids);
    },
    dump(tree) {
      const nodeSize = accessor('nodeSize');
      const _dump = i0 => node => {
        const i1 = i0 + '  ';
        const i2 = i0 + '    ';
        const {x, y} = node;
        const size = nodeSize(node);
        const kids = (node.children || []);
        const kdumps = (kids.length === 0) ? ' ' :
          `,${i1}children: [${i2}${kids.map(_dump(i2)).join(i2)}${i1}],${i0}`;
        return `{ size: [${size.join(', ')}],${i1}x: ${x}, y: ${y}${kdumps}},`;
      };
      return _dump('\n')(tree);
    },
  });
  return layout;
}
flextree.version = version;

const layoutChildren = (w, y = 0) => {
  w.y = y;
  (w.children || []).reduce((acc, kid) => {
    const [i, lastLows] = acc;
    layoutChildren(kid, w.y + w.ySize);
    // The lowest vertical coordinate while extreme nodes still point
    // in current subtree.
    const lowY = (i === 0 ? kid.lExt : kid.rExt).bottom;
    if (i !== 0) separate(w, i, lastLows);
    const lows = updateLows(lowY, i, lastLows);
    return [i + 1, lows];
  }, [0, null]);
  shiftChange(w);
  positionRoot(w);
  return w;
};

// Resolves the relative coordinate properties - relX and prelim --
// to set the final, absolute x coordinate for each node. This also sets
// `prelim` to 0, so that `relX` for each node is its x-coordinate relative
// to its parent.
const resolveX = (w, prevSum, parentX) => {
  // A call to resolveX without arguments is assumed to be for the root of
  // the tree. This will set the root's x-coord to zero.
  if (typeof prevSum === 'undefined') {
    prevSum = -w.relX - w.prelim;
    parentX = 0;
  }
  const sum = prevSum + w.relX;
  w.relX = sum + w.prelim - parentX;
  w.prelim = 0;
  w.x = parentX + w.relX;
  (w.children || []).forEach(k => resolveX(k, sum, w.x));
  return w;
};

// Process shift and change for all children, to add intermediate spacing to
// each child's modifier.
const shiftChange = w => {
  (w.children || []).reduce((acc, child) => {
    const [lastShiftSum, lastChangeSum] = acc;
    const shiftSum = lastShiftSum + child.shift;
    const changeSum = lastChangeSum + shiftSum + child.change;
    child.relX += changeSum;
    return [shiftSum, changeSum];
  }, [0, 0]);
};

// Separates the latest child from its previous sibling
/* eslint-disable complexity */
const separate = (w, i, lows) => {
  const lSib = w.children[i - 1];
  const curSubtree = w.children[i];
  let rContour = lSib;
  let rSumMods = lSib.relX;
  let lContour = curSubtree;
  let lSumMods = curSubtree.relX;
  let isFirst = true;
  while (rContour && lContour) {
    if (rContour.bottom > lows.lowY) lows = lows.next;
    // How far to the left of the right side of rContour is the left side
    // of lContour? First compute the center-to-center distance, then add
    // the "spacing"
    const dist =
      (rSumMods + rContour.prelim) - (lSumMods + lContour.prelim) +
      rContour.xSize / 2 + lContour.xSize / 2 +
      rContour.spacing(lContour);
    if (dist > 0 || (dist < 0 && isFirst)) {
      lSumMods += dist;
      // Move subtree by changing relX.
      moveSubtree(curSubtree, dist);
      distributeExtra(w, i, lows.index, dist);
    }
    isFirst = false;
    // Advance highest node(s) and sum(s) of modifiers
    const rightBottom = rContour.bottom;
    const leftBottom = lContour.bottom;
    if (rightBottom <= leftBottom) {
      rContour = nextRContour(rContour);
      if (rContour) rSumMods += rContour.relX;
    }
    if (rightBottom >= leftBottom) {
      lContour = nextLContour(lContour);
      if (lContour) lSumMods += lContour.relX;
    }
  }
  // Set threads and update extreme nodes. In the first case, the
  // current subtree is taller than the left siblings.
  if (!rContour && lContour) setLThr(w, i, lContour, lSumMods);
  // In the next case, the left siblings are taller than the current subtree
  else if (rContour && !lContour) setRThr(w, i, rContour, rSumMods);
};
/* eslint-enable complexity */

// Move subtree by changing relX.
const moveSubtree = (subtree, distance) => {
  subtree.relX += distance;
  subtree.lExtRelX += distance;
  subtree.rExtRelX += distance;
};

const distributeExtra = (w, curSubtreeI, leftSibI, dist) => {
  const curSubtree = w.children[curSubtreeI];
  const n = curSubtreeI - leftSibI;
  // Are there intermediate children?
  if (n > 1) {
    const delta = dist / n;
    w.children[leftSibI + 1].shift += delta;
    curSubtree.shift -= delta;
    curSubtree.change -= dist - delta;
  }
};

const nextLContour = w => {
  return w.hasChildren ? w.firstChild : w.lThr;
};

const nextRContour = w => {
  return w.hasChildren ? w.lastChild : w.rThr;
};

const setLThr = (w, i, lContour, lSumMods) => {
  const firstChild = w.firstChild;
  const lExt = firstChild.lExt;
  const curSubtree = w.children[i];
  lExt.lThr = lContour;
  // Change relX so that the sum of modifier after following thread is correct.
  const diff = lSumMods - lContour.relX - firstChild.lExtRelX;
  lExt.relX += diff;
  // Change preliminary x coordinate so that the node does not move.
  lExt.prelim -= diff;
  // Update extreme node and its sum of modifiers.
  firstChild.lExt = curSubtree.lExt;
  firstChild.lExtRelX = curSubtree.lExtRelX;
};

// Mirror image of setLThr.
const setRThr = (w, i, rContour, rSumMods) => {
  const curSubtree = w.children[i];
  const rExt = curSubtree.rExt;
  const lSib = w.children[i - 1];
  rExt.rThr = rContour;
  const diff = rSumMods - rContour.relX - curSubtree.rExtRelX;
  rExt.relX += diff;
  rExt.prelim -= diff;
  curSubtree.rExt = lSib.rExt;
  curSubtree.rExtRelX = lSib.rExtRelX;
};

// Position root between children, taking into account their modifiers
const positionRoot = w => {
  if (w.hasChildren) {
    const k0 = w.firstChild;
    const kf = w.lastChild;
    const prelim = (k0.prelim + k0.relX - k0.xSize / 2 +
      kf.relX + kf.prelim + kf.xSize / 2 ) / 2;
    Object.assign(w, {
      prelim,
      lExt: k0.lExt, lExtRelX: k0.lExtRelX,
      rExt: kf.rExt, rExtRelX: kf.rExtRelX,
    });
  }
};

// Make/maintain a linked list of the indexes of left siblings and their
// lowest vertical coordinate.
const updateLows = (lowY, index, lastLows) => {
  // Remove siblings that are hidden by the new subtree.
  while (lastLows !== null && lowY >= lastLows.lowY)
    lastLows = lastLows.next;
  // Prepend the new subtree.
  return {
    lowY,
    index,
    next: lastLows,
  };
};

const uniqId = Math.random().toString(36).slice(2, 8);
let globalIndex = 0;
function getId() {
  globalIndex += 1;
  return `mm-${uniqId}-${globalIndex}`;
}
function walkTree(tree, callback, key = 'c') {
  const walk = (item, parent) => callback(item, () => {
    var _item$key;

    (_item$key = item[key]) == null ? void 0 : _item$key.forEach(child => {
      walk(child, item);
    });
  }, parent);

  walk(tree);
}
function arrayFrom(arrayLike) {
  if (Array.from) return Array.from(arrayLike);
  const array = [];

  for (let i = 0; i < arrayLike.length; i += 1) {
    array.push(arrayLike[i]);
  }

  return array;
}
function flatMap(arrayLike, callback) {
  if (arrayLike.flatMap) return arrayLike.flatMap(callback);
  const array = [];

  for (let i = 0; i < arrayLike.length; i += 1) {
    const result = callback(arrayLike[i], i, arrayLike);
    if (Array.isArray(result)) array.push(...result);else array.push(result);
  }

  return array;
}
function addClass(className, ...rest) {
  const classList = (className || '').split(' ').filter(Boolean);
  rest.forEach(item => {
    if (item && classList.indexOf(item) < 0) classList.push(item);
  });
  return classList.join(' ');
}
function childSelector(filter) {
  if (typeof filter === 'string') {
    const tagName = filter;

    filter = el => el.tagName === tagName;
  }

  const filterFn = filter;
  return function selector() {
    let nodes = arrayFrom(this.childNodes);
    if (filterFn) nodes = nodes.filter(node => filterFn(node));
    return nodes;
  };
}

function memoize(fn) {
  const cache = {};
  return function memoized(...args) {
    const key = `${args[0]}`;
    let data = cache[key];

    if (!data) {
      data = {
        value: fn(...args)
      };
      cache[key] = data;
    }

    return data.value;
  };
}

function createElement(tagName, props, attrs) {
  const el = document.createElement(tagName);

  if (props) {
    Object.entries(props).forEach(([key, value]) => {
      el[key] = value;
    });
  }

  if (attrs) {
    Object.entries(attrs).forEach(([key, value]) => {
      el.setAttribute(key, value);
    });
  }

  return el;
}

const memoizedPreloadJS = memoize(url => {
  document.head.append(createElement('link', {
    rel: 'preload',
    as: 'script',
    href: url
  }));
});

function loadJSItem(item, context) {
  if (item.type === 'script') {
    return new Promise((resolve, reject) => {
      document.head.append(createElement('script', Object.assign(Object.assign({}, item.data), {}, {
        onload: resolve,
        onerror: reject
      })));
    });
  } else if (item.type === 'iife') {
    const {
      fn,
      getParams
    } = item.data;
    fn(...((getParams == null ? void 0 : getParams(context)) || []));
  }
}

function loadCSSItem(item) {
  if (item.type === 'style') {
    document.head.append(createElement('style', {
      textContent: item.data
    }));
  } else if (item.type === 'stylesheet') {
    document.head.append(createElement('link', Object.assign({
      rel: 'stylesheet'
    }, item.data)));
  }
}

async function loadJS(items, options) {
  const needPreload = items.filter(item => item.type === 'script');
  if (needPreload.length > 1) needPreload.forEach(item => memoizedPreloadJS(item.data.src));

  for (const item of items) {
    await loadJSItem(item, options);
  }
}
function loadCSS(items) {
  for (const item of items) {
    loadCSSItem(item);
  }
}
async function initializePlugins(Markmap, plugins, options) {
  options = Object.assign({}, options);
  await Promise.all(plugins.map(plugin => {
    loadCSS(plugin.styles);
    return loadJS(plugin.scripts, options);
  }));

  for (const {
    initialize
  } of plugins) {
    if (initialize) initialize(Markmap, options);
  }
}

const styles = [];
const scripts = [{
  type: 'iife',
  data: {
    fn: mathJax => {
      mathJax.options = Object.assign({
        skipHtmlTags: {
          '[-]': ['code', 'pre']
        }
      }, mathJax.options);
      mathJax.startup = Object.assign({
        typeset: false
      }, mathJax.startup);
      window.MathJax = mathJax;
    },
    getParams: context => [Object.assign({}, context.mathJax)]
  }
}, {
  type: 'script',
  data: {
    src: 'https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-svg.js'
  }
}];

function initialize(Markmap, options) {
  Markmap.transformHtml.tap((mm, nodes) => {
    var _MathJax$typeset, _MathJax;

    (_MathJax$typeset = (_MathJax = window.MathJax).typeset) == null ? void 0 : _MathJax$typeset.call(_MathJax, nodes);
  });
}

const plugin = {
  styles,
  scripts,
  initialize
};

const styles$1 = [{
  type: 'stylesheet',
  data: {
    href: 'https://cdn.jsdelivr.net/npm/prismjs@1/themes/prism.css'
  }
}];
const scripts$1 = [{
  type: 'iife',
  data: {
    fn: () => {
      window.Prism = {
        manual: true
      };
    }
  }
}, {
  type: 'script',
  data: {
    src: 'https://cdn.jsdelivr.net/npm/prismjs@1/components/prism-core.min.js'
  }
}, // components will be added by paths relative to path of autoloader
{
  type: 'script',
  data: {
    src: 'https://cdn.jsdelivr.net/npm/prismjs@1/plugins/autoloader/prism-autoloader.min.js'
  }
}];

function initialize$1(Markmap, options) {
  Markmap.transformHtml.tap((mm, nodes) => {
    const {
      Prism
    } = window;
    const langs = flatMap(nodes, node => arrayFrom(node.querySelectorAll('code[class*=language-]'))).map(code => {
      const lang = code.className.match(/(?:^|\s)language-(\S+)|$/)[1];

      if (Prism.languages[lang]) {
        Prism.highlightElement(code);
      } else {
        return lang;
      }
    }).filter(Boolean);

    if (langs.length) {
      Prism.plugins.autoloader.loadLanguages(langs, () => {
        mm.setData();
        mm.fit();
      });
    }
  });
}

const plugin$1 = {
  styles: styles$1,
  scripts: scripts$1,
  initialize: initialize$1
};

var plugins = /*#__PURE__*/Object.freeze({
__proto__: null,
mathJax: plugin,
prism: plugin$1
});

class Hook {
  constructor() {
    this.listeners = [];
  }

  tap(fn) {
    this.listeners.push(fn);
  }

  call(...args) {
    for (const fn of this.listeners) {
      fn(...args);
    }
  }

}

function linkWidth(nodeData) {
  const data = nodeData.data;
  return Math.max(6 - 2 * data.d, 1.5);
}

function adjustSpacing(tree, spacing) {
  walkTree(tree, (d, next) => {
    d.ySizeInner = d.ySize - spacing;
    d.y += spacing;
    next();
  }, 'children');
}

class Markmap {
  constructor(svg, opts) {
    this.options = void 0;
    this.state = void 0;
    this.svg = void 0;
    this.styleNode = void 0;
    this.g = void 0;
    this.zoom = void 0;
    ['handleZoom', 'handleClick'].forEach(key => {
      this[key] = this[key].bind(this);
    });
    this.svg = svg.datum ? svg : d3.select(svg);
    this.styleNode = this.svg.append('style');
    this.zoom = d3.zoom().on('zoom', this.handleZoom);
    this.options = Object.assign({
      duration: 500,
      nodeFont: '300 16px/20px sans-serif',
      nodeMinHeight: 16,
      spacingVertical: 5,
      spacingHorizontal: 80,
      autoFit: false,
      fitRatio: 0.95,
      color: (colorFn => node => colorFn(node.p.i))(d3.scaleOrdinal(d3.schemeCategory10)),
      paddingX: 8
    }, opts);
    this.state = {
      id: this.options.id || getId()
    };
    this.g = this.svg.append('g').attr('class', `${this.state.id}-g`);
    this.updateStyle();
    this.svg.call(this.zoom);
  }

  getStyleContent() {
    const {
      style,
      nodeFont
    } = this.options;
    const {
      id
    } = this.state;
    const extraStyle = typeof style === 'function' ? style(id) : '';
    const styleText = `\
.${id} a { color: #0097e6; }
.${id} a:hover { color: #00a8ff; }
.${id}-g > path { fill: none; }
.${id}-fo > div { font: ${nodeFont}; white-space: nowrap; }
.${id}-fo code { padding: .2em .4em; font-size: calc(1em - 2px); color: #555; background-color: #f0f0f0; border-radius: 2px; }
.${id}-fo del { text-decoration: line-through; }
.${id}-fo em { font-style: italic; }
.${id}-fo strong { font-weight: 500; }
.${id}-fo pre { margin: 0; }
.${id}-fo pre[class*=language-] { padding: 0; }
.${id}-g > g { cursor: pointer; }
${extraStyle}
`;
    return styleText;
  }

  updateStyle() {
    this.svg.attr('class', addClass(this.svg.attr('class'), this.state.id));
    this.styleNode.text(this.getStyleContent());
  }

  handleZoom() {
    const {
      transform
    } = d3.event;
    this.g.attr('transform', transform);
  }

  handleClick(d) {
    var _data$p;

    const {
      data
    } = d;
    data.p = Object.assign(Object.assign({}, data.p), {}, {
      f: !((_data$p = data.p) == null ? void 0 : _data$p.f)
    });
    this.renderData(d.data);
  }

  initializeData(node) {
    let i = 0;
    const {
      nodeFont,
      color,
      nodeMinHeight
    } = this.options;
    const {
      id
    } = this.state;
    const container = document.createElement('div');
    const containerClass = `${id}-container`;
    container.className = addClass(container.className, `${id}-fo`, containerClass);
    const style = document.createElement('style');
    style.textContent = `
${this.getStyleContent()}
.${containerClass} {
  position: absolute;
  width: 0;
  height: 0;
  top: -100px;
  left: -100px;
  overflow: hidden;
  font: ${nodeFont};
}
.${containerClass} > div {
  display: inline-block;
}
`;
    document.body.append(style, container);
    walkTree(node, (item, next) => {
      var _item$c;

      item.c = (_item$c = item.c) == null ? void 0 : _item$c.map(child => Object.assign({}, child));
      i += 1;
      const el = document.createElement('div');
      el.innerHTML = item.v;
      container.append(el);
      item.p = Object.assign(Object.assign({}, item.p), {}, {
        // unique ID
        i,
        el
      });
      color(item); // preload colors

      next();
    });
    const nodes = arrayFrom(container.childNodes);
    this.constructor.transformHtml.call(this, nodes);
    walkTree(node, (item, next, parent) => {
      var _parent$p;

      const rect = item.p.el.getBoundingClientRect();
      item.v = item.p.el.innerHTML;
      item.p.s = [Math.ceil(rect.width), Math.max(Math.ceil(rect.height), nodeMinHeight)]; // TODO keep keys for unchanged objects
      // unique key, should be based on content

      item.p.k = `${(parent == null ? void 0 : (_parent$p = parent.p) == null ? void 0 : _parent$p.i) || ''}.${item.p.i}:${item.v}`;
      next();
    });
    container.remove();
    style.remove();
  }

  setOptions(opts) {
    Object.assign(this.options, opts);
  }

  setData(data, opts) {
    if (!data) data = Object.assign({}, this.state.data);
    this.state.data = data;
    this.initializeData(data);
    if (opts) this.setOptions(opts);
    this.renderData();
  }

  renderData(originData) {
    var _origin$data$p$x, _origin$data$p$y;

    if (!this.state.data) return;
    const {
      spacingHorizontal,
      paddingX,
      spacingVertical,
      autoFit,
      color
    } = this.options;
    const {
      id
    } = this.state;
    const layout = flextree().children(d => {
      var _d$p;

      return !((_d$p = d.p) == null ? void 0 : _d$p.f) && d.c;
    }).nodeSize(d => {
      const [width, height] = d.data.p.s;
      return [height, width + (width ? paddingX * 2 : 0) + spacingHorizontal];
    }).spacing((a, b) => {
      return a.parent === b.parent ? spacingVertical : spacingVertical * 2;
    });
    const tree = layout.hierarchy(this.state.data);
    layout(tree);
    adjustSpacing(tree, spacingHorizontal);
    const descendants = tree.descendants().reverse();
    const links = tree.links();
    const linkShape = d3.linkHorizontal();
    const minX = d3.min(descendants, d => d.x - d.xSize / 2);
    const maxX = d3.max(descendants, d => d.x + d.xSize / 2);
    const minY = d3.min(descendants, d => d.y);
    const maxY = d3.max(descendants, d => d.y + d.ySizeInner);
    Object.assign(this.state, {
      minX,
      maxX,
      minY,
      maxY
    });
    if (autoFit) this.fit();
    const origin = originData && descendants.find(item => item.data === originData) || tree;
    const x0 = (_origin$data$p$x = origin.data.p.x0) != null ? _origin$data$p$x : origin.x;
    const y0 = (_origin$data$p$y = origin.data.p.y0) != null ? _origin$data$p$y : origin.y; // Update the nodes

    const node = this.g.selectAll(childSelector('g')).data(descendants, d => d.data.p.k);
    const nodeEnter = node.enter().append('g').attr('transform', d => `translate(${y0 + origin.ySizeInner - d.ySizeInner},${x0 + origin.xSize / 2 - d.xSize})`).on('click', this.handleClick);
    const nodeExit = this.transition(node.exit());
    nodeExit.select('rect').attr('width', 0).attr('x', d => d.ySizeInner);
    nodeExit.select('foreignObject').style('opacity', 0);
    nodeExit.attr('transform', d => `translate(${origin.y + origin.ySizeInner - d.ySizeInner},${origin.x + origin.xSize / 2 - d.xSize})`).remove();
    const nodeMerge = node.merge(nodeEnter);
    this.transition(nodeMerge).attr('transform', d => `translate(${d.y},${d.x - d.xSize / 2})`);
    const rect = nodeMerge.selectAll(childSelector('rect')).data(d => [d], d => d.data.p.k).join(enter => {
      return enter.append('rect').attr('x', d => d.ySizeInner).attr('y', d => d.xSize - linkWidth(d) / 2).attr('width', 0).attr('height', linkWidth);
    }, update => update, exit => exit.remove());
    this.transition(rect).attr('x', -1).attr('width', d => d.ySizeInner + 2).attr('fill', d => color(d.data));
    const circle = nodeMerge.selectAll(childSelector('circle')).data(d => d.data.c ? [d] : [], d => d.data.p.k).join(enter => {
      return enter.append('circle').attr('stroke-width', '1.5').attr('cx', d => d.ySizeInner).attr('cy', d => d.xSize).attr('r', 0);
    }, update => update, exit => exit.remove());
    this.transition(circle).attr('r', 6).attr('stroke', d => color(d.data)).attr('fill', d => {
      var _d$data$p;

      return ((_d$data$p = d.data.p) == null ? void 0 : _d$data$p.f) ? color(d.data) : '#fff';
    });
    const foreignObject = nodeMerge.selectAll(childSelector('foreignObject')).data(d => [d], d => d.data.p.k).join(enter => {
      const fo = enter.append('foreignObject').attr('class', `${id}-fo`).attr('x', paddingX).attr('y', 0).style('opacity', 0).attr('height', d => d.xSize);
      fo.append('xhtml:div').select(function (d) {
        const node = d.data.p.el.cloneNode(true);
        this.replaceWith(node);
        return node;
      }).attr('xmlns', 'http://www.w3.org/1999/xhtml');
      return fo;
    }, update => update, exit => exit.remove()).attr('width', d => Math.max(0, d.ySizeInner - paddingX * 2));
    this.transition(foreignObject).style('opacity', 1); // Update the links

    const path = this.g.selectAll(childSelector('path')).data(links, d => d.target.data.p.k).join(enter => {
      const source = [y0 + origin.ySizeInner, x0 + origin.xSize / 2];
      return enter.insert('path', 'g').attr('d', linkShape({
        source,
        target: source
      }));
    }, update => update, exit => {
      const source = [origin.y + origin.ySizeInner, origin.x + origin.xSize / 2];
      return this.transition(exit).attr('d', linkShape({
        source,
        target: source
      })).remove();
    });
    this.transition(path).attr('stroke', d => color(d.target.data)).attr('stroke-width', d => linkWidth(d.target)).attr('d', d => {
      const source = [d.source.y + d.source.ySizeInner, d.source.x + d.source.xSize / 2];
      const target = [d.target.y, d.target.x + d.target.xSize / 2];
      return linkShape({
        source,
        target
      });
    });
    descendants.forEach(d => {
      d.data.p.x0 = d.x;
      d.data.p.y0 = d.y;
    });
  }

  transition(sel) {
    const {
      duration
    } = this.options;
    return sel.transition().duration(duration);
  }

  fit() {
    const svgNode = this.svg.node();
    const {
      width: offsetWidth,
      height: offsetHeight
    } = svgNode.getBoundingClientRect();
    const {
      fitRatio
    } = this.options;
    const {
      minX,
      maxX,
      minY,
      maxY
    } = this.state;
    const naturalWidth = maxY - minY;
    const naturalHeight = maxX - minX;
    const scale = Math.min(offsetWidth / naturalWidth * fitRatio, offsetHeight / naturalHeight * fitRatio, 2);
    const initialZoom = d3.zoomIdentity.translate((offsetWidth - naturalWidth * scale) / 2 - minY * scale, (offsetHeight - naturalHeight * scale) / 2 - minX * scale).scale(scale);
    return this.transition(this.svg).call(this.zoom.transform, initialZoom).end();
  }

  rescale(scale) {
    const svgNode = this.svg.node();
    const {
      width: offsetWidth,
      height: offsetHeight
    } = svgNode.getBoundingClientRect();
    const halfWidth = offsetWidth / 2;
    const halfHeight = offsetHeight / 2;
    const transform = d3.zoomTransform(svgNode);
    const newTransform = transform.translate((halfWidth - transform.x) * (1 - scale) / transform.k, (halfHeight - transform.y) * (1 - scale) / transform.k).scale(scale);
    return this.transition(this.svg).call(this.zoom.transform, newTransform).end();
  }

  static create(svg, opts, data) {
    const mm = new Markmap(svg, opts);

    if (data) {
      mm.setData(data);
      mm.fit(); // always fit for the first render
    }

    return mm;
  }

}
Markmap.transformHtml = new Hook();
function markmap(svg, data, opts) {
  return Markmap.create(svg, opts, data);
}
async function loadPlugins(items, options) {
  items = items.map(item => {
    if (typeof item === 'string') {
      const name = item;
      item = plugins[name];

      if (!item) {
        console.warn(`[markmap] Unknown plugin: ${name}`);
      }
    }

    return item;
  }).filter(Boolean);
  return initializePlugins(Markmap, items, options);
}

exports.Markmap = Markmap;
exports.loadPlugins = loadPlugins;
exports.markmap = markmap;
exports.plugins = plugins;

}(this.markmap = this.markmap || {}, d3));
