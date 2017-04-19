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

	module.exports = __webpack_require__(71);


/***/ },

/***/ 71:
/***/ function(module, exports, __webpack_require__) {

	'use strict';

	exports.__esModule = true;

	var _col = __webpack_require__(72);

	var _col2 = _interopRequireDefault(_col);

	function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

	/* istanbul ignore next */
	_col2.default.install = function (Vue) {
	  Vue.component(_col2.default.name, _col2.default);
	};

	exports.default = _col2.default;

/***/ },

/***/ 72:
/***/ function(module, exports) {

	'use strict';

	exports.__esModule = true;

	var _typeof = typeof Symbol === "function" && typeof Symbol.iterator === "symbol" ? function (obj) { return typeof obj; } : function (obj) { return obj && typeof Symbol === "function" && obj.constructor === Symbol && obj !== Symbol.prototype ? "symbol" : typeof obj; };

	exports.default = {
	  name: 'ElCol',

	  props: {
	    span: {
	      type: Number,
	      default: 24
	    },
	    offset: Number,
	    pull: Number,
	    push: Number,
	    xs: [Number, Object],
	    sm: [Number, Object],
	    md: [Number, Object],
	    lg: [Number, Object]
	  },

	  computed: {
	    gutter: function gutter() {
	      return this.$parent.gutter;
	    },
	    style: function style() {
	      var ret = {};

	      if (this.gutter) {
	        ret.paddingLeft = this.gutter / 2 + 'px';
	        ret.paddingRight = ret.paddingLeft;
	      }

	      return ret;
	    }
	  },
	  render: function render(h) {
	    var _this = this;

	    var style = this.style;

	    var classList = [];

	    ['span', 'offset', 'pull', 'push'].forEach(function (prop) {
	      if (_this[prop]) {
	        classList.push(prop !== 'span' ? 'el-col-' + prop + '-' + _this[prop] : 'el-col-' + _this[prop]);
	      }
	    });

	    ['xs', 'sm', 'md', 'lg'].forEach(function (size) {
	      if (typeof _this[size] === 'number') {
	        classList.push('el-col-' + size + '-' + _this[size]);
	      } else if (_typeof(_this[size]) === 'object') {
	        (function () {
	          var props = _this[size];
	          Object.keys(props).forEach(function (prop) {
	            classList.push(prop !== 'span' ? 'el-col-' + size + '-' + prop + '-' + props[prop] : 'el-col-' + size + '-' + props[prop]);
	          });
	        })();
	      }
	    });

	    return h(
	      'div',
	      {
	        'class': ['el-col', classList],
	        style: style },
	      [this.$slots.default]
	    );
	  }
	};

/***/ }

/******/ });