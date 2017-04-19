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

	module.exports = __webpack_require__(147);


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

/***/ 147:
/***/ function(module, exports, __webpack_require__) {

	'use strict';

	exports.__esModule = true;

	var _form = __webpack_require__(148);

	var _form2 = _interopRequireDefault(_form);

	function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

	/* istanbul ignore next */
	_form2.default.install = function (Vue) {
	  Vue.component(_form2.default.name, _form2.default);
	};

	exports.default = _form2.default;

/***/ },

/***/ 148:
/***/ function(module, exports, __webpack_require__) {

	var Component = __webpack_require__(3)(
	  /* script */
	  __webpack_require__(149),
	  /* template */
	  __webpack_require__(150),
	  /* scopeId */
	  null,
	  /* cssModules */
	  null
	)

	module.exports = Component.exports


/***/ },

/***/ 149:
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

	exports.default = {
	  name: 'ElForm',

	  componentName: 'ElForm',

	  props: {
	    model: Object,
	    rules: Object,
	    labelPosition: String,
	    labelWidth: String,
	    labelSuffix: {
	      type: String,
	      default: ''
	    },
	    inline: Boolean,
	    showMessage: {
	      type: Boolean,
	      default: true
	    }
	  },
	  watch: {
	    rules: function rules() {
	      this.validate();
	    }
	  },
	  data: function data() {
	    return {
	      fields: []
	    };
	  },
	  created: function created() {
	    var _this = this;

	    this.$on('el.form.addField', function (field) {
	      if (field) {
	        _this.fields.push(field);
	      }
	    });
	    /* istanbul ignore next */
	    this.$on('el.form.removeField', function (field) {
	      if (field.prop) {
	        _this.fields.splice(_this.fields.indexOf(field), 1);
	      }
	    });
	  },

	  methods: {
	    resetFields: function resetFields() {
	      this.fields.forEach(function (field) {
	        field.resetField();
	      });
	    },
	    validate: function validate(callback) {
	      var _this2 = this;

	      var valid = true;
	      var count = 0;
	      // 如果需要验证的fields为空，调用验证时立刻返回callback
	      if (this.fields.length === 0 && callback) {
	        callback(true);
	      }
	      this.fields.forEach(function (field, index) {
	        field.validate('', function (errors) {
	          if (errors) {
	            valid = false;
	          }
	          if (typeof callback === 'function' && ++count === _this2.fields.length) {
	            callback(valid);
	          }
	        });
	      });
	    },
	    validateField: function validateField(prop, cb) {
	      var field = this.fields.filter(function (field) {
	        return field.prop === prop;
	      })[0];
	      if (!field) {
	        throw new Error('must call validateField with valid prop string!');
	      }

	      field.validate('', cb);
	    }
	  }
	};

/***/ },

/***/ 150:
/***/ function(module, exports) {

	module.exports={render:function (){var _vm=this;var _h=_vm.$createElement;var _c=_vm._self._c||_h;
	  return _c('form', {
	    staticClass: "el-form",
	    class: [
	      _vm.labelPosition ? 'el-form--label-' + _vm.labelPosition : '', {
	        'el-form--inline': _vm.inline
	      }
	    ]
	  }, [_vm._t("default")], 2)
	},staticRenderFns: []}

/***/ }

/******/ });