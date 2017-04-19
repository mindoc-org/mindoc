module.exports =
/******/ (function(modules) { // webpackBootstrap
/******/ 	// The module cache
/******/ 	var installedModules = {};

/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {

/******/ 		// Check if module is in cache
/******/ 		if(installedModules[moduleId])
/******/ 			return installedModules[moduleId].exports;

/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = installedModules[moduleId] = {
/******/ 			exports: {},
/******/ 			id: moduleId,
/******/ 			loaded: false
/******/ 		};

/******/ 		// Execute the module function
/******/ 		modules[moduleId].call(module.exports, module, module.exports, __webpack_require__);

/******/ 		// Flag the module as loaded
/******/ 		module.loaded = true;

/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}


/******/ 	// expose the modules object (__webpack_modules__)
/******/ 	__webpack_require__.m = modules;

/******/ 	// expose the module cache
/******/ 	__webpack_require__.c = installedModules;

/******/ 	// __webpack_public_path__
/******/ 	__webpack_require__.p = "/dist/";

/******/ 	// Load entry module and return exports
/******/ 	return __webpack_require__(0);
/******/ })
/************************************************************************/
/******/ ({

/***/ 0:
/***/ function(module, exports, __webpack_require__) {

	module.exports = __webpack_require__(49);


/***/ },

/***/ 3:
/***/ function(module, exports) {

	module.exports = function normalizeComponent (
	  rawScriptExports,
	  compiledTemplate,
	  scopeId,
	  cssModules
	) {
	  var esModule
	  var scriptExports = rawScriptExports = rawScriptExports || {}

	  // ES6 modules interop
	  var type = typeof rawScriptExports.default
	  if (type === 'object' || type === 'function') {
	    esModule = rawScriptExports
	    scriptExports = rawScriptExports.default
	  }

	  // Vue.extend constructor export interop
	  var options = typeof scriptExports === 'function'
	    ? scriptExports.options
	    : scriptExports

	  // render functions
	  if (compiledTemplate) {
	    options.render = compiledTemplate.render
	    options.staticRenderFns = compiledTemplate.staticRenderFns
	  }

	  // scopedId
	  if (scopeId) {
	    options._scopeId = scopeId
	  }

	  // inject cssModules
	  if (cssModules) {
	    var computed = options.computed || (options.computed = {})
	    Object.keys(cssModules).forEach(function (key) {
	      var module = cssModules[key]
	      computed[key] = function () { return module }
	    })
	  }

	  return {
	    esModule: esModule,
	    exports: scriptExports,
	    options: options
	  }
	}


/***/ },

/***/ 49:
/***/ function(module, exports, __webpack_require__) {

	'use strict';

	exports.__esModule = true;

	var _item = __webpack_require__(50);

	var _item2 = _interopRequireDefault(_item);

	function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

	/* istanbul ignore next */
	_item2.default.install = function (Vue) {
	  Vue.component(_item2.default.name, _item2.default);
	};

	exports.default = _item2.default;

/***/ },

/***/ 50:
/***/ function(module, exports, __webpack_require__) {

	var Component = __webpack_require__(3)(
	  /* script */
	  __webpack_require__(51),
	  /* template */
	  __webpack_require__(52),
	  /* scopeId */
	  null,
	  /* cssModules */
	  null
	)

	module.exports = Component.exports


/***/ },

/***/ 51:
/***/ function(module, exports) {

	'use strict';

	exports.__esModule = true;
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//

	var CARD_SCALE = 0.83;
	exports.default = {
	  name: 'ElCarouselItem',

	  props: {
	    name: String
	  },

	  data: function data() {
	    return {
	      hover: false,
	      translate: 0,
	      scale: 1,
	      active: false,
	      ready: false,
	      inStage: false
	    };
	  },


	  methods: {
	    processIndex: function processIndex(index, activeIndex, length) {
	      if (activeIndex === 0 && index === length - 1) {
	        return -1;
	      } else if (activeIndex === length - 1 && index === 0) {
	        return length;
	      } else if (index < activeIndex - 1 && activeIndex - index >= length / 2) {
	        return length + 1;
	      } else if (index > activeIndex + 1 && index - activeIndex >= length / 2) {
	        return -2;
	      }
	      return index;
	    },
	    calculateTranslate: function calculateTranslate(index, activeIndex, parentWidth) {
	      if (this.inStage) {
	        return parentWidth * ((2 - CARD_SCALE) * (index - activeIndex) + 1) / 4;
	      } else if (index < activeIndex) {
	        return -(1 + CARD_SCALE) * parentWidth / 4;
	      } else {
	        return (3 + CARD_SCALE) * parentWidth / 4;
	      }
	    },
	    translateItem: function translateItem(index, activeIndex) {
	      var parentWidth = this.$parent.$el.offsetWidth;
	      var length = this.$parent.items.length;
	      if (index !== activeIndex && length > 2) {
	        index = this.processIndex(index, activeIndex, length);
	      }
	      if (this.$parent.type === 'card') {
	        this.inStage = Math.round(Math.abs(index - activeIndex)) <= 1;
	        this.active = index === activeIndex;
	        this.translate = this.calculateTranslate(index, activeIndex, parentWidth);
	        this.scale = this.active ? 1 : CARD_SCALE;
	      } else {
	        this.active = index === activeIndex;
	        this.translate = parentWidth * (index - activeIndex);
	      }
	      this.ready = true;
	    },
	    handleItemClick: function handleItemClick() {
	      var parent = this.$parent;
	      if (parent && parent.type === 'card') {
	        var index = parent.items.indexOf(this);
	        parent.setActiveItem(index);
	      }
	    }
	  },

	  created: function created() {
	    this.$parent && this.$parent.handleItemChange();
	  },
	  destroyed: function destroyed() {
	    this.$parent && this.$parent.handleItemChange();
	  }
	};

/***/ },

/***/ 52:
/***/ function(module, exports) {

	module.exports={render:function (){var _vm=this;var _h=_vm.$createElement;var _c=_vm._self._c||_h;
	  return _c('div', {
	    directives: [{
	      name: "show",
	      rawName: "v-show",
	      value: (_vm.ready),
	      expression: "ready"
	    }],
	    staticClass: "el-carousel__item",
	    class: {
	      'is-active': _vm.active,
	      'el-carousel__item--card': _vm.$parent.type === 'card',
	        'is-in-stage': _vm.inStage,
	        'is-hover': _vm.hover
	    },
	    style: ({
	      msTransform: ("translateX(" + _vm.translate + "px) scale(" + _vm.scale + ")"),
	      webkitTransform: ("translateX(" + _vm.translate + "px) scale(" + _vm.scale + ")"),
	      transform: ("translateX(" + _vm.translate + "px) scale(" + _vm.scale + ")")
	    }),
	    on: {
	      "click": _vm.handleItemClick
	    }
	  }, [(_vm.$parent.type === 'card') ? _c('div', {
	    directives: [{
	      name: "show",
	      rawName: "v-show",
	      value: (!_vm.active),
	      expression: "!active"
	    }],
	    staticClass: "el-carousel__mask"
	  }) : _vm._e(), _vm._t("default")], 2)
	},staticRenderFns: []}

/***/ }

/******/ });