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

	module.exports = __webpack_require__(166);


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

/***/ 9:
/***/ function(module, exports) {

	module.exports = require("element-ui/lib/input");

/***/ },

/***/ 117:
/***/ function(module, exports) {

	module.exports = require("element-ui/lib/utils/dom");

/***/ },

/***/ 166:
/***/ function(module, exports, __webpack_require__) {

	'use strict';

	exports.__esModule = true;

	var _inputNumber = __webpack_require__(167);

	var _inputNumber2 = _interopRequireDefault(_inputNumber);

	function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

	/* istanbul ignore next */
	_inputNumber2.default.install = function (Vue) {
	  Vue.component(_inputNumber2.default.name, _inputNumber2.default);
	};

	exports.default = _inputNumber2.default;

/***/ },

/***/ 167:
/***/ function(module, exports, __webpack_require__) {

	var Component = __webpack_require__(3)(
	  /* script */
	  __webpack_require__(168),
	  /* template */
	  __webpack_require__(169),
	  /* scopeId */
	  null,
	  /* cssModules */
	  null
	)

	module.exports = Component.exports


/***/ },

/***/ 168:
/***/ function(module, exports, __webpack_require__) {

	'use strict';

	exports.__esModule = true;

	var _input = __webpack_require__(9);

	var _input2 = _interopRequireDefault(_input);

	var _dom = __webpack_require__(117);

	function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

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

	exports.default = {
	  name: 'ElInputNumber',
	  directives: {
	    repeatClick: {
	      bind: function bind(el, binding, vnode) {
	        var interval = null;
	        var startTime = void 0;
	        var handler = function handler() {
	          return vnode.context[binding.expression].apply();
	        };
	        var clear = function clear() {
	          if (new Date() - startTime < 100) {
	            handler();
	          }
	          clearInterval(interval);
	          interval = null;
	        };

	        (0, _dom.on)(el, 'mousedown', function () {
	          startTime = new Date();
	          (0, _dom.once)(document, 'mouseup', clear);
	          clearInterval(interval);
	          interval = setInterval(handler, 100);
	        });
	      }
	    }
	  },
	  components: {
	    ElInput: _input2.default
	  },
	  props: {
	    step: {
	      type: Number,
	      default: 1
	    },
	    max: {
	      type: Number,
	      default: Infinity
	    },
	    min: {
	      type: Number,
	      default: -Infinity
	    },
	    value: {
	      default: 0
	    },
	    disabled: Boolean,
	    size: String,
	    controls: {
	      type: Boolean,
	      default: true
	    }
	  },
	  data: function data() {
	    return {
	      currentValue: 0
	    };
	  },

	  watch: {
	    value: {
	      immediate: true,
	      handler: function handler(value) {
	        var newVal = Number(value);
	        if (isNaN(newVal)) return;
	        if (newVal >= this.max) newVal = this.max;
	        if (newVal <= this.min) newVal = this.min;
	        this.currentValue = newVal;
	        this.$emit('input', newVal);
	      }
	    }
	  },
	  computed: {
	    minDisabled: function minDisabled() {
	      return this._decrease(this.value, this.step) < this.min;
	    },
	    maxDisabled: function maxDisabled() {
	      return this._increase(this.value, this.step) > this.max;
	    },
	    precision: function precision() {
	      var value = this.value,
	          step = this.step,
	          getPrecision = this.getPrecision;

	      return Math.max(getPrecision(value), getPrecision(step));
	    }
	  },
	  methods: {
	    toPrecision: function toPrecision(num, precision) {
	      if (precision === undefined) precision = this.precision;
	      return parseFloat(parseFloat(Number(num).toFixed(precision)));
	    },
	    getPrecision: function getPrecision(value) {
	      var valueString = value.toString();
	      var dotPosition = valueString.indexOf('.');
	      var precision = 0;
	      if (dotPosition !== -1) {
	        precision = valueString.length - dotPosition - 1;
	      }
	      return precision;
	    },
	    _increase: function _increase(val, step) {
	      if (typeof val !== 'number') return this.currentValue;

	      var precisionFactor = Math.pow(10, this.precision);

	      return this.toPrecision((precisionFactor * val + precisionFactor * step) / precisionFactor);
	    },
	    _decrease: function _decrease(val, step) {
	      if (typeof val !== 'number') return this.currentValue;

	      var precisionFactor = Math.pow(10, this.precision);

	      return this.toPrecision((precisionFactor * val - precisionFactor * step) / precisionFactor);
	    },
	    increase: function increase() {
	      if (this.disabled || this.maxDisabled) return;
	      var value = this.value || 0;
	      var newVal = this._increase(value, this.step);
	      if (newVal > this.max) return;
	      this.setCurrentValue(newVal);
	    },
	    decrease: function decrease() {
	      if (this.disabled || this.minDisabled) return;
	      var value = this.value || 0;
	      var newVal = this._decrease(value, this.step);
	      if (newVal < this.min) return;
	      this.setCurrentValue(newVal);
	    },
	    handleBlur: function handleBlur() {
	      this.$refs.input.setCurrentValue(this.currentValue);
	    },
	    setCurrentValue: function setCurrentValue(newVal) {
	      var oldVal = this.currentValue;
	      if (newVal >= this.max) newVal = this.max;
	      if (newVal <= this.min) newVal = this.min;
	      if (oldVal === newVal) return;
	      this.$emit('change', newVal, oldVal);
	      this.$emit('input', newVal);
	      this.currentValue = newVal;
	    },
	    handleInput: function handleInput(value) {
	      var newVal = Number(value);
	      if (!isNaN(newVal)) {
	        this.setCurrentValue(newVal);
	      }
	    }
	  }
	};

/***/ },

/***/ 169:
/***/ function(module, exports) {

	module.exports={render:function (){var _vm=this;var _h=_vm.$createElement;var _c=_vm._self._c||_h;
	  return _c('div', {
	    staticClass: "el-input-number",
	    class: [
	      _vm.size ? 'el-input-number--' + _vm.size : '', {
	        'is-disabled': _vm.disabled
	      }, {
	        'is-without-controls': !_vm.controls
	      }
	    ]
	  }, [(_vm.controls) ? _c('span', {
	    directives: [{
	      name: "repeat-click",
	      rawName: "v-repeat-click",
	      value: (_vm.decrease),
	      expression: "decrease"
	    }],
	    staticClass: "el-input-number__decrease",
	    class: {
	      'is-disabled': _vm.minDisabled
	    }
	  }, [_c('i', {
	    staticClass: "el-icon-minus"
	  })]) : _vm._e(), (_vm.controls) ? _c('span', {
	    directives: [{
	      name: "repeat-click",
	      rawName: "v-repeat-click",
	      value: (_vm.increase),
	      expression: "increase"
	    }],
	    staticClass: "el-input-number__increase",
	    class: {
	      'is-disabled': _vm.maxDisabled
	    }
	  }, [_c('i', {
	    staticClass: "el-icon-plus"
	  })]) : _vm._e(), _c('el-input', {
	    ref: "input",
	    attrs: {
	      "value": _vm.currentValue,
	      "disabled": _vm.disabled,
	      "size": _vm.size,
	      "max": _vm.max,
	      "min": _vm.min
	    },
	    on: {
	      "blur": _vm.handleBlur,
	      "input": _vm.handleInput
	    },
	    nativeOn: {
	      "keydown": [function($event) {
	        if (_vm._k($event.keyCode, "up", 38)) { return; }
	        $event.preventDefault();
	        _vm.increase($event)
	      }, function($event) {
	        if (_vm._k($event.keyCode, "down", 40)) { return; }
	        $event.preventDefault();
	        _vm.decrease($event)
	      }]
	    }
	  }, [(_vm.$slots.prepend) ? _c('template', {
	    slot: "prepend"
	  }, [_vm._t("prepend")], 2) : _vm._e(), (_vm.$slots.append) ? _c('template', {
	    slot: "append"
	  }, [_vm._t("append")], 2) : _vm._e()], 2)], 1)
	},staticRenderFns: []}

/***/ }

/******/ });