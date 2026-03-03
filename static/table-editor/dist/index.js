(function webpackUniversalModuleDefinition(root, factory) {
	if(typeof exports === 'object' && typeof module === 'object')
		module.exports = factory();
	else if(typeof define === 'function' && define.amd)
		define("TableEditor", [], factory);
	else if(typeof exports === 'object')
		exports["TableEditor"] = factory();
	else
		root["TableEditor"] = factory();
})(typeof self !== 'undefined' ? self : this, function() {
return /******/ (function(modules) { // webpackBootstrap
/******/ 	// The module cache
/******/ 	var installedModules = {};
/******/
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/
/******/ 		// Check if module is in cache
/******/ 		if(installedModules[moduleId]) {
/******/ 			return installedModules[moduleId].exports;
/******/ 		}
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = installedModules[moduleId] = {
/******/ 			i: moduleId,
/******/ 			l: false,
/******/ 			exports: {}
/******/ 		};
/******/
/******/ 		// Execute the module function
/******/ 		modules[moduleId].call(module.exports, module, module.exports, __webpack_require__);
/******/
/******/ 		// Flag the module as loaded
/******/ 		module.l = true;
/******/
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/
/******/
/******/ 	// expose the modules object (__webpack_modules__)
/******/ 	__webpack_require__.m = modules;
/******/
/******/ 	// expose the module cache
/******/ 	__webpack_require__.c = installedModules;
/******/
/******/ 	// define getter function for harmony exports
/******/ 	__webpack_require__.d = function(exports, name, getter) {
/******/ 		if(!__webpack_require__.o(exports, name)) {
/******/ 			Object.defineProperty(exports, name, {
/******/ 				configurable: false,
/******/ 				enumerable: true,
/******/ 				get: getter
/******/ 			});
/******/ 		}
/******/ 	};
/******/
/******/ 	// getDefaultExport function for compatibility with non-harmony modules
/******/ 	__webpack_require__.n = function(module) {
/******/ 		var getter = module && module.__esModule ?
/******/ 			function getDefault() { return module['default']; } :
/******/ 			function getModuleExports() { return module; };
/******/ 		__webpack_require__.d(getter, 'a', getter);
/******/ 		return getter;
/******/ 	};
/******/
/******/ 	// Object.prototype.hasOwnProperty.call
/******/ 	__webpack_require__.o = function(object, property) { return Object.prototype.hasOwnProperty.call(object, property); };
/******/
/******/ 	// __webpack_public_path__
/******/ 	__webpack_require__.p = "";
/******/
/******/ 	// Load entry module and return exports
/******/ 	return __webpack_require__(__webpack_require__.s = 0);
/******/ })
/************************************************************************/
/******/ ([
/* 0 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
Object.defineProperty(__webpack_exports__, "__esModule", { value: true });
/* harmony export (immutable) */ __webpack_exports__["initTableEditor"] = initTableEditor;
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__susisu_mte_kernel__ = __webpack_require__(1);
/* global CodeMirror, $ */


// port of the code from: https://github.com/susisu/mte-demo/blob/master/src/main.js

// text editor interface
// see https://doc.esdoc.org/github.com/susisu/mte-kernel/class/lib/text-editor.js~ITextEditor.html
class TextEditorInterface {
  constructor (editor) {
    this.editor = editor
    this.doc = editor.getDoc()
    this.transaction = false
    this.onDidFinishTransaction = null
  }

  getCursorPosition () {
    const { line, ch } = this.doc.getCursor()
    return new __WEBPACK_IMPORTED_MODULE_0__susisu_mte_kernel__["c" /* Point */](line, ch)
  }

  setCursorPosition (pos) {
    this.doc.setCursor({ line: pos.row, ch: pos.column })
  }

  setSelectionRange (range) {
    this.doc.setSelection(
      { line: range.start.row, ch: range.start.column },
      { line: range.end.row, ch: range.end.column }
    )
  }

  getLastRow () {
    return this.doc.lineCount() - 1
  }

  acceptsTableEdit () {
    return true
  }

  getLine (row) {
    return this.doc.getLine(row)
  }

  insertLine (row, line) {
    const lastRow = this.getLastRow()
    if (row > lastRow) {
      const lastLine = this.getLine(lastRow)
      this.doc.replaceRange(
        '\n' + line,
        { line: lastRow, ch: lastLine.length },
        { line: lastRow, ch: lastLine.length }
      )
    } else {
      this.doc.replaceRange(
        line + '\n',
        { line: row, ch: 0 },
        { line: row, ch: 0 }
      )
    }
  }

  deleteLine (row) {
    const lastRow = this.getLastRow()
    if (row >= lastRow) {
      if (lastRow > 0) {
        const preLastLine = this.getLine(lastRow - 1)
        const lastLine = this.getLine(lastRow)
        this.doc.replaceRange(
          '',
          { line: lastRow - 1, ch: preLastLine.length },
          { line: lastRow, ch: lastLine.length }
        )
      } else {
        const lastLine = this.getLine(lastRow)
        this.doc.replaceRange(
          '',
          { line: lastRow, ch: 0 },
          { line: lastRow, ch: lastLine.length }
        )
      }
    } else {
      this.doc.replaceRange(
        '',
        { line: row, ch: 0 },
        { line: row + 1, ch: 0 }
      )
    }
  }

  replaceLines (startRow, endRow, lines) {
    const lastRow = this.getLastRow()
    if (endRow > lastRow) {
      const lastLine = this.getLine(lastRow)
      this.doc.replaceRange(
        lines.join('\n'),
        { line: startRow, ch: 0 },
        { line: lastRow, ch: lastLine.length }
      )
    } else {
      this.doc.replaceRange(
        lines.join('\n') + '\n',
        { line: startRow, ch: 0 },
        { line: endRow, ch: 0 }
      )
    }
  }

  transact (func) {
    this.transaction = true
    func()
    this.transaction = false
    if (this.onDidFinishTransaction) {
      this.onDidFinishTransaction.call(undefined)
    }
  }
}

function initTableEditor (editor) {
  // create an interface to the text editor
  const editorIntf = new TextEditorInterface(editor)
  // create a table editor object
  const tableEditor = new __WEBPACK_IMPORTED_MODULE_0__susisu_mte_kernel__["d" /* TableEditor */](editorIntf)
  // options for the table editor
  const opts = Object(__WEBPACK_IMPORTED_MODULE_0__susisu_mte_kernel__["e" /* options */])({
    smartCursor: true,
    formatType: __WEBPACK_IMPORTED_MODULE_0__susisu_mte_kernel__["b" /* FormatType */].NORMAL
  })
  // keymap of the commands
  // from https://github.com/susisu/mte-demo/blob/master/src/main.js
  const keyMap = CodeMirror.normalizeKeyMap({
    Tab: () => { tableEditor.nextCell(opts) },
    'Shift-Tab': () => { tableEditor.previousCell(opts) },
    Enter: () => { tableEditor.nextRow(opts) },
    'Ctrl-Enter': () => { tableEditor.escape(opts) },
    'Cmd-Enter': () => { tableEditor.escape(opts) },
    'Shift-Ctrl-Left': () => { tableEditor.alignColumn(__WEBPACK_IMPORTED_MODULE_0__susisu_mte_kernel__["a" /* Alignment */].LEFT, opts) },
    'Shift-Cmd-Left': () => { tableEditor.alignColumn(__WEBPACK_IMPORTED_MODULE_0__susisu_mte_kernel__["a" /* Alignment */].LEFT, opts) },
    'Shift-Ctrl-Right': () => { tableEditor.alignColumn(__WEBPACK_IMPORTED_MODULE_0__susisu_mte_kernel__["a" /* Alignment */].RIGHT, opts) },
    'Shift-Cmd-Right': () => { tableEditor.alignColumn(__WEBPACK_IMPORTED_MODULE_0__susisu_mte_kernel__["a" /* Alignment */].RIGHT, opts) },
    'Shift-Ctrl-Up': () => { tableEditor.alignColumn(__WEBPACK_IMPORTED_MODULE_0__susisu_mte_kernel__["a" /* Alignment */].CENTER, opts) },
    'Shift-Cmd-Up': () => { tableEditor.alignColumn(__WEBPACK_IMPORTED_MODULE_0__susisu_mte_kernel__["a" /* Alignment */].CENTER, opts) },
    'Shift-Ctrl-Down': () => { tableEditor.alignColumn(__WEBPACK_IMPORTED_MODULE_0__susisu_mte_kernel__["a" /* Alignment */].NONE, opts) },
    'Shift-Cmd-Down': () => { tableEditor.alignColumn(__WEBPACK_IMPORTED_MODULE_0__susisu_mte_kernel__["a" /* Alignment */].NONE, opts) },
    'Ctrl-Left': () => { tableEditor.moveFocus(0, -1, opts) },
    'Cmd-Left': () => { tableEditor.moveFocus(0, -1, opts) },
    'Ctrl-Right': () => { tableEditor.moveFocus(0, 1, opts) },
    'Cmd-Right': () => { tableEditor.moveFocus(0, 1, opts) },
    'Ctrl-Up': () => { tableEditor.moveFocus(-1, 0, opts) },
    'Cmd-Up': () => { tableEditor.moveFocus(-1, 0, opts) },
    'Ctrl-Down': () => { tableEditor.moveFocus(1, 0, opts) },
    'Cmd-Down': () => { tableEditor.moveFocus(1, 0, opts) },
    'Ctrl-K Ctrl-I': () => { tableEditor.insertRow(opts) },
    'Cmd-K Cmd-I': () => { tableEditor.insertRow(opts) },
    'Ctrl-L Ctrl-I': () => { tableEditor.deleteRow(opts) },
    'Cmd-L Cmd-I': () => { tableEditor.deleteRow(opts) },
    'Ctrl-K Ctrl-J': () => { tableEditor.insertColumn(opts) },
    'Cmd-K Cmd-J': () => { tableEditor.insertColumn(opts) },
    'Ctrl-L Ctrl-J': () => { tableEditor.deleteColumn(opts) },
    'Cmd-L Cmd-J': () => { tableEditor.deleteColumn(opts) },
    'Alt-Shift-Ctrl-Left': () => { tableEditor.moveColumn(-1, opts) },
    'Alt-Shift-Cmd-Left': () => { tableEditor.moveColumn(-1, opts) },
    'Alt-Shift-Ctrl-Right': () => { tableEditor.moveColumn(1, opts) },
    'Alt-Shift-Cmd-Right': () => { tableEditor.moveColumn(1, opts) },
    'Alt-Shift-Ctrl-Up': () => { tableEditor.moveRow(-1, opts) },
    'Alt-Shift-Cmd-Up': () => { tableEditor.moveRow(-1, opts) },
    'Alt-Shift-Ctrl-Down': () => { tableEditor.moveRow(1, opts) },
    'Alt-Shift-Cmd-Down': () => { tableEditor.moveRow(1, opts) }
  })

  // enable keymap if the cursor is in a table
  function updateActiveState() {
    const active = tableEditor.cursorIsInTable();
    if (active) {
      editor.setOption("extraKeys", keyMap);
    }
    else {
      editor.setOption("extraKeys", null);
      tableEditor.resetSmartCursor();
    }
  }
  // event subscriptions
  editor.on("cursorActivity", () => {
    if (!editorIntf.transaction) {
      updateActiveState();
    }
  });
  editor.on("changes", () => {
    if (!editorIntf.transaction) {
      updateActiveState();
    }
  });
  editorIntf.onDidFinishTransaction = () => {
    updateActiveState();
  };

  return tableEditor
}



/***/ }),
/* 1 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "c", function() { return Point; });
/* unused harmony export Range */
/* unused harmony export Focus */
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return Alignment; });
/* unused harmony export DefaultAlignment */
/* unused harmony export HeaderAlignment */
/* unused harmony export TableCell */
/* unused harmony export TableRow */
/* unused harmony export Table */
/* unused harmony export readTable */
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "b", function() { return FormatType; });
/* unused harmony export completeTable */
/* unused harmony export formatTable */
/* unused harmony export alterAlignment */
/* unused harmony export insertRow */
/* unused harmony export deleteRow */
/* unused harmony export moveRow */
/* unused harmony export insertColumn */
/* unused harmony export deleteColumn */
/* unused harmony export moveColumn */
/* unused harmony export Insert */
/* unused harmony export Delete */
/* unused harmony export applyEditScript */
/* unused harmony export shortestEditScript */
/* unused harmony export ITextEditor */
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "e", function() { return options; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "d", function() { return TableEditor; });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_meaw__ = __webpack_require__(2);


/**
 * A `Point` represents a point in the text editor.
 */
class Point {
  /**
   * Creates a new `Point` object.
   *
   * @param {number} row - Row of the point, starts from 0.
   * @param {number} column - Column of the point, starts from 0.
   */
  constructor(row, column) {
    /** @private */
    this._row = row;
    /** @private */
    this._column = column;
  }

  /**
   * Row of the point.
   *
   * @type {number}
   */
  get row() {
    return this._row;
  }

  /**
   * Column of the point.
   *
   * @type {number}
   */
  get column() {
    return this._column;
  }

  /**
   * Checks if the point is equal to another point.
   *
   * @param {Point} point - A point object.
   * @returns {boolean} `true` if two points are equal.
   */
  equals(point) {
    return this.row === point.row && this.column === point.column;
  }
}

/**
 * A `Range` object represents a range in the text editor.
 */
class Range {
  /**
   * Creates a new `Range` object.
   *
   * @param {Point} start - The start point of the range.
   * @param {Point} end - The end point of the range.
   */
  constructor(start, end) {
    /** @private */
    this._start = start;
    /** @private */
    this._end = end;
  }

  /**
   * The start point of the range.
   *
   * @type {Point}
   */
  get start() {
    return this._start;
  }

  /**
   * The end point of the range.
   *
   * @type {Point}
   */
  get end() {
    return this._end;
  }
}

/**
 * A `Focus` object represents which cell is focused in the table.
 *
 * Note that `row` and `column` properties specifiy a cell's position in the table, not the cursor's
 * position in the text editor as {@link Point} class.
 *
 * @private
 */
class Focus {
  /**
   * Creates a new `Focus` object.
   *
   * @param {number} row - Row of the focused cell.
   * @param {number} column - Column of the focused cell.
   * @param {number} offset - Raw offset in the cell.
   */
  constructor(row, column, offset) {
    /** @private */
    this._row = row;
    /** @private */
    this._column = column;
    /** @private */
    this._offset = offset;
  }

  /**
   * Row of the focused cell.
   *
   * @type {number}
   */
  get row() {
    return this._row;
  }

  /**
   * Column of the focused cell.
   *
   * @type {number}
   */
  get column() {
    return this._column;
  }

  /**
   * Raw offset in the cell.
   *
   * @type {number}
   */
  get offset() {
    return this._offset;
  }

  /**
   * Checks if two focuses point the same cell.
   * Offsets are ignored.
   *
   * @param {Focus} focus - A focus object.
   * @returns {boolean}
   */
  posEquals(focus) {
    return this.row === focus.row && this.column === focus.column;
  }

  /**
   * Creates a copy of the focus object by setting its row to the specified value.
   *
   * @param {number} row - Row of the focused cell.
   * @returns {Focus} A new focus object with the specified row.
   */
  setRow(row) {
    return new Focus(row, this.column, this.offset);
  }

  /**
   * Creates a copy of the focus object by setting its column to the specified value.
   *
   * @param {number} column - Column of the focused cell.
   * @returns {Focus} A new focus object with the specified column.
   */
  setColumn(column) {
    return new Focus(this.row, column, this.offset);
  }

  /**
   * Creates a copy of the focus object by setting its offset to the specified value.
   *
   * @param {number} offset - Offset in the focused cell.
   * @returns {Focus} A new focus object with the specified offset.
   */
  setOffset(offset) {
    return new Focus(this.row, this.column, offset);
  }
}

/**
 * Represents column alignment.
 *
 * - `Alignment.NONE` - Use default alignment.
 * - `Alignment.LEFT` - Align left.
 * - `Alignment.RIGHT` - Align right.
 * - `Alignment.CENTER` - Align center.
 *
 * @type {Object}
 */
const Alignment = Object.freeze({
  NONE  : "none",
  LEFT  : "left",
  RIGHT : "right",
  CENTER: "center"
});

/**
 * Represents default column alignment
 *
 * - `DefaultAlignment.LEFT` - Align left.
 * - `DefaultAlignment.RIGHT` - Align right.
 * - `DefaultAlignment.CENTER` - Align center.
 *
 * @type {Object}
 */
const DefaultAlignment = Object.freeze({
  LEFT  : Alignment.LEFT,
  RIGHT : Alignment.RIGHT,
  CENTER: Alignment.CENTER
});

/**
 * Represents alignment of header cells.
 *
 * - `HeaderAlignment.FOLLOW` - Follow column's alignment.
 * - `HeaderAlignment.LEFT` - Align left.
 * - `HeaderAlignment.RIGHT` - Align right.
 * - `HeaderAlignment.CENTER` - Align center.
 *
 * @type {Object}
 */
const HeaderAlignment = Object.freeze({
  FOLLOW: "follow",
  LEFT  : Alignment.LEFT,
  RIGHT : Alignment.RIGHT,
  CENTER: Alignment.CENTER
});

/**
 * A `TableCell` object represents a table cell.
 *
 * @private
 */
class TableCell {
  /**
   * Creates a new `TableCell` object.
   *
   * @param {string} rawContent - Raw content of the cell.
   */
  constructor(rawContent) {
    /** @private */
    this._rawContent = rawContent;
    /** @private */
    this._content = rawContent.trim();
    /** @private */
    this._paddingLeft = this._content === ""
      ? (this._rawContent === "" ? 0 : 1)
      : this._rawContent.length - this._rawContent.trimLeft().length;
    /** @private */
    this._paddingRight = this._rawContent.length - this._content.length - this._paddingLeft;
  }

  /**
   * Raw content of the cell.
   *
   * @type {string}
   */
  get rawContent() {
    return this._rawContent;
  }

  /**
   * Trimmed content of the cell.
   *
   * @type {string}
   */
  get content() {
    return this._content;
  }

  /**
   * Width of the left padding of the cell.
   *
   * @type {number}
   */
  get paddingLeft() {
    return this._paddingLeft;
  }

  /**
   * Width of the right padding of the cell.
   *
   * @type {number}
   */
  get paddingRight() {
    return this._paddingRight;
  }

  /**
   * Convers the cell to a text representation.
   *
   * @returns {string} The raw content of the cell.
   */
  toText() {
    return this.rawContent;
  }

  /**
   * Checks if the cell is a delimiter i.e. it only contains hyphens `-` with optional one
   * leading and trailing colons `:`.
   *
   * @returns {boolean} `true` if the cell is a delimiter.
   */
  isDelimiter() {
    return /^\s*:?-+:?\s*$/.test(this.rawContent);
  }

  /**
   * Returns the alignment the cell represents.
   *
   * @returns {Alignment|undefined} The alignment the cell represents;
   * `undefined` if the cell is not a delimiter.
   */
  getAlignment() {
    if (!this.isDelimiter()) {
      return undefined;
    }
    if (this.content[0] === ":") {
      if (this.content[this.content.length - 1] === ":") {
        return Alignment.CENTER;
      }
      else {
        return Alignment.LEFT;
      }
    }
    else {
      if (this.content[this.content.length - 1] === ":") {
        return Alignment.RIGHT;
      }
      else {
        return Alignment.NONE;
      }
    }
  }

  /**
   * Computes a relative position in the trimmed content from that in the raw content.
   *
   * @param {number} rawOffset - Relative position in the raw content.
   * @returns {number} - Relative position in the trimmed content.
   */
  computeContentOffset(rawOffset) {
    if (this.content === "") {
      return 0;
    }
    if (rawOffset < this.paddingLeft) {
      return 0;
    }
    if (rawOffset < this.paddingLeft + this.content.length) {
      return rawOffset - this.paddingLeft;
    }
    else {
      return this.content.length;
    }
  }

  /**
   * Computes a relative position in the raw content from that in the trimmed content.
   *
   * @param {number} contentOffset - Relative position in the trimmed content.
   * @returns {number} - Relative position in the raw content.
   */
  computeRawOffset(contentOffset) {
    return contentOffset + this.paddingLeft;
  }
}

/**
 * A `TableRow` object represents a table row.
 *
 * @private
 */
class TableRow {
  /**
   * Creates a new `TableRow` objec.
   *
   * @param {Array<TableCell>} cells - Cells that the row contains.
   * @param {string} marginLeft - Margin string at the left of the row.
   * @param {string} marginRight - Margin string at the right of the row.
   */
  constructor(cells, marginLeft, marginRight) {
    /** @private */
    this._cells = cells.slice();
    /** @private */
    this._marginLeft = marginLeft;
    /** @private */
    this._marginRight = marginRight;
  }

  /**
   * Margin string at the left of the row.
   *
   * @type {string}
   */
  get marginLeft() {
    return this._marginLeft;
  }

  /**
   * Margin string at the right of the row.
   *
   * @type {string}
   */
  get marginRight() {
    return this._marginRight;
  }

  /**
   * Gets the number of the cells in the row.
   *
   * @returns {number} Number of the cells.
   */
  getWidth() {
    return this._cells.length;
  }

  /**
   * Returns the cells that the row contains.
   *
   * @returns {Array<TableCell>} An array of cells that the row contains.
   */
  getCells() {
    return this._cells.slice();
  }

  /**
   * Gets a cell at the specified index.
   *
   * @param {number} index - Index.
   * @returns {TableCell|undefined} The cell at the specified index if exists;
   * `undefined` if no cell is found.
   */
  getCellAt(index) {
    return this._cells[index];
  }

  /**
   * Convers the row to a text representation.
   *
   * @returns {string} A text representation of the row.
   */
  toText() {
    if (this._cells.length === 0) {
      return this.marginLeft;
    }
    else {
      const cells = this._cells.map(cell => cell.toText()).join("|");
      return `${this.marginLeft}|${cells}|${this.marginRight}`;
    }
  }

  /**
   * Checks if the row is a delimiter or not.
   *
   * @returns {boolean} `true` if the row is a delimiter i.e. all the cells contained are delimiters.
   */
  isDelimiter() {
    return this._cells.every(cell => cell.isDelimiter());
  }
}

/**
 * A `Table` object represents a table.
 *
 * @private
 */
class Table {
  /**
   * Creates a new `Table` object.
   *
   * @param {Array<TableRow>} rows - An array of rows that the table contains.
   */
  constructor(rows) {
    /** @private */
    this._rows = rows.slice();
  }

  /**
   * Gets the number of rows in the table.
   *
   * @returns {number} The number of rows.
   */
  getHeight() {
    return this._rows.length;
  }

  /**
   * Gets the maximum width of the rows in the table.
   *
   * @returns {number} The maximum width of the rows.
   */
  getWidth() {
    return this._rows.map(row => row.getWidth())
      .reduce((x, y) => Math.max(x, y), 0);
  }

  /**
   * Gets the width of the header row.
   *
   * @returns {number|undefined} The width of the header row;
   * `undefined` if there is no header row.
   */
  getHeaderWidth() {
    if (this._rows.length === 0) {
      return undefined;
    }
    return this._rows[0].getWidth();
  }

  /**
   * Gets the rows that the table contains.
   *
   * @returns {Array<TableRow>} An array of the rows.
   */
  getRows() {
    return this._rows.slice();
  }

  /**
   * Gets a row at the specified index.
   *
   * @param {number} index - Row index.
   * @returns {TableRow|undefined} The row at the specified index;
   * `undefined` if not found.
   */
  getRowAt(index) {
    return this._rows[index];
  }

  /**
   * Gets the delimiter row of the table.
   *
   * @returns {TableRow|undefined} The delimiter row;
   * `undefined` if there is not delimiter row.
   */
  getDelimiterRow() {
    const row = this._rows[1];
    if (row === undefined) {
      return undefined;
    }
    if (row.isDelimiter()) {
      return row;
    }
    else {
      return undefined;
    }
  }

  /**
   * Gets a cell at the specified index.
   *
   * @param {number} rowIndex - Row index of the cell.
   * @param {number} columnIndex - Column index of the cell.
   * @returns {TableCell|undefined} The cell at the specified index;
   * `undefined` if not found.
   */
  getCellAt(rowIndex, columnIndex) {
    const row = this._rows[rowIndex];
    if (row === undefined) {
      return undefined;
    }
    return row.getCellAt(columnIndex);
  }

  /**
   * Gets the cell at the focus.
   *
   * @param {Focus} focus - Focus object.
   * @returns {TableCell|undefined} The cell at the focus;
   * `undefined` if not found.
   */
  getFocusedCell(focus) {
    return this.getCellAt(focus.row, focus.column);
  }

  /**
   * Converts the table to an array of text representations of the rows.
   *
   * @returns {Array<string>} An array of text representations of the rows.
   */
  toLines() {
    return this._rows.map(row => row.toText());
  }

  /**
   * Computes a focus from a point in the text editor.
   *
   * @param {Point} pos - A point in the text editor.
   * @param {number} rowOffset - The row index where the table starts in the text editor.
   * @returns {Focus|undefined} A focus object that corresponds to the specified point;
   * `undefined` if the row index is out of bounds.
   */
  focusOfPosition(pos, rowOffset) {
    const rowIndex = pos.row - rowOffset;
    const row = this._rows[rowIndex];
    if (row === undefined) {
      return undefined;
    }
    if (pos.column < row.marginLeft.length + 1) {
      return new Focus(rowIndex, -1, pos.column);
    }
    else {
      const cellWidths = row.getCells().map(cell => cell.rawContent.length);
      let columnPos = row.marginLeft.length + 1; // left margin + a pipe
      let columnIndex = 0;
      for (; columnIndex < cellWidths.length; columnIndex++) {
        if (columnPos + cellWidths[columnIndex] + 1 > pos.column) {
          break;
        }
        columnPos += cellWidths[columnIndex] + 1;
      }
      const offset = pos.column - columnPos;
      return new Focus(rowIndex, columnIndex, offset);
    }
  }

  /**
   * Computes a position in the text editor from a focus.
   *
   * @param {Focus} focus - A focus object.
   * @param {number} rowOffset - The row index where the table starts in the text editor.
   * @returns {Point|undefined} A position in the text editor that corresponds to the focus;
   * `undefined` if the focused row  is out of the table.
   */
  positionOfFocus(focus, rowOffset) {
    const row = this._rows[focus.row];
    if (row === undefined) {
      return undefined;
    }
    const rowPos = focus.row + rowOffset;
    if (focus.column < 0) {
      return new Point(rowPos, focus.offset);
    }
    const cellWidths = row.getCells().map(cell => cell.rawContent.length);
    const maxIndex = Math.min(focus.column, cellWidths.length);
    let columnPos = row.marginLeft.length + 1;
    for (let columnIndex = 0; columnIndex < maxIndex; columnIndex++) {
      columnPos += cellWidths[columnIndex] + 1;
    }
    return new Point(rowPos, columnPos + focus.offset);
  }

  /**
   * Computes a selection range from a focus.
   *
   * @param {Focus} focus - A focus object.
   * @param {number} rowOffset - The row index where the table starts in the text editor.
   * @returns {Range|undefined} A range to be selected that corresponds to the focus;
   * `undefined` if the focus does not specify any cell or the specified cell is empty.
   */
  selectionRangeOfFocus(focus, rowOffset) {
    const row = this._rows[focus.row];
    if (row === undefined) {
      return undefined;
    }
    const cell = row.getCellAt(focus.column);
    if (cell === undefined) {
      return undefined;
    }
    if (cell.content === "") {
      return undefined;
    }
    const rowPos = focus.row + rowOffset;
    const cellWidths = row.getCells().map(cell => cell.rawContent.length);
    let columnPos = row.marginLeft.length + 1;
    for (let columnIndex = 0; columnIndex < focus.column; columnIndex++) {
      columnPos += cellWidths[columnIndex] + 1;
    }
    columnPos += cell.paddingLeft;
    return new Range(
      new Point(rowPos, columnPos),
      new Point(rowPos, columnPos + cell.content.length)
    );
  }
}

/**
 * Splits a text into cells.
 *
 * @private
 * @param {string} text
 * @returns {Array<string>}
 */
function _splitCells(text) {
  const cells = [];
  let buf = "";
  let rest = text;
  while (rest !== "") {
    switch (rest[0]) {
    case "`":
      // read code span
      {
        const start = rest.match(/^`*/)[0];
        let buf1 = start;
        let rest1 = rest.substr(start.length);
        let closed = false;
        while (rest1 !== "") {
          if (rest1[0] === "`") {
            const end = rest1.match(/^`*/)[0];
            buf1 += end;
            rest1 = rest1.substr(end.length);
            if (end.length === start.length) {
              closed = true;
              break;
            }
          }
          else {
            buf1 += rest1[0];
            rest1 = rest1.substr(1);
          }
        }
        if (closed) {
          buf += buf1;
          rest = rest1;
        }
        else {
          buf += "`";
          rest = rest.substr(1);
        }
      }
      break;
    case "\\":
      // escape next character
      if (rest.length >= 2) {
        buf += rest.substr(0, 2);
        rest = rest.substr(2);
      }
      else {
        buf += "\\";
        rest = rest.substr(1);
      }
      break;
    case "|":
      // flush buffer
      cells.push(buf);
      buf = "";
      rest = rest.substr(1);
      break;
    default:
      buf += rest[0];
      rest = rest.substr(1);
    }
  }
  cells.push(buf);
  return cells;
}

/**
 * Reads a table row.
 *
 * @private
 * @param {string} text - A text
 * @returns {TableRow}
 */
function _readRow(text) {
  let cells = _splitCells(text);
  let marginLeft;
  if (cells.length > 0 && /^\s*$/.test(cells[0])) {
    marginLeft = cells[0];
    cells = cells.slice(1);
  }
  else {
    marginLeft = "";
  }
  let marginRight;
  if (cells.length > 1 && /^\s*$/.test(cells[cells.length - 1])) {
    marginRight = cells[cells.length - 1];
    cells = cells.slice(0, cells.length - 1);
  }
  else {
    marginRight = "";
  }
  return new TableRow(cells.map(cell => new TableCell(cell)), marginLeft, marginRight);
}

/**
 * Reads a table from lines.
 *
 * @private
 * @param {Array<string>} lines - An array of texts, each text represents a row.
 * @returns {Table} The table red from the lines.
 */
function readTable(lines) {
  return new Table(lines.map(_readRow));
}

/**
 * Creates a delimiter text.
 *
 * @private
 * @param {Alignment} alignment
 * @param {number} width - Width of the horizontal bar of delimiter.
 * @returns {string}
 * @throws {Error} Unknown alignment.
 */
function _delimiterText(alignment, width) {
  const bar = "-".repeat(width);
  switch (alignment) {
  case Alignment.NONE:
    return ` ${bar} `;
  case Alignment.LEFT:
    return `:${bar} `;
  case Alignment.RIGHT:
    return ` ${bar}:`;
  case Alignment.CENTER:
    return `:${bar}:`;
  default:
    throw new Error("Unknown alignment: " + alignment);
  }
}

/**
 * Extends array size.
 *
 * @private
 * @param {Array} arr
 * @param {number} size
 * @param {Function} callback - Callback function to fill newly created cells.
 * @returns {Array} Extended array.
 */
function _extendArray(arr, size, callback) {
  const extended = arr.slice();
  for (let i = arr.length; i < size; i++) {
    extended.push(callback(i, arr));
  }
  return extended;
}

/**
 * Completes a table by adding missing delimiter and cells.
 * After completion, all rows in the table have the same width.
 *
 * @private
 * @param {Table} table - A table object.
 * @param {Object} options - An object containing options for completion.
 *
 * | property name       | type           | description                                               |
 * | ------------------- | -------------- | --------------------------------------------------------- |
 * | `minDelimiterWidth` | {@link number} | Width of delimiters used when completing delimiter cells. |
 *
 * @returns {Object} An object that represents the result of the completion.
 *
 * | property name       | type            | description                            |
 * | ------------------- | --------------- | -------------------------------------- |
 * | `table`             | {@link Table}   | A completed table object.              |
 * | `delimiterInserted` | {@link boolean} | `true` if a delimiter row is inserted. |
 *
 * @throws {Error} Empty table.
 */
function completeTable(table, options) {
  const tableHeight = table.getHeight();
  const tableWidth = table.getWidth();
  if (tableHeight === 0) {
    throw new Error("Empty table");
  }
  const rows = table.getRows();
  const newRows = [];
  // header
  const headerRow = rows[0];
  const headerCells = headerRow.getCells();
  newRows.push(new TableRow(
    _extendArray(headerCells, tableWidth, j => new TableCell(
      j === headerCells.length ? headerRow.marginRight : ""
    )),
    headerRow.marginLeft,
    headerCells.length < tableWidth ? "" : headerRow.marginRight
  ));
  // delimiter
  const delimiterRow = table.getDelimiterRow();
  if (delimiterRow !== undefined) {
    const delimiterCells = delimiterRow.getCells();
    newRows.push(new TableRow(
      _extendArray(delimiterCells, tableWidth, j => new TableCell(
        _delimiterText(
          Alignment.NONE,
          j === delimiterCells.length
            ? Math.max(options.minDelimiterWidth, delimiterRow.marginRight.length - 2)
            : options.minDelimiterWidth
        )
      )),
      delimiterRow.marginLeft,
      delimiterCells.length < tableWidth ? "" : delimiterRow.marginRight
    ));
  }
  else {
    newRows.push(new TableRow(
      _extendArray([], tableWidth, () => new TableCell(
        _delimiterText(Alignment.NONE, options.minDelimiterWidth)
      )),
      "",
      ""
    ));
  }
  // body
  for (let i = delimiterRow !== undefined ? 2 : 1; i < tableHeight; i++) {
    const row = rows[i];
    const cells = row.getCells();
    newRows.push(new TableRow(
      _extendArray(cells, tableWidth, j => new TableCell(
        j === cells.length ? row.marginRight : ""
      )),
      row.marginLeft,
      cells.length < tableWidth ? "" : row.marginRight
    ));
  }
  return {
    table            : new Table(newRows),
    delimiterInserted: delimiterRow === undefined
  };
}

/**
 * Calculates the width of a text based on characters' EAW properties.
 *
 * @private
 * @param {string} text
 * @param {Object} options -
 *
 * | property name     | type                               |
 * | ----------------- | ---------------------------------- |
 * | `normalize`       | {@link boolean}                    |
 * | `wideChars`       | {@link Set}&lt;{@link string} &gt; |
 * | `narrowChars`     | {@link Set}&lt;{@link string} &gt; |
 * | `ambiguousAsWide` | {@link boolean}                    |
 *
 * @returns {number} Calculated width of the text.
 */
function _computeTextWidth(text, options) {
  const normalized = options.normalize ? text.normalize("NFC") : text;
  let w = 0;
  for (const char of normalized) {
    if (options.wideChars.has(char)) {
      w += 2;
      continue;
    }
    if (options.narrowChars.has(char)) {
      w += 1;
      continue;
    }
    switch (Object(__WEBPACK_IMPORTED_MODULE_0_meaw__["a" /* getEAW */])(char)) {
    case "F":
    case "W":
      w += 2;
      break;
    case "A":
      w += options.ambiguousAsWide ? 2 : 1;
      break;
    default:
      w += 1;
    }
  }
  return w;
}

/**
 * Returns a aligned cell content.
 *
 * @private
 * @param {string} text
 * @param {number} width
 * @param {Alignment} alignment
 * @param {Object} options - Options for computing text width.
 * @returns {string}
 * @throws {Error} Unknown alignment.
 * @throws {Error} Unexpected default alignment.
 */
function _alignText(text, width, alignment, options) {
  const space = width - _computeTextWidth(text, options);
  if (space < 0) {
    return text;
  }
  switch (alignment) {
  case Alignment.NONE:
    throw new Error("Unexpected default alignment");
  case Alignment.LEFT:
    return text + " ".repeat(space);
  case Alignment.RIGHT:
    return " ".repeat(space) + text;
  case Alignment.CENTER:
    return " ".repeat(Math.floor(space / 2))
      + text
      + " ".repeat(Math.ceil(space / 2));
  default:
    throw new Error("Unknown alignment: " + alignment);
  }
}

/**
 * Just adds one space paddings to both sides of a text.
 *
 * @private
 * @param {string} text
 * @returns {string}
 */
function _padText(text) {
  return ` ${text} `;
}

/**
 * Formats a table.
 *
 * @private
 * @param {Table} table - A table object.
 * @param {Object} options - An object containing options for formatting.
 *
 * | property name       | type                     | description                                             |
 * | ------------------- | ------------------------ | ------------------------------------------------------- |
 * | `minDelimiterWidth` | {@link number}           | Minimum width of delimiters.                            |
 * | `defaultAlignment`  | {@link DefaultAlignment} | Default alignment of columns.                           |
 * | `headerAlignment`   | {@link HeaderAlignment}  | Alignment of header cells.                              |
 * | `textWidthOptions`  | {@link Object}           | An object containing options for computing text widths. |
 *
 * `options.textWidthOptions` must contain the following options.
 *
 * | property name     | type                              | description                                         |
 * | ----------------- | --------------------------------- | --------------------------------------------------- |
 * | `normalize`       | {@link boolean}                   | Normalize texts before computing text widths.       |
 * | `wideChars`       | {@link Set}&lt;{@link string}&gt; | Set of characters that should be treated as wide.   |
 * | `narrowChars`     | {@link Set}&lt;{@link string}&gt; | Set of characters that should be treated as narrow. |
 * | `ambiguousAsWide` | {@link boolean}                   | Treat East Asian Ambiguous characters as wide.      |
 *
 * @returns {Object} An object that represents the result of formatting.
 *
 * | property name   | type           | description                                    |
 * | --------------- | -------------- | ---------------------------------------------- |
 * | `table`         | {@link Table}  | A formatted table object.                      |
 * | `marginLeft`    | {@link string} | The common left margin of the formatted table. |
 */
function _formatTable(table, options) {
  const tableHeight = table.getHeight();
  const tableWidth = table.getWidth();
  if (tableHeight === 0) {
    return {
      table,
      marginLeft: ""
    };
  }
  const marginLeft = table.getRowAt(0).marginLeft;
  if (tableWidth === 0) {
    const rows = new Array(tableHeight).fill()
      .map(() => new TableRow([], marginLeft, ""));
    return {
      table: new Table(rows),
      marginLeft
    };
  }
  // compute column widths
  const delimiterRow = table.getDelimiterRow();
  const columnWidths = new Array(tableWidth).fill(0);
  if (delimiterRow !== undefined) {
    const delimiterRowWidth = delimiterRow.getWidth();
    for (let j = 0; j < delimiterRowWidth; j++) {
      columnWidths[j] = options.minDelimiterWidth;
    }
  }
  for (let i = 0; i < tableHeight; i++) {
    if (delimiterRow !== undefined && i === 1) {
      continue;
    }
    const row = table.getRowAt(i);
    const rowWidth = row.getWidth();
    for (let j = 0; j < rowWidth; j++) {
      columnWidths[j] = Math.max(
        columnWidths[j],
        _computeTextWidth(row.getCellAt(j).content, options.textWidthOptions)
      );
    }
  }
  // get column alignments
  const alignments = delimiterRow !== undefined
    ? _extendArray(
      delimiterRow.getCells().map(cell => cell.getAlignment()),
      tableWidth,
      () => options.defaultAlignment
    )
    : new Array(tableWidth).fill(options.defaultAlignment);
  // format
  const rows = [];
  // header
  const headerRow = table.getRowAt(0);
  rows.push(new TableRow(
    headerRow.getCells().map((cell, j) =>
      new TableCell(_padText(_alignText(
        cell.content,
        columnWidths[j],
        options.headerAlignment === HeaderAlignment.FOLLOW
          ? (alignments[j] === Alignment.NONE ? options.defaultAlignment : alignments[j])
          : options.headerAlignment,
        options.textWidthOptions
      )))
    ),
    marginLeft,
    ""
  ));
  // delimiter
  if (delimiterRow !== undefined) {
    rows.push(new TableRow(
      delimiterRow.getCells().map((cell, j) =>
        new TableCell(_delimiterText(alignments[j], columnWidths[j]))
      ),
      marginLeft,
      ""
    ));
  }
  // body
  for (let i = delimiterRow !== undefined ? 2 : 1; i < tableHeight; i++) {
    const row = table.getRowAt(i);
    rows.push(new TableRow(
      row.getCells().map((cell, j) =>
        new TableCell(_padText(_alignText(
          cell.content,
          columnWidths[j],
          alignments[j] === Alignment.NONE ? options.defaultAlignment : alignments[j],
          options.textWidthOptions
        )))
      ),
      marginLeft,
      ""
    ));
  }
  return {
    table: new Table(rows),
    marginLeft
  };
}

/**
 * Formats a table weakly.
 * Rows are formatted independently to each other, cell contents are just trimmed and not aligned.
 * This is useful when using a non-monospaced font or dealing with wide tables.
 *
 * @private
 * @param {Table} table - A table object.
 * @param {Object} options - An object containing options for formatting.
 * The function accepts the same option object for {@link formatTable}, but properties not listed
 * here are just ignored.
 *
 * | property name       | type           | description          |
 * | ------------------- | -------------- | -------------------- |
 * | `minDelimiterWidth` | {@link number} | Width of delimiters. |
 *
 * @returns {Object} An object that represents the result of formatting.
 *
 * | property name   | type           | description                                    |
 * | --------------- | -------------- | ---------------------------------------------- |
 * | `table`         | {@link Table}  | A formatted table object.                      |
 * | `marginLeft`    | {@link string} | The common left margin of the formatted table. |
 */
function _weakFormatTable(table, options) {
  const tableHeight = table.getHeight();
  const tableWidth = table.getWidth();
  if (tableHeight === 0) {
    return {
      table,
      marginLeft: ""
    };
  }
  const marginLeft = table.getRowAt(0).marginLeft;
  if (tableWidth === 0) {
    const rows = new Array(tableHeight).fill()
      .map(() => new TableRow([], marginLeft, ""));
    return {
      table: new Table(rows),
      marginLeft
    };
  }
  const delimiterRow = table.getDelimiterRow();
  // format
  const rows = [];
  // header
  const headerRow = table.getRowAt(0);
  rows.push(new TableRow(
    headerRow.getCells().map(cell =>
      new TableCell(_padText(cell.content))
    ),
    marginLeft,
    ""
  ));
  // delimiter
  if (delimiterRow !== undefined) {
    rows.push(new TableRow(
      delimiterRow.getCells().map(cell =>
        new TableCell(_delimiterText(cell.getAlignment(), options.minDelimiterWidth))
      ),
      marginLeft,
      ""
    ));
  }
  // body
  for (let i = delimiterRow !== undefined ? 2 : 1; i < tableHeight; i++) {
    const row = table.getRowAt(i);
    rows.push(new TableRow(
      row.getCells().map(cell =>
        new TableCell(_padText(cell.content))
      ),
      marginLeft,
      ""
    ));
  }
  return {
    table: new Table(rows),
    marginLeft
  };
}

/**
 * Represents table format type.
 *
 * - `FormatType.NORMAL` - Formats table normally.
 * - `FormatType.WEAK` - Formats table weakly, rows are formatted independently to each other, cell
 *   contents are just trimmed and not aligned.
 *
 * @type {Object}
 */
const FormatType = Object.freeze({
  NORMAL: "normal",
  WEAK  : "weak"
});


/**
 * Formats a table.
 *
 * @private
 * @param {Table} table - A table object.
 * @param {Object} options - An object containing options for formatting.
 *
 * | property name       | type                     | description                                             |
 * | ------------------- | ------------------------ | ------------------------------------------------------- |
 * | `formatType`        | {@link FormatType}       | Format type, normal or weak.                            |
 * | `minDelimiterWidth` | {@link number}           | Minimum width of delimiters.                            |
 * | `defaultAlignment`  | {@link DefaultAlignment} | Default alignment of columns.                           |
 * | `headerAlignment`   | {@link HeaderAlignment}  | Alignment of header cells.                              |
 * | `textWidthOptions`  | {@link Object}           | An object containing options for computing text widths. |
 *
 * `options.textWidthOptions` must contain the following options.
 *
 * | property name     | type                              | description                                         |
 * | ----------------- | --------------------------------- | --------------------------------------------------- |
 * | `normalize`       | {@link boolean}                   | Normalize texts before computing text widths.       |
 * | `wideChars`       | {@link Set}&lt;{@link string}&gt; | Set of characters that should be treated as wide.   |
 * | `narrowChars`     | {@link Set}&lt;{@link string}&gt; | Set of characters that should be treated as narrow. |
 * | `ambiguousAsWide` | {@link boolean}                   | Treat East Asian Ambiguous characters as wide.      |
 *
 * @returns {Object} An object that represents the result of formatting.
 *
 * | property name   | type           | description                                    |
 * | --------------- | -------------- | ---------------------------------------------- |
 * | `table`         | {@link Table}  | A formatted table object.                      |
 * | `marginLeft`    | {@link string} | The common left margin of the formatted table. |
 *
 * @throws {Error} Unknown format type.
 */
function formatTable(table, options) {
  switch (options.formatType) {
  case FormatType.NORMAL:
    return _formatTable(table, options);
  case FormatType.WEAK:
    return _weakFormatTable(table, options);
  default:
    throw new Error("Unknown format type: " + options.formatType);
  }
}

/**
 * Alters a column's alignment of a table.
 *
 * @private
 * @param {Table} table - A completed non-empty table.
 * @param {number} columnIndex - An index of the column.
 * @param {Alignment} alignment - A new alignment of the column.
 * @param {Object} options - An object containing options for completion.
 *
 * | property name       | type           | description          |
 * | ------------------- | -------------- | -------------------- |
 * | `minDelimiterWidth` | {@link number} | Width of delimiters. |
 *
 * @returns {Table} An altered table object.
 * If the column index is out of range, returns the original table.
 */
function alterAlignment(table, columnIndex, alignment, options) {
  const delimiterRow = table.getRowAt(1);
  if (columnIndex < 0 || delimiterRow.getWidth() - 1 < columnIndex) {
    return table;
  }
  const delimiterCells = delimiterRow.getCells();
  delimiterCells[columnIndex] = new TableCell(_delimiterText(alignment, options.minDelimiterWidth));
  const rows = table.getRows();
  rows[1] = new TableRow(delimiterCells, delimiterRow.marginLeft, delimiterRow.marginRight);
  return new Table(rows);
}

/**
 * Inserts a row to a table.
 * The row is always inserted after the header and the delimiter rows, even if the index specifies
 * the header or the delimiter.
 *
 * @private
 * @param {Table} table - A completed non-empty table.
 * @param {number} rowIndex - An row index at which a new row will be inserted.
 * @param {TableRow} row - A table row to be inserted.
 * @returns {Table} An altered table obejct.
 */
function insertRow(table, rowIndex, row) {
  const rows = table.getRows();
  rows.splice(Math.max(rowIndex, 2), 0, row);
  return new Table(rows);
}

/**
 * Deletes a row in a table.
 * If the index specifies the header row, the cells are emptied but the row will not be removed.
 * If the index specifies the delimiter row, it does nothing.
 *
 * @private
 * @param {Table} table - A completed non-empty table.
 * @param {number} rowIndex - An index of the row to be deleted.
 * @returns {Table} An altered table obejct.
 */
function deleteRow(table, rowIndex) {
  if (rowIndex === 1) {
    return table;
  }
  const rows = table.getRows();
  if (rowIndex === 0) {
    const headerRow = rows[0];
    rows[0] = new TableRow(
      new Array(headerRow.getWidth()).fill(new TableCell("")),
      headerRow.marginLeft,
      headerRow.marginRight
    );
  }
  else {
    rows.splice(rowIndex, 1);
  }
  return new Table(rows);
}

/**
 * Moves a row at the index to the specified destination.
 *
 * @private
 * @param {Table} table - A completed non-empty table.
 * @param {number} rowIndex - Index of the row to be moved.
 * @param {number} destIndex - Index of the destination.
 * @returns {Table} An altered table object.
 */
function moveRow(table, rowIndex, destIndex) {
  if (rowIndex <= 1 || destIndex <= 1 || rowIndex === destIndex) {
    return table;
  }
  const rows = table.getRows();
  const row = rows[rowIndex];
  rows.splice(rowIndex, 1);
  rows.splice(destIndex, 0, row);
  return new Table(rows);
}

/**
 * Inserts a column to a table.
 *
 * @private
 * @param {Table} table - A completed non-empty table.
 * @param {number} columnIndex - An column index at which the new column will be inserted.
 * @param {Array<TableCell>} column - An array of cells.
 * @param {Object} options - An object containing options for completion.
 *
 * | property name       | type           | description             |
 * | ------------------- | -------------- | ----------------------- |
 * | `minDelimiterWidth` | {@link number} | Width of the delimiter. |
 *
 * @returns {Table} An altered table obejct.
 */
function insertColumn(table, columnIndex, column, options) {
  const rows = table.getRows();
  for (let i = 0; i < rows.length; i++) {
    const row = rows[i];
    const cells = rows[i].getCells();
    const cell = i === 1
      ? new TableCell(_delimiterText(Alignment.NONE, options.minDelimiterWidth))
      : column[i > 1 ? i - 1 : i];
    cells.splice(columnIndex, 0, cell);
    rows[i] = new TableRow(cells, row.marginLeft, row.marginRight);
  }
  return new Table(rows);
}

/**
 * Deletes a column in a table.
 * If there will be no columns after the deletion, the cells are emptied but the column will not be
 * removed.
 *
 * @private
 * @param {Table} table - A completed non-empty table.
 * @param {number} columnIndex - An index of the column to be deleted.
 * @param {Object} options - An object containing options for completion.
 *
 * | property name       | type           | description             |
 * | ------------------- | -------------- | ----------------------- |
 * | `minDelimiterWidth` | {@link number} | Width of the delimiter. |
 *
 * @returns {Table} An altered table object.
 */
function deleteColumn(table, columnIndex, options) {
  const rows = table.getRows();
  for (let i = 0; i < rows.length; i++) {
    const row = rows[i];
    let cells = row.getCells();
    if (cells.length <= 1) {
      cells = [new TableCell(i === 1
        ? _delimiterText(Alignment.NONE, options.minDelimiterWidth)
        : ""
      )];
    }
    else {
      cells.splice(columnIndex, 1);
    }
    rows[i] = new TableRow(cells, row.marginLeft, row.marginRight);
  }
  return new Table(rows);
}

/**
 * Moves a column at the index to the specified destination.
 *
 * @private
 * @param {Table} table - A completed non-empty table.
 * @param {number} columnIndex - Index of the column to be moved.
 * @param {number} destIndex - Index of the destination.
 * @returns {Table} An altered table object.
 */
function moveColumn(table, columnIndex, destIndex) {
  if (columnIndex === destIndex) {
    return table;
  }
  const rows = table.getRows();
  for (let i = 0; i < rows.length; i++) {
    const row = rows[i];
    const cells = row.getCells();
    const cell = cells[columnIndex];
    cells.splice(columnIndex, 1);
    cells.splice(destIndex, 0, cell);
    rows[i] = new TableRow(cells, row.marginLeft, row.marginRight);
  }
  return new Table(rows);
}

/**
 * The `Insert` class represents an insertion of a line.
 *
 * @private
 */
class Insert {
  /**
   * Creats a new `Insert` object.
   *
   * @param {number} row - Row index, starts from `0`.
   * @param {string} line - A string to be inserted at the row.
   */
  constructor(row, line) {
    /** @private */
    this._row = row;
    /** @private */
    this._line = line;
  }

  /**
   * Row index, starts from `0`.
   *
   * @type {number}
   */
  get row() {
    return this._row;
  }

  /**
   * A string to be inserted.
   *
   * @type {string}
   */
  get line() {
    return this._line;
  }
}

/**
 * The `Delete` class represents a deletion of a line.
 *
 * @private
 */
class Delete {
  /**
   * Creates a new `Delete` object.
   *
   * @param {number} row - Row index, starts from `0`.
   */
  constructor(row) {
    /** @private */
    this._row = row;
  }

  /**
   * Row index, starts from `0`.
   *
   * @type {number}
   */
  get row() {
    return this._row;
  }
}

/**
 * Applies a command to the text editor.
 *
 * @private
 * @param {ITextEditor} textEditor - An interface to the text editor.
 * @param {Insert|Delete} command - A command.
 * @param {number} rowOffset - Offset to the row index of the command.
 * @returns {undefined}
 */
function _applyCommand(textEditor, command, rowOffset) {
  if (command instanceof Insert) {
    textEditor.insertLine(rowOffset + command.row, command.line);
  }
  else if (command instanceof Delete) {
    textEditor.deleteLine(rowOffset + command.row);
  }
  else {
    throw new Error("Unknown command");
  }
}

/**
 * Apply an edit script (array of commands) to the text editor.
 *
 * @private
 * @param {ITextEditor} textEditor - An interface to the text editor.
 * @param {Array<Insert|Delete>} script - An array of commands.
 * The commands are applied sequentially in the order of the array.
 * @param {number} rowOffset - Offset to the row index of the commands.
 * @returns {undefined}
 */
function applyEditScript(textEditor, script, rowOffset) {
  for (const command of script) {
    _applyCommand(textEditor, command, rowOffset);
  }
}


/**
 * Linked list used to remember edit script.
 *
 * @private
 */
class IList {
  get car() {
    throw new Error("Not implemented");
  }

  get cdr() {
    throw new Error("Not implemented");
  }

  isEmpty() {
    throw new Error("Not implemented");
  }

  unshift(value) {
    return new Cons(value, this);
  }

  toArray() {
    const arr = [];
    let rest = this;
    while (!rest.isEmpty()) {
      arr.push(rest.car);
      rest = rest.cdr;
    }
    return arr;
  }
}

/**
 * @private
 */
class Nil extends IList {
  constructor() {
    super();
  }

  get car() {
    throw new Error("Empty list");
  }

  get cdr() {
    throw new Error("Empty list");
  }

  isEmpty() {
    return true;
  }
}

/**
 * @private
 */
class Cons extends IList {
  constructor(car, cdr) {
    super();
    this._car = car;
    this._cdr = cdr;
  }

  get car() {
    return this._car;
  }

  get cdr() {
    return this._cdr;
  }

  isEmpty() {
    return false;
  }
}

const nil = new Nil();


/**
 * Computes the shortest edit script between two arrays of strings.
 *
 * @private
 * @param {Array<string>} from - An array of string the edit starts from.
 * @param {Array<string>} to - An array of string the edit goes to.
 * @param {number} [limit=-1] - Upper limit of edit distance to be searched.
 * If negative, there is no limit.
 * @returns {Array<Insert|Delete>|undefined} The shortest edit script that turns `from` into `to`;
 * `undefined` if no edit script is found in the given range.
 */
function shortestEditScript(from, to, limit = -1) {
  const fromLen = from.length;
  const toLen = to.length;
  const maxd = limit >= 0 ? Math.min(limit, fromLen + toLen) : fromLen + toLen;
  const mem = new Array(Math.min(maxd, fromLen) + Math.min(maxd, toLen) + 1);
  const offset = Math.min(maxd, fromLen);
  for (let d = 0; d <= maxd; d++) {
    const mink = d <= fromLen ? -d :  d - 2 * fromLen;
    const maxk = d <= toLen   ?  d : -d + 2 * toLen;
    for (let k = mink; k <= maxk; k += 2) {
      let i;
      let script;
      if (d === 0) {
        i = 0;
        script = nil;
      }
      else if (k === -d) {
        i = mem[offset + k + 1].i + 1;
        script = mem[offset + k + 1].script.unshift(new Delete(i + k));
      }
      else if (k === d) {
        i = mem[offset + k - 1].i;
        script = mem[offset + k - 1].script.unshift(new Insert(i + k - 1, to[i + k - 1]));
      }
      else {
        const vi = mem[offset + k + 1].i + 1;
        const hi = mem[offset + k - 1].i;
        if (vi > hi) {
          i = vi;
          script = mem[offset + k + 1].script.unshift(new Delete(i + k));
        }
        else {
          i = hi;
          script = mem[offset + k - 1].script.unshift(new Insert(i + k - 1, to[i + k - 1]));
        }
      }
      while (i < fromLen && i + k < toLen && from[i] === to[i + k]) {
        i += 1;
      }
      if (k === toLen - fromLen && i === fromLen) {
        return script.toArray().reverse();
      }
      mem[offset + k] = { i, script };
    }
  }
  return undefined;
}

/**
 * The `ITextEditor` represents an interface to a text editor.
 *
 * @interface
 */
class ITextEditor {
  /**
   * Gets the current cursor position.
   *
   * @returns {Point} A point object that represents the cursor position.
   */
  getCursorPosition() {
    throw new Error("Not implemented: getCursorPosition");
  }

  /**
   * Sets the cursor position to a specified one.
   *
   * @param {Point} pos - A point object which the cursor position is set to.
   * @returns {undefined}
   */
  setCursorPosition(pos) {
    throw new Error("Not implemented: setCursorPosition");
  }

  /**
   * Sets the selection range.
   * This method also expects the cursor position to be moved as the end of the selection range.
   *
   * @param {Range} range - A range object that describes a selection range.
   * @returns {undefined}
   */
  setSelectionRange(range) {
    throw new Error("Not implemented: setSelectionRange");
  }

  /**
   * Gets the last row index of the text editor.
   *
   * @returns {number} The last row index.
   */
  getLastRow() {
    throw new Error("Not implemented: getLastRow");
  }

  /**
   * Checks if the editor accepts a table at a row to be editted.
   * It should return `false` if, for example, the row is in a code block (not Markdown).
   *
   * @param {number} row - A row index in the text editor.
   * @returns {boolean} `true` if the table at the row can be editted.
   */
  acceptsTableEdit(row) {
    throw new Error("Not implemented: acceptsTableEdit");
  }

  /**
   * Gets a line string at a row.
   *
   * @param {number} row - Row index, starts from `0`.
   * @returns {string} The line at the specified row.
   * The line must not contain an EOL like `"\n"` or `"\r"`.
   */
  getLine(row) {
    throw new Error("Not implemented: getLine");
  }

  /**
   * Inserts a line at a specified row.
   *
   * @param {number} row - Row index, starts from `0`.
   * @param {string} line - A string to be inserted.
   * This must not contain an EOL like `"\n"` or `"\r"`.
   * @return {undefined}
   */
  insertLine(row, line) {
    throw new Error("Not implemented: insertLine");
  }

  /**
   * Deletes a line at a specified row.
   *
   * @param {number} row - Row index, starts from `0`.
   * @returns {undefined}
   */
  deleteLine(row) {
    throw new Error("Not implemented: deleteLine");
  }

  /**
   * Replace lines in a specified range.
   *
   * @param {number} startRow - Start row index, starts from `0`.
   * @param {number} endRow - End row index.
   * Lines from `startRow` to `endRow - 1` is replaced.
   * @param {Array<string>} lines - An array of string.
   * Each strings must not contain an EOL like `"\n"` or `"\r"`.
   * @returns {undefined}
   */
  replaceLines(startRow, endRow, lines) {
    throw new Error("Not implemented: replaceLines");
  }

  /**
   * Batches multiple operations as a single undo/redo step.
   *
   * @param {Function} func - A callback function that executes some operations on the text editor.
   * @returns {undefined}
   */
  transact(func) {
    throw new Error("Not implemented: transact");
  }
}

/**
 * Reads a property of an object if exists; otherwise uses a default value.
 *
 * @private
 * @param {*} obj - An object. If a non-object value is specified, the default value is used.
 * @param {string} key - A key (or property name).
 * @param {*} defaultVal - A default value that is used when a value does not exist.
 * @returns {*} A read value or the default value.
 */
function _value(obj, key, defaultVal) {
  return (typeof obj === "object" && obj !== null && obj[key] !== undefined)
    ? obj[key]
    : defaultVal;
}

/**
 * Reads multiple properties of an object if exists; otherwise uses default values.
 *
 * @private
 * @param {*} obj - An object. If a non-object value is specified, the default value is used.
 * @param {Object} keys - An object that consists of pairs of a key and a default value.
 * @returns {Object} A new object that contains read values.
 */
function _values(obj, keys) {
  const res = {};
  for (const [key, defaultVal] of Object.entries(keys)) {
    res[key] = _value(obj, key, defaultVal);
  }
  return res;
}

/**
 * Reads options for the formatter from an object.
 * The default values are used for options that are not specified.
 *
 * @param {Object} obj - An object containing options.
 * The available options and default values are listed below.
 *
 * | property name       | type                     | description                                             | default value            |
 * | ------------------- | ------------------------ | ------------------------------------------------------- | ------------------------ |
 * | `formatType`        | {@link FormatType}       | Format type, normal or weak.                            | `FormatType.NORMAL`      |
 * | `minDelimiterWidth` | {@link number}           | Minimum width of delimiters.                            | `3`                      |
 * | `defaultAlignment`  | {@link DefaultAlignment} | Default alignment of columns.                           | `DefaultAlignment.LEFT`  |
 * | `headerAlignment`   | {@link HeaderAlignment}  | Alignment of header cells.                              | `HeaderAlignment.FOLLOW` |
 * | `textWidthOptions`  | {@link Object}           | An object containing options for computing text widths. |                          |
 * | `smartCursor`       | {@link boolean}          | Enables "Smart Cursor" feature.                         | `false`                  |
 *
 * The available options for `obj.textWidthOptions` are the following ones.
 *
 * | property name     | type                              | description                                           | default value |
 * | ----------------- | --------------------------------- | ----------------------------------------------------- | ------------- |
 * | `normalize`       | {@link boolean}                   | Normalizes texts before computing text widths.        | `true`        |
 * | `wideChars`       | {@link Set}&lt;{@link string}&gt; | A set of characters that should be treated as wide.   | `new Set()`   |
 * | `narrowChars`     | {@link Set}&lt;{@link string}&gt; | A set of characters that should be treated as narrow. | `new Set()`   |
 * | `ambiguousAsWide` | {@link boolean}                   | Treats East Asian Ambiguous characters as wide.       | `false`       |
 *
 * @returns {Object} - An object that contains complete options.
 */
function options(obj) {
  const res = _values(obj, {
    formatType       : FormatType.NORMAL,
    minDelimiterWidth: 3,
    defaultAlignment : DefaultAlignment.LEFT,
    headerAlignment  : HeaderAlignment.FOLLOW,
    smartCursor      : false
  });
  res.textWidthOptions = _values(obj.textWidthOptions, {
    normalize      : true,
    wideChars      : new Set(),
    narrowChars    : new Set(),
    ambiguousAsWide: false
  });
  return res;
}

/**
 * Checks if a line is a table row.
 *
 * @private
 * @param {string} line - A string.
 * @returns {boolean} `true` if the given line starts with a pipe `|`.
 */
function _isTableRow(line) {
  return line.trimLeft()[0] === "|";
}

/**
 * Computes new focus offset from information of completed and formatted tables.
 *
 * @private
 * @param {Focus} focus - A focus.
 * @param {Table} table - A completed but not formatted table with original cell contents.
 * @param {Object} formatted - Information of the formatted table.
 * @param {boolean} moved - Indicates whether the focus position is moved by a command or not.
 * @returns {number}
 */
function _computeNewOffset(focus, table, formatted, moved) {
  if (moved) {
    const formattedFocusedCell = formatted.table.getFocusedCell(focus);
    if (formattedFocusedCell !== undefined) {
      return formattedFocusedCell.computeRawOffset(0);
    }
    else {
      return focus.column < 0 ? formatted.marginLeft.length : 0;
    }
  }
  else {
    const focusedCell = table.getFocusedCell(focus);
    const formattedFocusedCell = formatted.table.getFocusedCell(focus);
    if (focusedCell !== undefined && formattedFocusedCell !== undefined) {
      const contentOffset = Math.min(
        focusedCell.computeContentOffset(focus.offset),
        formattedFocusedCell.content.length
      );
      return formattedFocusedCell.computeRawOffset(contentOffset);
    }
    else {
      return focus.column < 0 ? formatted.marginLeft.length : 0;
    }
  }
}

/**
 * The `TableEditor` class is at the center of the markdown-table-editor.
 * When a command is executed, it reads a table from the text editor, does some operation on the
 * table, and then apply the result to the text editor.
 *
 * To use this class, the text editor (or an interface to it) must implement {@link ITextEditor}.
 */
class TableEditor {
  /**
   * Creates a new table editor instance.
   *
   * @param {ITextEditor} textEditor - A text editor interface.
   */
  constructor(textEditor) {
    /** @private */
    this._textEditor = textEditor;
    // smart cursor
    /** @private */
    this._scActive = false;
    /** @private */
    this._scTablePos = null;
    /** @private */
    this._scStartFocus = null;
    /** @private */
    this._scLastFocus = null;
  }

  /**
   * Resets the smart cursor.
   * Call this method when the table editor is inactivated.
   *
   * @returns {undefined}
   */
  resetSmartCursor() {
    this._scActive = false;
  }

  /**
   * Checks if the cursor is in a table row.
   * This is useful to check whether the table editor should be activated or not.
   *
   * @returns {boolean} `true` if the cursor is in a table row.
   */
  cursorIsInTable() {
    const pos = this._textEditor.getCursorPosition();
    return this._textEditor.acceptsTableEdit(pos.row)
      && _isTableRow(this._textEditor.getLine(pos.row));
  }

  /**
   * Finds a table under the current cursor position.
   *
   * @private
   * @returns {Object|undefined} An object that contains information about the table;
   * `undefined` if there is no table.
   * The return object contains the properties listed in the table.
   *
   * | property name   | type                                | description                                                              |
   * | --------------- | ----------------------------------- | ------------------------------------------------------------------------ |
   * | `range`         | {@link Range}                       | The range of the table.                                                  |
   * | `lines`         | {@link Array}&lt;{@link string}&gt; | An array of the lines in the range.                                      |
   * | `table`         | {@link Table}                       | A table object read from the text editor.                                |
   * | `focus`         | {@link Focus}                       | A focus object that represents the current cursor position in the table. |
   */
  _findTable() {
    const pos = this._textEditor.getCursorPosition();
    const lastRow = this._textEditor.getLastRow();
    const lines = [];
    let startRow = pos.row;
    let endRow = pos.row;
    // current line
    {
      const line = this._textEditor.getLine(pos.row);
      if (!this._textEditor.acceptsTableEdit(pos.row) || !_isTableRow(line)) {
        return undefined;
      }
      lines.push(line);
    }
    // previous lines
    for (let row = pos.row - 1; row >= 0; row--) {
      const line = this._textEditor.getLine(row);
      if (!this._textEditor.acceptsTableEdit(row) || !_isTableRow(line)) {
        break;
      }
      lines.unshift(line);
      startRow = row;
    }
    // next lines
    for (let row = pos.row + 1; row <= lastRow; row++) {
      const line = this._textEditor.getLine(row);
      if (!this._textEditor.acceptsTableEdit(row) || !_isTableRow(line)) {
        break;
      }
      lines.push(line);
      endRow = row;
    }
    const range = new Range(
      new Point(startRow, 0),
      new Point(endRow, lines[lines.length - 1].length)
    );
    const table = readTable(lines);
    const focus = table.focusOfPosition(pos, startRow);
    return { range, lines, table, focus };
  }

  /**
   * Finds a table and does an operation with it.
   *
   * @private
   * @param {Function} func - A function that does some operation on table information obtained by
   * {@link TableEditor#_findTable}.
   * @returns {undefined}
   */
  _withTable(func) {
    const info = this._findTable();
    if (info === undefined) {
      return;
    }
    func(info);
  }

  /**
   * Updates lines in a given range in the text editor.
   *
   * @private
   * @param {number} startRow - Start row index, starts from `0`.
   * @param {number} endRow - End row index.
   * Lines from `startRow` to `endRow - 1` are replaced.
   * @param {Array<string>} newLines - New lines.
   * @param {Array<string>} [oldLines=undefined] - Old lines to be replaced.
   * @returns {undefined}
   */
  _updateLines(startRow, endRow, newLines, oldLines = undefined) {
    if (oldLines !== undefined) {
      // apply the shortest edit script
      // if a table is edited in a normal manner, the edit distance never exceeds 3
      const ses = shortestEditScript(oldLines, newLines, 3);
      if (ses !== undefined) {
        applyEditScript(this._textEditor, ses, startRow);
        return;
      }
    }
    this._textEditor.replaceLines(startRow, endRow, newLines);
  }

  /**
   * Moves the cursor position to the focused cell,
   *
   * @private
   * @param {number} startRow - Row index where the table starts in the text editor.
   * @param {Table} table - A table.
   * @param {Focus} focus - A focus to which the cursor will be moved.
   * @returns {undefined}
   */
  _moveToFocus(startRow, table, focus) {
    const pos = table.positionOfFocus(focus, startRow);
    if (pos !== undefined) {
      this._textEditor.setCursorPosition(pos);
    }
  }

  /**
   * Selects the focused cell.
   * If the cell has no content to be selected, then just moves the cursor position.
   *
   * @private
   * @param {number} startRow - Row index where the table starts in the text editor.
   * @param {Table} table - A table.
   * @param {Focus} focus - A focus to be selected.
   * @returns {undefined}
   */
  _selectFocus(startRow, table, focus) {
    const range = table.selectionRangeOfFocus(focus, startRow);
    if (range !== undefined) {
      this._textEditor.setSelectionRange(range);
    }
    else {
      this._moveToFocus(startRow, table, focus);
    }
  }

  /**
   * Formats the table under the cursor.
   *
   * @param {Object} options - See {@link options}.
   * @returns {undefined}
   */
  format(options) {
    this._withTable(({ range, lines, table, focus }) => {
      let newFocus = focus;
      // complete
      const completed = completeTable(table, options);
      if (completed.delimiterInserted && newFocus.row > 0) {
        newFocus = newFocus.setRow(newFocus.row + 1);
      }
      // format
      const formatted = formatTable(completed.table, options);
      newFocus = newFocus.setOffset(_computeNewOffset(newFocus, completed.table, formatted, false));
      // apply
      this._textEditor.transact(() => {
        this._updateLines(range.start.row, range.end.row + 1, formatted.table.toLines(), lines);
        this._moveToFocus(range.start.row, formatted.table, newFocus);
      });
    });
  }

  /**
   * Formats and escapes from the table.
   *
   * @param {Object} options - See {@link options}.
   * @returns {undefined}
   */
  escape(options) {
    this._withTable(({ range, lines, table, focus }) => {
      // complete
      const completed = completeTable(table, options);
      // format
      const formatted = formatTable(completed.table, options);
      // apply
      const newPos = new Point(range.end.row + (completed.delimiterInserted ? 2 : 1), 0);
      this._textEditor.transact(() => {
        this._updateLines(range.start.row, range.end.row + 1, formatted.table.toLines(), lines);
        if (newPos.row > this._textEditor.getLastRow()) {
          this._textEditor.insertLine(newPos.row, "");
        }
        this._textEditor.setCursorPosition(newPos);
      });
      this.resetSmartCursor();
    });
  }

  /**
   * Alters the alignment of the focused column.
   *
   * @param {Alignment} alignment - New alignment.
   * @param {Object} options - See {@link options}.
   * @returns {undefined}
   */
  alignColumn(alignment, options) {
    this._withTable(({ range, lines, table, focus }) => {
      let newFocus = focus;
      // complete
      const completed = completeTable(table, options);
      if (completed.delimiterInserted && newFocus.row > 0) {
        newFocus = newFocus.setRow(newFocus.row + 1);
      }
      // alter alignment
      let altered = completed.table;
      if (0 <= newFocus.column && newFocus.column <= altered.getHeaderWidth() - 1) {
        altered = alterAlignment(completed.table, newFocus.column, alignment, options);
      }
      // format
      const formatted = formatTable(altered, options);
      newFocus = newFocus.setOffset(_computeNewOffset(newFocus, completed.table, formatted, false));
      // apply
      this._textEditor.transact(() => {
        this._updateLines(range.start.row, range.end.row + 1, formatted.table.toLines(), lines);
        this._moveToFocus(range.start.row, formatted.table, newFocus);
      });
    });
  }

  /**
   * Selects the focused cell content.
   *
   * @param {Object} options - See {@link options}.
   * @returns {undefined}
   */
  selectCell(options) {
    this._withTable(({ range, lines, table, focus }) => {
      let newFocus = focus;
      // complete
      const completed = completeTable(table, options);
      if (completed.delimiterInserted && newFocus.row > 0) {
        newFocus = newFocus.setRow(newFocus.row + 1);
      }
      // format
      const formatted = formatTable(completed.table, options);
      newFocus = newFocus.setOffset(_computeNewOffset(newFocus, completed.table, formatted, false));
      // apply
      this._textEditor.transact(() => {
        this._updateLines(range.start.row, range.end.row + 1, formatted.table.toLines(), lines);
        this._selectFocus(range.start.row, formatted.table, newFocus);
      });
    });
  }

  /**
   * Moves the focus to another cell.
   *
   * @param {number} rowOffset - Offset in row.
   * @param {number} columnOffset - Offset in column.
   * @param {Object} options - See {@link options}.
   * @returns {undefined}
   */
  moveFocus(rowOffset, columnOffset, options) {
    this._withTable(({ range, lines, table, focus }) => {
      let newFocus = focus;
      // complete
      const completed = completeTable(table, options);
      if (completed.delimiterInserted && newFocus.row > 0) {
        newFocus = newFocus.setRow(newFocus.row + 1);
      }
      const startFocus = newFocus;
      // move focus
      if (rowOffset !== 0) {
        const height = completed.table.getHeight();
        // skip delimiter row
        const skip =
            newFocus.row < 1 && newFocus.row + rowOffset >= 1 ? 1
          : newFocus.row > 1 && newFocus.row + rowOffset <= 1 ? -1
          : 0;
        newFocus = newFocus.setRow(
          Math.min(Math.max(newFocus.row + rowOffset + skip, 0), height <= 2 ? 0 : height - 1)
        );
      }
      if (columnOffset !== 0) {
        const width = completed.table.getHeaderWidth();
        if (!(newFocus.column < 0 && columnOffset < 0)
          && !(newFocus.column > width - 1 && columnOffset > 0)) {
          newFocus = newFocus.setColumn(
            Math.min(Math.max(newFocus.column + columnOffset, 0), width - 1)
          );
        }
      }
      const moved = !newFocus.posEquals(startFocus);
      // format
      const formatted = formatTable(completed.table, options);
      newFocus = newFocus.setOffset(_computeNewOffset(newFocus, completed.table, formatted, moved));
      // apply
      this._textEditor.transact(() => {
        this._updateLines(range.start.row, range.end.row + 1, formatted.table.toLines(), lines);
        if (moved) {
          this._selectFocus(range.start.row, formatted.table, newFocus);
        }
        else {
          this._moveToFocus(range.start.row, formatted.table, newFocus);
        }
      });
      if (moved) {
        this.resetSmartCursor();
      }
    });
  }

  /**
   * Moves the focus to the next cell.
   *
   * @param {Object} options - See {@link options}.
   * @returns {undefined}
   */
  nextCell(options) {
    this._withTable(({ range, lines, table, focus }) => {
      // reset smart cursor if moved
      const focusMoved = (this._scTablePos !== null && !range.start.equals(this._scTablePos))
        || (this._scLastFocus !== null && !focus.posEquals(this._scLastFocus));
      if (this._scActive && focusMoved) {
        this.resetSmartCursor();
      }
      let newFocus = focus;
      // complete
      const completed = completeTable(table, options);
      if (completed.delimiterInserted && newFocus.row > 0) {
        newFocus = newFocus.setRow(newFocus.row + 1);
      }
      const startFocus = newFocus;
      let altered = completed.table;
      // move focus
      if (newFocus.row === 1) {
        // move to next row
        newFocus = newFocus.setRow(2);
        if (options.smartCursor) {
          if (newFocus.column < 0 || altered.getHeaderWidth() - 1 < newFocus.column) {
            newFocus = newFocus.setColumn(0);
          }
        }
        else {
          newFocus = newFocus.setColumn(0);
        }
        // insert an empty row if needed
        if (newFocus.row > altered.getHeight() - 1) {
          const row = new Array(altered.getHeaderWidth()).fill(new TableCell(""));
          altered = insertRow(altered, altered.getHeight(), new TableRow(row, "", ""));
        }
      }
      else {
        // insert an empty column if needed
        if (newFocus.column > altered.getHeaderWidth() - 1) {
          const column = new Array(altered.getHeight() - 1).fill(new TableCell(""));
          altered = insertColumn(altered, altered.getHeaderWidth(), column, options);
        }
        // move to next column
        newFocus = newFocus.setColumn(newFocus.column + 1);
      }
      // format
      const formatted = formatTable(altered, options);
      newFocus = newFocus.setOffset(_computeNewOffset(newFocus, altered, formatted, true));
      // apply
      const newLines = formatted.table.toLines();
      if (newFocus.column > formatted.table.getHeaderWidth() - 1) {
        // add margin
        newLines[newFocus.row] += " ";
        newFocus = newFocus.setOffset(1);
      }
      this._textEditor.transact(() => {
        this._updateLines(range.start.row, range.end.row + 1, newLines, lines);
        this._selectFocus(range.start.row, formatted.table, newFocus);
      });
      if (options.smartCursor) {
        if (!this._scActive) {
          // activate smart cursor
          this._scActive = true;
          this._scTablePos = range.start;
          if (startFocus.column < 0 || formatted.table.getHeaderWidth() - 1 < startFocus.column) {
            this._scStartFocus = new Focus(startFocus.row, 0, 0);
          }
          else {
            this._scStartFocus = startFocus;
          }
        }
        this._scLastFocus = newFocus;
      }
    });
  }

  /**
   * Moves the focus to the previous cell.
   *
   * @param {Object} options - See {@link options}.
   * @returns {undefined}
   */
  previousCell(options) {
    this._withTable(({ range, lines, table, focus }) => {
      let newFocus = focus;
      // complete
      const completed = completeTable(table, options);
      if (completed.delimiterInserted && newFocus.row > 0) {
        newFocus = newFocus.setRow(newFocus.row + 1);
      }
      const startFocus = newFocus;
      // move focus
      if (newFocus.row === 0) {
        if (newFocus.column > 0) {
          newFocus = newFocus.setColumn(newFocus.column - 1);
        }
      }
      else if (newFocus.row === 1) {
        newFocus = new Focus(0, completed.table.getHeaderWidth() - 1, newFocus.offset);
      }
      else {
        if (newFocus.column > 0) {
          newFocus = newFocus.setColumn(newFocus.column - 1);
        }
        else {
          newFocus = new Focus(
            newFocus.row === 2 ? 0 : newFocus.row - 1,
            completed.table.getHeaderWidth() - 1,
            newFocus.offset
          );
        }
      }
      const moved = !newFocus.posEquals(startFocus);
      // format
      const formatted = formatTable(completed.table, options);
      newFocus = newFocus.setOffset(_computeNewOffset(newFocus, completed.table, formatted, moved));
      // apply
      this._textEditor.transact(() => {
        this._updateLines(range.start.row, range.end.row + 1, formatted.table.toLines(), lines);
        if (moved) {
          this._selectFocus(range.start.row, formatted.table, newFocus);
        }
        else {
          this._moveToFocus(range.start.row, formatted.table, newFocus);
        }
      });
      if (moved) {
        this.resetSmartCursor();
      }
    });
  }

  /**
   * Moves the focus to the next row.
   *
   * @param {Object} options - See {@link options}.
   * @returns {undefined}
   */
  nextRow(options) {
    this._withTable(({ range, lines, table, focus }) => {
      // reset smart cursor if moved
      const focusMoved = (this._scTablePos !== null && !range.start.equals(this._scTablePos))
        || (this._scLastFocus !== null && !focus.posEquals(this._scLastFocus));
      if (this._scActive && focusMoved) {
        this.resetSmartCursor();
      }
      let newFocus = focus;
      // complete
      const completed = completeTable(table, options);
      if (completed.delimiterInserted && newFocus.row > 0) {
        newFocus = newFocus.setRow(newFocus.row + 1);
      }
      const startFocus = newFocus;
      let altered = completed.table;
      // move focus
      if (newFocus.row === 0) {
        newFocus = newFocus.setRow(2);
      }
      else {
        newFocus = newFocus.setRow(newFocus.row + 1);
      }
      if (options.smartCursor) {
        if (this._scActive) {
          newFocus = newFocus.setColumn(this._scStartFocus.column);
        }
        else if (newFocus.column < 0 || altered.getHeaderWidth() - 1 < newFocus.column) {
          newFocus = newFocus.setColumn(0);
        }
      }
      else {
        newFocus = newFocus.setColumn(0);
      }
      // insert empty row if needed
      if (newFocus.row > altered.getHeight() - 1) {
        const row = new Array(altered.getHeaderWidth()).fill(new TableCell(""));
        altered = insertRow(altered, altered.getHeight(), new TableRow(row, "", ""));
      }
      // format
      const formatted = formatTable(altered, options);
      newFocus = newFocus.setOffset(_computeNewOffset(newFocus, altered, formatted, true));
      // apply
      this._textEditor.transact(() => {
        this._updateLines(range.start.row, range.end.row + 1, formatted.table.toLines(), lines);
        this._selectFocus(range.start.row, formatted.table, newFocus);
      });
      if (options.smartCursor) {
        if (!this._scActive) {
          // activate smart cursor
          this._scActive = true;
          this._scTablePos = range.start;
          if (startFocus.column < 0 || formatted.table.getHeaderWidth() - 1 < startFocus.column) {
            this._scStartFocus = new Focus(startFocus.row, 0, 0);
          }
          else {
            this._scStartFocus = startFocus;
          }
        }
        this._scLastFocus = newFocus;
      }
    });
  }

  /**
   * Inserts an empty row at the current focus.
   *
   * @param {Object} options - See {@link options}.
   * @returns {undefined}
   */
  insertRow(options) {
    this._withTable(({ range, lines, table, focus }) => {
      let newFocus = focus;
      // complete
      const completed = completeTable(table, options);
      if (completed.delimiterInserted && newFocus.row > 0) {
        newFocus = newFocus.setRow(newFocus.row + 1);
      }
      // move focus
      if (newFocus.row <= 1) {
        newFocus = newFocus.setRow(2);
      }
      newFocus = newFocus.setColumn(0);
      // insert an empty row
      const row = new Array(completed.table.getHeaderWidth()).fill(new TableCell(""));
      const altered = insertRow(completed.table, newFocus.row, new TableRow(row, "", ""));
      // format
      const formatted = formatTable(altered, options);
      newFocus = newFocus.setOffset(_computeNewOffset(newFocus, altered, formatted, true));
      // apply
      this._textEditor.transact(() => {
        this._updateLines(range.start.row, range.end.row + 1, formatted.table.toLines(), lines);
        this._moveToFocus(range.start.row, formatted.table, newFocus);
      });
      this.resetSmartCursor();
    });
  }

  /**
   * Deletes a row at the current focus.
   *
   * @param {Object} options - See {@link options}.
   * @returns {undefined}
   */
  deleteRow(options) {
    this._withTable(({ range, lines, table, focus }) => {
      let newFocus = focus;
      // complete
      const completed = completeTable(table, options);
      if (completed.delimiterInserted && newFocus.row > 0) {
        newFocus = newFocus.setRow(newFocus.row + 1);
      }
      // delete a row
      let altered = completed.table;
      let moved = false;
      if (newFocus.row !== 1) {
        altered = deleteRow(altered, newFocus.row);
        moved = true;
        if (newFocus.row > altered.getHeight() - 1) {
          newFocus = newFocus.setRow(newFocus.row === 2 ? 0 : newFocus.row - 1);
        }
      }
      // format
      const formatted = formatTable(altered, options);
      newFocus = newFocus.setOffset(_computeNewOffset(newFocus, altered, formatted, moved));
      // apply
      this._textEditor.transact(() => {
        this._updateLines(range.start.row, range.end.row + 1, formatted.table.toLines(), lines);
        if (moved) {
          this._selectFocus(range.start.row, formatted.table, newFocus);
        }
        else {
          this._moveToFocus(range.start.row, formatted.table, newFocus);
        }
      });
      this.resetSmartCursor();
    });
  }

  /**
   * Moves the focused row by the specified offset.
   *
   * @param {number} offset - An offset the row is moved by.
   * @param {Object} options - See {@link options}.
   * @returns {undefined}
   */
  moveRow(offset, options) {
    this._withTable(({ range, lines, table, focus }) => {
      let newFocus = focus;
      // complete
      const completed = completeTable(table, options);
      if (completed.delimiterInserted && newFocus.row > 0) {
        newFocus = newFocus.setRow(newFocus.row + 1);
      }
      // move row
      let altered = completed.table;
      if (newFocus.row > 1) {
        const dest = Math.min(Math.max(newFocus.row + offset, 2), altered.getHeight() - 1);
        altered = moveRow(altered, newFocus.row, dest);
        newFocus = newFocus.setRow(dest);
      }
      // format
      const formatted = formatTable(altered, options);
      newFocus = newFocus.setOffset(_computeNewOffset(newFocus, altered, formatted, false));
      // apply
      this._textEditor.transact(() => {
        this._updateLines(range.start.row, range.end.row + 1, formatted.table.toLines(), lines);
        this._moveToFocus(range.start.row, formatted.table, newFocus);
      });
      this.resetSmartCursor();
    });
  }

  /**
   * Inserts an empty column at the current focus.
   *
   * @param {Object} options - See {@link options}.
   * @returns {undefined}
   */
  insertColumn(options) {
    this._withTable(({ range, lines, table, focus }) => {
      let newFocus = focus;
      // complete
      const completed = completeTable(table, options);
      if (completed.delimiterInserted && newFocus.row > 0) {
        newFocus = newFocus.setRow(newFocus.row + 1);
      }
      // move focus
      if (newFocus.row === 1) {
        newFocus = newFocus.setRow(0);
      }
      if (newFocus.column < 0) {
        newFocus = newFocus.setColumn(0);
      }
      // insert an empty column
      const column = new Array(completed.table.getHeight() - 1).fill(new TableCell(""));
      const altered = insertColumn(completed.table, newFocus.column, column, options);
      // format
      const formatted = formatTable(altered, options);
      newFocus = newFocus.setOffset(_computeNewOffset(newFocus, altered, formatted, true));
      // apply
      this._textEditor.transact(() => {
        this._updateLines(range.start.row, range.end.row + 1, formatted.table.toLines(), lines);
        this._moveToFocus(range.start.row, formatted.table, newFocus);
      });
      this.resetSmartCursor();
    });
  }

  /**
   * Deletes a column at the current focus.
   *
   * @param {Object} options - See {@link options}.
   * @returns {undefined}
   */
  deleteColumn(options) {
    this._withTable(({ range, lines, table, focus }) => {
      let newFocus = focus;
      // complete
      const completed = completeTable(table, options);
      if (completed.delimiterInserted && newFocus.row > 0) {
        newFocus = newFocus.setRow(newFocus.row + 1);
      }
      // move focus
      if (newFocus.row === 1) {
        newFocus = newFocus.setRow(0);
      }
      // delete a column
      let altered = completed.table;
      let moved = false;
      if (0 <= newFocus.column && newFocus.column <= altered.getHeaderWidth() - 1) {
        altered = deleteColumn(completed.table, newFocus.column, options);
        moved = true;
        if (newFocus.column > altered.getHeaderWidth() - 1) {
          newFocus = newFocus.setColumn(altered.getHeaderWidth() - 1);
        }
      }
      // format
      const formatted = formatTable(altered, options);
      newFocus = newFocus.setOffset(_computeNewOffset(newFocus, altered, formatted, moved));
      // apply
      this._textEditor.transact(() => {
        this._updateLines(range.start.row, range.end.row + 1, formatted.table.toLines(), lines);
        if (moved) {
          this._selectFocus(range.start.row, formatted.table, newFocus);
        }
        else {
          this._moveToFocus(range.start.row, formatted.table, newFocus);
        }
      });
      this.resetSmartCursor();
    });
  }

  /**
   * Moves the focused column by the specified offset.
   *
   * @param {number} offset - An offset the column is moved by.
   * @param {Object} options - See {@link options}.
   * @returns {undefined}
   */
  moveColumn(offset, options) {
    this._withTable(({ range, lines, table, focus }) => {
      let newFocus = focus;
      // complete
      const completed = completeTable(table, options);
      if (completed.delimiterInserted && newFocus.row > 0) {
        newFocus = newFocus.setRow(newFocus.row + 1);
      }
      // move column
      let altered = completed.table;
      if (0 <= newFocus.column && newFocus.column <= altered.getHeaderWidth() - 1) {
        const dest = Math.min(Math.max(newFocus.column + offset, 0), altered.getHeaderWidth() - 1);
        altered = moveColumn(altered, newFocus.column, dest);
        newFocus = newFocus.setColumn(dest);
      }
      // format
      const formatted = formatTable(altered, options);
      newFocus = newFocus.setOffset(_computeNewOffset(newFocus, altered, formatted, false));
      // apply
      this._textEditor.transact(() => {
        this._updateLines(range.start.row, range.end.row + 1, formatted.table.toLines(), lines);
        this._moveToFocus(range.start.row, formatted.table, newFocus);
      });
      this.resetSmartCursor();
    });
  }

  /**
   * Formats all the tables in the text editor.
   *
   * @param {Object} options - See {@link options}.
   * @returns {undefined}
   */
  formatAll(options) {
    this._textEditor.transact(() => {
      let pos = this._textEditor.getCursorPosition();
      let lines = [];
      let startRow = undefined;
      let lastRow = this._textEditor.getLastRow();
      // find tables
      for (let row = 0; row <= lastRow; row++) {
        const line = this._textEditor.getLine(row);
        if (this._textEditor.acceptsTableEdit(row) && _isTableRow(line)) {
          lines.push(line);
          if (startRow === undefined) {
            startRow = row;
          }
        }
        else if (startRow !== undefined) {
          // get table info
          const endRow = row - 1;
          const range = new Range(
            new Point(startRow, 0),
            new Point(endRow, lines[lines.length - 1].length)
          );
          const table = readTable(lines);
          const focus = table.focusOfPosition(pos, startRow);
          const focused = focus !== undefined;
          // format
          let newFocus = focus;
          const completed = completeTable(table, options);
          if (focused && completed.delimiterInserted && newFocus.row > 0) {
            newFocus = newFocus.setRow(newFocus.row + 1);
          }
          const formatted = formatTable(completed.table, options);
          if (focused) {
            newFocus = newFocus.setOffset(
              _computeNewOffset(newFocus, completed.table, formatted, false)
            );
          }
          // apply
          const newLines = formatted.table.toLines();
          this._updateLines(range.start.row, range.end.row + 1, newLines, lines);
          // update cursor position
          const diff = newLines.length - lines.length;
          if (focused) {
            pos = formatted.table.positionOfFocus(newFocus, startRow);
          }
          else if (pos.row > endRow) {
            pos = new Point(pos.row + diff, pos.column);
          }
          // reset
          lines = [];
          startRow = undefined;
          // update
          lastRow += diff;
          row += diff;
        }
      }
      if (startRow !== undefined) {
        // get table info
        const endRow = lastRow;
        const range = new Range(
          new Point(startRow, 0),
          new Point(endRow, lines[lines.length - 1].length)
        );
        const table = readTable(lines);
        const focus = table.focusOfPosition(pos, startRow);
        // format
        let newFocus = focus;
        const completed = completeTable(table, options);
        if (completed.delimiterInserted && newFocus.row > 0) {
          newFocus = newFocus.setRow(newFocus.row + 1);
        }
        const formatted = formatTable(completed.table, options);
        newFocus = newFocus.setOffset(
          _computeNewOffset(newFocus, completed.table, formatted, false)
        );
        // apply
        const newLines = formatted.table.toLines();
        this._updateLines(range.start.row, range.end.row + 1, newLines, lines);
        pos = formatted.table.positionOfFocus(newFocus, startRow);
      }
      this._textEditor.setCursorPosition(pos);
    });
  }
}


//# sourceMappingURL=mte-kernel.mjs.map


/***/ }),
/* 2 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return getEAW; });
/* unused harmony export computeWidth */
/*
 * This file is generated by a script. DO NOT EDIT BY HAND!
 */

/*
 * This part (from BEGIN to END) is derived from the Unicode Data Files:
 *
 * UNICODE, INC. LICENSE AGREEMENT
 *
 * Copyright © 1991-2017 Unicode, Inc. All rights reserved.
 * Distributed under the Terms of Use in http://www.unicode.org/copyright.html.
 *
 * Permission is hereby granted, free of charge, to any person obtaining
 * a copy of the Unicode data files and any associated documentation
 * (the "Data Files") or Unicode software and any associated documentation
 * (the "Software") to deal in the Data Files or Software
 * without restriction, including without limitation the rights to use,
 * copy, modify, merge, publish, distribute, and/or sell copies of
 * the Data Files or Software, and to permit persons to whom the Data Files
 * or Software are furnished to do so, provided that either
 * (a) this copyright and permission notice appear with all copies
 * of the Data Files or Software, or
 * (b) this copyright and permission notice appear in associated
 * Documentation.
 *
 * THE DATA FILES AND SOFTWARE ARE PROVIDED "AS IS", WITHOUT WARRANTY OF
 * ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
 * WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT OF THIRD PARTY RIGHTS.
 * IN NO EVENT SHALL THE COPYRIGHT HOLDER OR HOLDERS INCLUDED IN THIS
 * NOTICE BE LIABLE FOR ANY CLAIM, OR ANY SPECIAL INDIRECT OR CONSEQUENTIAL
 * DAMAGES, OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE,
 * DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER
 * TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
 * PERFORMANCE OF THE DATA FILES OR SOFTWARE.
 */

/* BEGIN */
var defs = [
  { start: 0, end: 31, prop: "N" },
  { start: 32, end: 126, prop: "Na" },
  { start: 127, end: 160, prop: "N" },
  { start: 161, end: 161, prop: "A" },
  { start: 162, end: 163, prop: "Na" },
  { start: 164, end: 164, prop: "A" },
  { start: 165, end: 166, prop: "Na" },
  { start: 167, end: 168, prop: "A" },
  { start: 169, end: 169, prop: "N" },
  { start: 170, end: 170, prop: "A" },
  { start: 171, end: 171, prop: "N" },
  { start: 172, end: 172, prop: "Na" },
  { start: 173, end: 174, prop: "A" },
  { start: 175, end: 175, prop: "Na" },
  { start: 176, end: 180, prop: "A" },
  { start: 181, end: 181, prop: "N" },
  { start: 182, end: 186, prop: "A" },
  { start: 187, end: 187, prop: "N" },
  { start: 188, end: 191, prop: "A" },
  { start: 192, end: 197, prop: "N" },
  { start: 198, end: 198, prop: "A" },
  { start: 199, end: 207, prop: "N" },
  { start: 208, end: 208, prop: "A" },
  { start: 209, end: 214, prop: "N" },
  { start: 215, end: 216, prop: "A" },
  { start: 217, end: 221, prop: "N" },
  { start: 222, end: 225, prop: "A" },
  { start: 226, end: 229, prop: "N" },
  { start: 230, end: 230, prop: "A" },
  { start: 231, end: 231, prop: "N" },
  { start: 232, end: 234, prop: "A" },
  { start: 235, end: 235, prop: "N" },
  { start: 236, end: 237, prop: "A" },
  { start: 238, end: 239, prop: "N" },
  { start: 240, end: 240, prop: "A" },
  { start: 241, end: 241, prop: "N" },
  { start: 242, end: 243, prop: "A" },
  { start: 244, end: 246, prop: "N" },
  { start: 247, end: 250, prop: "A" },
  { start: 251, end: 251, prop: "N" },
  { start: 252, end: 252, prop: "A" },
  { start: 253, end: 253, prop: "N" },
  { start: 254, end: 254, prop: "A" },
  { start: 255, end: 256, prop: "N" },
  { start: 257, end: 257, prop: "A" },
  { start: 258, end: 272, prop: "N" },
  { start: 273, end: 273, prop: "A" },
  { start: 274, end: 274, prop: "N" },
  { start: 275, end: 275, prop: "A" },
  { start: 276, end: 282, prop: "N" },
  { start: 283, end: 283, prop: "A" },
  { start: 284, end: 293, prop: "N" },
  { start: 294, end: 295, prop: "A" },
  { start: 296, end: 298, prop: "N" },
  { start: 299, end: 299, prop: "A" },
  { start: 300, end: 304, prop: "N" },
  { start: 305, end: 307, prop: "A" },
  { start: 308, end: 311, prop: "N" },
  { start: 312, end: 312, prop: "A" },
  { start: 313, end: 318, prop: "N" },
  { start: 319, end: 322, prop: "A" },
  { start: 323, end: 323, prop: "N" },
  { start: 324, end: 324, prop: "A" },
  { start: 325, end: 327, prop: "N" },
  { start: 328, end: 331, prop: "A" },
  { start: 332, end: 332, prop: "N" },
  { start: 333, end: 333, prop: "A" },
  { start: 334, end: 337, prop: "N" },
  { start: 338, end: 339, prop: "A" },
  { start: 340, end: 357, prop: "N" },
  { start: 358, end: 359, prop: "A" },
  { start: 360, end: 362, prop: "N" },
  { start: 363, end: 363, prop: "A" },
  { start: 364, end: 461, prop: "N" },
  { start: 462, end: 462, prop: "A" },
  { start: 463, end: 463, prop: "N" },
  { start: 464, end: 464, prop: "A" },
  { start: 465, end: 465, prop: "N" },
  { start: 466, end: 466, prop: "A" },
  { start: 467, end: 467, prop: "N" },
  { start: 468, end: 468, prop: "A" },
  { start: 469, end: 469, prop: "N" },
  { start: 470, end: 470, prop: "A" },
  { start: 471, end: 471, prop: "N" },
  { start: 472, end: 472, prop: "A" },
  { start: 473, end: 473, prop: "N" },
  { start: 474, end: 474, prop: "A" },
  { start: 475, end: 475, prop: "N" },
  { start: 476, end: 476, prop: "A" },
  { start: 477, end: 592, prop: "N" },
  { start: 593, end: 593, prop: "A" },
  { start: 594, end: 608, prop: "N" },
  { start: 609, end: 609, prop: "A" },
  { start: 610, end: 707, prop: "N" },
  { start: 708, end: 708, prop: "A" },
  { start: 709, end: 710, prop: "N" },
  { start: 711, end: 711, prop: "A" },
  { start: 712, end: 712, prop: "N" },
  { start: 713, end: 715, prop: "A" },
  { start: 716, end: 716, prop: "N" },
  { start: 717, end: 717, prop: "A" },
  { start: 718, end: 719, prop: "N" },
  { start: 720, end: 720, prop: "A" },
  { start: 721, end: 727, prop: "N" },
  { start: 728, end: 731, prop: "A" },
  { start: 732, end: 732, prop: "N" },
  { start: 733, end: 733, prop: "A" },
  { start: 734, end: 734, prop: "N" },
  { start: 735, end: 735, prop: "A" },
  { start: 736, end: 767, prop: "N" },
  { start: 768, end: 879, prop: "A" },
  { start: 880, end: 912, prop: "N" },
  { start: 913, end: 929, prop: "A" },
  { start: 930, end: 930, prop: "N" },
  { start: 931, end: 937, prop: "A" },
  { start: 938, end: 944, prop: "N" },
  { start: 945, end: 961, prop: "A" },
  { start: 962, end: 962, prop: "N" },
  { start: 963, end: 969, prop: "A" },
  { start: 970, end: 1024, prop: "N" },
  { start: 1025, end: 1025, prop: "A" },
  { start: 1026, end: 1039, prop: "N" },
  { start: 1040, end: 1103, prop: "A" },
  { start: 1104, end: 1104, prop: "N" },
  { start: 1105, end: 1105, prop: "A" },
  { start: 1106, end: 4351, prop: "N" },
  { start: 4352, end: 4447, prop: "W" },
  { start: 4448, end: 8207, prop: "N" },
  { start: 8208, end: 8208, prop: "A" },
  { start: 8209, end: 8210, prop: "N" },
  { start: 8211, end: 8214, prop: "A" },
  { start: 8215, end: 8215, prop: "N" },
  { start: 8216, end: 8217, prop: "A" },
  { start: 8218, end: 8219, prop: "N" },
  { start: 8220, end: 8221, prop: "A" },
  { start: 8222, end: 8223, prop: "N" },
  { start: 8224, end: 8226, prop: "A" },
  { start: 8227, end: 8227, prop: "N" },
  { start: 8228, end: 8231, prop: "A" },
  { start: 8232, end: 8239, prop: "N" },
  { start: 8240, end: 8240, prop: "A" },
  { start: 8241, end: 8241, prop: "N" },
  { start: 8242, end: 8243, prop: "A" },
  { start: 8244, end: 8244, prop: "N" },
  { start: 8245, end: 8245, prop: "A" },
  { start: 8246, end: 8250, prop: "N" },
  { start: 8251, end: 8251, prop: "A" },
  { start: 8252, end: 8253, prop: "N" },
  { start: 8254, end: 8254, prop: "A" },
  { start: 8255, end: 8307, prop: "N" },
  { start: 8308, end: 8308, prop: "A" },
  { start: 8309, end: 8318, prop: "N" },
  { start: 8319, end: 8319, prop: "A" },
  { start: 8320, end: 8320, prop: "N" },
  { start: 8321, end: 8324, prop: "A" },
  { start: 8325, end: 8360, prop: "N" },
  { start: 8361, end: 8361, prop: "H" },
  { start: 8362, end: 8363, prop: "N" },
  { start: 8364, end: 8364, prop: "A" },
  { start: 8365, end: 8450, prop: "N" },
  { start: 8451, end: 8451, prop: "A" },
  { start: 8452, end: 8452, prop: "N" },
  { start: 8453, end: 8453, prop: "A" },
  { start: 8454, end: 8456, prop: "N" },
  { start: 8457, end: 8457, prop: "A" },
  { start: 8458, end: 8466, prop: "N" },
  { start: 8467, end: 8467, prop: "A" },
  { start: 8468, end: 8469, prop: "N" },
  { start: 8470, end: 8470, prop: "A" },
  { start: 8471, end: 8480, prop: "N" },
  { start: 8481, end: 8482, prop: "A" },
  { start: 8483, end: 8485, prop: "N" },
  { start: 8486, end: 8486, prop: "A" },
  { start: 8487, end: 8490, prop: "N" },
  { start: 8491, end: 8491, prop: "A" },
  { start: 8492, end: 8530, prop: "N" },
  { start: 8531, end: 8532, prop: "A" },
  { start: 8533, end: 8538, prop: "N" },
  { start: 8539, end: 8542, prop: "A" },
  { start: 8543, end: 8543, prop: "N" },
  { start: 8544, end: 8555, prop: "A" },
  { start: 8556, end: 8559, prop: "N" },
  { start: 8560, end: 8569, prop: "A" },
  { start: 8570, end: 8584, prop: "N" },
  { start: 8585, end: 8585, prop: "A" },
  { start: 8586, end: 8591, prop: "N" },
  { start: 8592, end: 8601, prop: "A" },
  { start: 8602, end: 8631, prop: "N" },
  { start: 8632, end: 8633, prop: "A" },
  { start: 8634, end: 8657, prop: "N" },
  { start: 8658, end: 8658, prop: "A" },
  { start: 8659, end: 8659, prop: "N" },
  { start: 8660, end: 8660, prop: "A" },
  { start: 8661, end: 8678, prop: "N" },
  { start: 8679, end: 8679, prop: "A" },
  { start: 8680, end: 8703, prop: "N" },
  { start: 8704, end: 8704, prop: "A" },
  { start: 8705, end: 8705, prop: "N" },
  { start: 8706, end: 8707, prop: "A" },
  { start: 8708, end: 8710, prop: "N" },
  { start: 8711, end: 8712, prop: "A" },
  { start: 8713, end: 8714, prop: "N" },
  { start: 8715, end: 8715, prop: "A" },
  { start: 8716, end: 8718, prop: "N" },
  { start: 8719, end: 8719, prop: "A" },
  { start: 8720, end: 8720, prop: "N" },
  { start: 8721, end: 8721, prop: "A" },
  { start: 8722, end: 8724, prop: "N" },
  { start: 8725, end: 8725, prop: "A" },
  { start: 8726, end: 8729, prop: "N" },
  { start: 8730, end: 8730, prop: "A" },
  { start: 8731, end: 8732, prop: "N" },
  { start: 8733, end: 8736, prop: "A" },
  { start: 8737, end: 8738, prop: "N" },
  { start: 8739, end: 8739, prop: "A" },
  { start: 8740, end: 8740, prop: "N" },
  { start: 8741, end: 8741, prop: "A" },
  { start: 8742, end: 8742, prop: "N" },
  { start: 8743, end: 8748, prop: "A" },
  { start: 8749, end: 8749, prop: "N" },
  { start: 8750, end: 8750, prop: "A" },
  { start: 8751, end: 8755, prop: "N" },
  { start: 8756, end: 8759, prop: "A" },
  { start: 8760, end: 8763, prop: "N" },
  { start: 8764, end: 8765, prop: "A" },
  { start: 8766, end: 8775, prop: "N" },
  { start: 8776, end: 8776, prop: "A" },
  { start: 8777, end: 8779, prop: "N" },
  { start: 8780, end: 8780, prop: "A" },
  { start: 8781, end: 8785, prop: "N" },
  { start: 8786, end: 8786, prop: "A" },
  { start: 8787, end: 8799, prop: "N" },
  { start: 8800, end: 8801, prop: "A" },
  { start: 8802, end: 8803, prop: "N" },
  { start: 8804, end: 8807, prop: "A" },
  { start: 8808, end: 8809, prop: "N" },
  { start: 8810, end: 8811, prop: "A" },
  { start: 8812, end: 8813, prop: "N" },
  { start: 8814, end: 8815, prop: "A" },
  { start: 8816, end: 8833, prop: "N" },
  { start: 8834, end: 8835, prop: "A" },
  { start: 8836, end: 8837, prop: "N" },
  { start: 8838, end: 8839, prop: "A" },
  { start: 8840, end: 8852, prop: "N" },
  { start: 8853, end: 8853, prop: "A" },
  { start: 8854, end: 8856, prop: "N" },
  { start: 8857, end: 8857, prop: "A" },
  { start: 8858, end: 8868, prop: "N" },
  { start: 8869, end: 8869, prop: "A" },
  { start: 8870, end: 8894, prop: "N" },
  { start: 8895, end: 8895, prop: "A" },
  { start: 8896, end: 8977, prop: "N" },
  { start: 8978, end: 8978, prop: "A" },
  { start: 8979, end: 8985, prop: "N" },
  { start: 8986, end: 8987, prop: "W" },
  { start: 8988, end: 9000, prop: "N" },
  { start: 9001, end: 9002, prop: "W" },
  { start: 9003, end: 9192, prop: "N" },
  { start: 9193, end: 9196, prop: "W" },
  { start: 9197, end: 9199, prop: "N" },
  { start: 9200, end: 9200, prop: "W" },
  { start: 9201, end: 9202, prop: "N" },
  { start: 9203, end: 9203, prop: "W" },
  { start: 9204, end: 9311, prop: "N" },
  { start: 9312, end: 9449, prop: "A" },
  { start: 9450, end: 9450, prop: "N" },
  { start: 9451, end: 9547, prop: "A" },
  { start: 9548, end: 9551, prop: "N" },
  { start: 9552, end: 9587, prop: "A" },
  { start: 9588, end: 9599, prop: "N" },
  { start: 9600, end: 9615, prop: "A" },
  { start: 9616, end: 9617, prop: "N" },
  { start: 9618, end: 9621, prop: "A" },
  { start: 9622, end: 9631, prop: "N" },
  { start: 9632, end: 9633, prop: "A" },
  { start: 9634, end: 9634, prop: "N" },
  { start: 9635, end: 9641, prop: "A" },
  { start: 9642, end: 9649, prop: "N" },
  { start: 9650, end: 9651, prop: "A" },
  { start: 9652, end: 9653, prop: "N" },
  { start: 9654, end: 9655, prop: "A" },
  { start: 9656, end: 9659, prop: "N" },
  { start: 9660, end: 9661, prop: "A" },
  { start: 9662, end: 9663, prop: "N" },
  { start: 9664, end: 9665, prop: "A" },
  { start: 9666, end: 9669, prop: "N" },
  { start: 9670, end: 9672, prop: "A" },
  { start: 9673, end: 9674, prop: "N" },
  { start: 9675, end: 9675, prop: "A" },
  { start: 9676, end: 9677, prop: "N" },
  { start: 9678, end: 9681, prop: "A" },
  { start: 9682, end: 9697, prop: "N" },
  { start: 9698, end: 9701, prop: "A" },
  { start: 9702, end: 9710, prop: "N" },
  { start: 9711, end: 9711, prop: "A" },
  { start: 9712, end: 9724, prop: "N" },
  { start: 9725, end: 9726, prop: "W" },
  { start: 9727, end: 9732, prop: "N" },
  { start: 9733, end: 9734, prop: "A" },
  { start: 9735, end: 9736, prop: "N" },
  { start: 9737, end: 9737, prop: "A" },
  { start: 9738, end: 9741, prop: "N" },
  { start: 9742, end: 9743, prop: "A" },
  { start: 9744, end: 9747, prop: "N" },
  { start: 9748, end: 9749, prop: "W" },
  { start: 9750, end: 9755, prop: "N" },
  { start: 9756, end: 9756, prop: "A" },
  { start: 9757, end: 9757, prop: "N" },
  { start: 9758, end: 9758, prop: "A" },
  { start: 9759, end: 9791, prop: "N" },
  { start: 9792, end: 9792, prop: "A" },
  { start: 9793, end: 9793, prop: "N" },
  { start: 9794, end: 9794, prop: "A" },
  { start: 9795, end: 9799, prop: "N" },
  { start: 9800, end: 9811, prop: "W" },
  { start: 9812, end: 9823, prop: "N" },
  { start: 9824, end: 9825, prop: "A" },
  { start: 9826, end: 9826, prop: "N" },
  { start: 9827, end: 9829, prop: "A" },
  { start: 9830, end: 9830, prop: "N" },
  { start: 9831, end: 9834, prop: "A" },
  { start: 9835, end: 9835, prop: "N" },
  { start: 9836, end: 9837, prop: "A" },
  { start: 9838, end: 9838, prop: "N" },
  { start: 9839, end: 9839, prop: "A" },
  { start: 9840, end: 9854, prop: "N" },
  { start: 9855, end: 9855, prop: "W" },
  { start: 9856, end: 9874, prop: "N" },
  { start: 9875, end: 9875, prop: "W" },
  { start: 9876, end: 9885, prop: "N" },
  { start: 9886, end: 9887, prop: "A" },
  { start: 9888, end: 9888, prop: "N" },
  { start: 9889, end: 9889, prop: "W" },
  { start: 9890, end: 9897, prop: "N" },
  { start: 9898, end: 9899, prop: "W" },
  { start: 9900, end: 9916, prop: "N" },
  { start: 9917, end: 9918, prop: "W" },
  { start: 9919, end: 9919, prop: "A" },
  { start: 9920, end: 9923, prop: "N" },
  { start: 9924, end: 9925, prop: "W" },
  { start: 9926, end: 9933, prop: "A" },
  { start: 9934, end: 9934, prop: "W" },
  { start: 9935, end: 9939, prop: "A" },
  { start: 9940, end: 9940, prop: "W" },
  { start: 9941, end: 9953, prop: "A" },
  { start: 9954, end: 9954, prop: "N" },
  { start: 9955, end: 9955, prop: "A" },
  { start: 9956, end: 9959, prop: "N" },
  { start: 9960, end: 9961, prop: "A" },
  { start: 9962, end: 9962, prop: "W" },
  { start: 9963, end: 9969, prop: "A" },
  { start: 9970, end: 9971, prop: "W" },
  { start: 9972, end: 9972, prop: "A" },
  { start: 9973, end: 9973, prop: "W" },
  { start: 9974, end: 9977, prop: "A" },
  { start: 9978, end: 9978, prop: "W" },
  { start: 9979, end: 9980, prop: "A" },
  { start: 9981, end: 9981, prop: "W" },
  { start: 9982, end: 9983, prop: "A" },
  { start: 9984, end: 9988, prop: "N" },
  { start: 9989, end: 9989, prop: "W" },
  { start: 9990, end: 9993, prop: "N" },
  { start: 9994, end: 9995, prop: "W" },
  { start: 9996, end: 10023, prop: "N" },
  { start: 10024, end: 10024, prop: "W" },
  { start: 10025, end: 10044, prop: "N" },
  { start: 10045, end: 10045, prop: "A" },
  { start: 10046, end: 10059, prop: "N" },
  { start: 10060, end: 10060, prop: "W" },
  { start: 10061, end: 10061, prop: "N" },
  { start: 10062, end: 10062, prop: "W" },
  { start: 10063, end: 10066, prop: "N" },
  { start: 10067, end: 10069, prop: "W" },
  { start: 10070, end: 10070, prop: "N" },
  { start: 10071, end: 10071, prop: "W" },
  { start: 10072, end: 10101, prop: "N" },
  { start: 10102, end: 10111, prop: "A" },
  { start: 10112, end: 10132, prop: "N" },
  { start: 10133, end: 10135, prop: "W" },
  { start: 10136, end: 10159, prop: "N" },
  { start: 10160, end: 10160, prop: "W" },
  { start: 10161, end: 10174, prop: "N" },
  { start: 10175, end: 10175, prop: "W" },
  { start: 10176, end: 10213, prop: "N" },
  { start: 10214, end: 10221, prop: "Na" },
  { start: 10222, end: 10628, prop: "N" },
  { start: 10629, end: 10630, prop: "Na" },
  { start: 10631, end: 11034, prop: "N" },
  { start: 11035, end: 11036, prop: "W" },
  { start: 11037, end: 11087, prop: "N" },
  { start: 11088, end: 11088, prop: "W" },
  { start: 11089, end: 11092, prop: "N" },
  { start: 11093, end: 11093, prop: "W" },
  { start: 11094, end: 11097, prop: "A" },
  { start: 11098, end: 11903, prop: "N" },
  { start: 11904, end: 11929, prop: "W" },
  { start: 11930, end: 11930, prop: "N" },
  { start: 11931, end: 12019, prop: "W" },
  { start: 12020, end: 12031, prop: "N" },
  { start: 12032, end: 12245, prop: "W" },
  { start: 12246, end: 12271, prop: "N" },
  { start: 12272, end: 12283, prop: "W" },
  { start: 12284, end: 12287, prop: "N" },
  { start: 12288, end: 12288, prop: "F" },
  { start: 12289, end: 12350, prop: "W" },
  { start: 12351, end: 12352, prop: "N" },
  { start: 12353, end: 12438, prop: "W" },
  { start: 12439, end: 12440, prop: "N" },
  { start: 12441, end: 12543, prop: "W" },
  { start: 12544, end: 12548, prop: "N" },
  { start: 12549, end: 12590, prop: "W" },
  { start: 12591, end: 12592, prop: "N" },
  { start: 12593, end: 12686, prop: "W" },
  { start: 12687, end: 12687, prop: "N" },
  { start: 12688, end: 12730, prop: "W" },
  { start: 12731, end: 12735, prop: "N" },
  { start: 12736, end: 12771, prop: "W" },
  { start: 12772, end: 12783, prop: "N" },
  { start: 12784, end: 12830, prop: "W" },
  { start: 12831, end: 12831, prop: "N" },
  { start: 12832, end: 12871, prop: "W" },
  { start: 12872, end: 12879, prop: "A" },
  { start: 12880, end: 13054, prop: "W" },
  { start: 13055, end: 13055, prop: "N" },
  { start: 13056, end: 19903, prop: "W" },
  { start: 19904, end: 19967, prop: "N" },
  { start: 19968, end: 42124, prop: "W" },
  { start: 42125, end: 42127, prop: "N" },
  { start: 42128, end: 42182, prop: "W" },
  { start: 42183, end: 43359, prop: "N" },
  { start: 43360, end: 43388, prop: "W" },
  { start: 43389, end: 44031, prop: "N" },
  { start: 44032, end: 55203, prop: "W" },
  { start: 55204, end: 57343, prop: "N" },
  { start: 57344, end: 63743, prop: "A" },
  { start: 63744, end: 64255, prop: "W" },
  { start: 64256, end: 65023, prop: "N" },
  { start: 65024, end: 65039, prop: "A" },
  { start: 65040, end: 65049, prop: "W" },
  { start: 65050, end: 65071, prop: "N" },
  { start: 65072, end: 65106, prop: "W" },
  { start: 65107, end: 65107, prop: "N" },
  { start: 65108, end: 65126, prop: "W" },
  { start: 65127, end: 65127, prop: "N" },
  { start: 65128, end: 65131, prop: "W" },
  { start: 65132, end: 65280, prop: "N" },
  { start: 65281, end: 65376, prop: "F" },
  { start: 65377, end: 65470, prop: "H" },
  { start: 65471, end: 65473, prop: "N" },
  { start: 65474, end: 65479, prop: "H" },
  { start: 65480, end: 65481, prop: "N" },
  { start: 65482, end: 65487, prop: "H" },
  { start: 65488, end: 65489, prop: "N" },
  { start: 65490, end: 65495, prop: "H" },
  { start: 65496, end: 65497, prop: "N" },
  { start: 65498, end: 65500, prop: "H" },
  { start: 65501, end: 65503, prop: "N" },
  { start: 65504, end: 65510, prop: "F" },
  { start: 65511, end: 65511, prop: "N" },
  { start: 65512, end: 65518, prop: "H" },
  { start: 65519, end: 65532, prop: "N" },
  { start: 65533, end: 65533, prop: "A" },
  { start: 65534, end: 94175, prop: "N" },
  { start: 94176, end: 94177, prop: "W" },
  { start: 94178, end: 94207, prop: "N" },
  { start: 94208, end: 100332, prop: "W" },
  { start: 100333, end: 100351, prop: "N" },
  { start: 100352, end: 101106, prop: "W" },
  { start: 101107, end: 110591, prop: "N" },
  { start: 110592, end: 110878, prop: "W" },
  { start: 110879, end: 110959, prop: "N" },
  { start: 110960, end: 111355, prop: "W" },
  { start: 111356, end: 126979, prop: "N" },
  { start: 126980, end: 126980, prop: "W" },
  { start: 126981, end: 127182, prop: "N" },
  { start: 127183, end: 127183, prop: "W" },
  { start: 127184, end: 127231, prop: "N" },
  { start: 127232, end: 127242, prop: "A" },
  { start: 127243, end: 127247, prop: "N" },
  { start: 127248, end: 127277, prop: "A" },
  { start: 127278, end: 127279, prop: "N" },
  { start: 127280, end: 127337, prop: "A" },
  { start: 127338, end: 127343, prop: "N" },
  { start: 127344, end: 127373, prop: "A" },
  { start: 127374, end: 127374, prop: "W" },
  { start: 127375, end: 127376, prop: "A" },
  { start: 127377, end: 127386, prop: "W" },
  { start: 127387, end: 127404, prop: "A" },
  { start: 127405, end: 127487, prop: "N" },
  { start: 127488, end: 127490, prop: "W" },
  { start: 127491, end: 127503, prop: "N" },
  { start: 127504, end: 127547, prop: "W" },
  { start: 127548, end: 127551, prop: "N" },
  { start: 127552, end: 127560, prop: "W" },
  { start: 127561, end: 127567, prop: "N" },
  { start: 127568, end: 127569, prop: "W" },
  { start: 127570, end: 127583, prop: "N" },
  { start: 127584, end: 127589, prop: "W" },
  { start: 127590, end: 127743, prop: "N" },
  { start: 127744, end: 127776, prop: "W" },
  { start: 127777, end: 127788, prop: "N" },
  { start: 127789, end: 127797, prop: "W" },
  { start: 127798, end: 127798, prop: "N" },
  { start: 127799, end: 127868, prop: "W" },
  { start: 127869, end: 127869, prop: "N" },
  { start: 127870, end: 127891, prop: "W" },
  { start: 127892, end: 127903, prop: "N" },
  { start: 127904, end: 127946, prop: "W" },
  { start: 127947, end: 127950, prop: "N" },
  { start: 127951, end: 127955, prop: "W" },
  { start: 127956, end: 127967, prop: "N" },
  { start: 127968, end: 127984, prop: "W" },
  { start: 127985, end: 127987, prop: "N" },
  { start: 127988, end: 127988, prop: "W" },
  { start: 127989, end: 127991, prop: "N" },
  { start: 127992, end: 128062, prop: "W" },
  { start: 128063, end: 128063, prop: "N" },
  { start: 128064, end: 128064, prop: "W" },
  { start: 128065, end: 128065, prop: "N" },
  { start: 128066, end: 128252, prop: "W" },
  { start: 128253, end: 128254, prop: "N" },
  { start: 128255, end: 128317, prop: "W" },
  { start: 128318, end: 128330, prop: "N" },
  { start: 128331, end: 128334, prop: "W" },
  { start: 128335, end: 128335, prop: "N" },
  { start: 128336, end: 128359, prop: "W" },
  { start: 128360, end: 128377, prop: "N" },
  { start: 128378, end: 128378, prop: "W" },
  { start: 128379, end: 128404, prop: "N" },
  { start: 128405, end: 128406, prop: "W" },
  { start: 128407, end: 128419, prop: "N" },
  { start: 128420, end: 128420, prop: "W" },
  { start: 128421, end: 128506, prop: "N" },
  { start: 128507, end: 128591, prop: "W" },
  { start: 128592, end: 128639, prop: "N" },
  { start: 128640, end: 128709, prop: "W" },
  { start: 128710, end: 128715, prop: "N" },
  { start: 128716, end: 128716, prop: "W" },
  { start: 128717, end: 128719, prop: "N" },
  { start: 128720, end: 128722, prop: "W" },
  { start: 128723, end: 128746, prop: "N" },
  { start: 128747, end: 128748, prop: "W" },
  { start: 128749, end: 128755, prop: "N" },
  { start: 128756, end: 128760, prop: "W" },
  { start: 128761, end: 129295, prop: "N" },
  { start: 129296, end: 129342, prop: "W" },
  { start: 129343, end: 129343, prop: "N" },
  { start: 129344, end: 129356, prop: "W" },
  { start: 129357, end: 129359, prop: "N" },
  { start: 129360, end: 129387, prop: "W" },
  { start: 129388, end: 129407, prop: "N" },
  { start: 129408, end: 129431, prop: "W" },
  { start: 129432, end: 129471, prop: "N" },
  { start: 129472, end: 129472, prop: "W" },
  { start: 129473, end: 129487, prop: "N" },
  { start: 129488, end: 129510, prop: "W" },
  { start: 129511, end: 131071, prop: "N" },
  { start: 131072, end: 196605, prop: "W" },
  { start: 196606, end: 196607, prop: "N" },
  { start: 196608, end: 262141, prop: "W" },
  { start: 262142, end: 917759, prop: "N" },
  { start: 917760, end: 917999, prop: "A" },
  { start: 918000, end: 983039, prop: "N" },
  { start: 983040, end: 1048573, prop: "A" },
  { start: 1048574, end: 1048575, prop: "N" },
  { start: 1048576, end: 1114109, prop: "A" },
  { start: 1114110, end: 1114111, prop: "N" }
];
/* END */

/**
 * Returns The EAW property of a code point.
 * @private
 * @param {string} codePoint A code point
 * @return {string} The EAW property of the specified code point
 */
function _getEAWOfCodePoint(codePoint) {
  let min = 0;
  let max = defs.length - 1;
  while (min !== max) {
    const i   = min + ((max - min) >> 1);
    const def = defs[i];
    if (codePoint < def.start) {
      max = i - 1;
    }
    else if (codePoint > def.end) {
      min = i + 1;
    }
    else {
      return def.prop;
    }
  }
  return defs[min].prop;
}

/**
 * Returns the EAW property of a character.
 * @param {string} str A string in which the character is contained
 * @param {number} [at = 0] The position (in code unit) of the character in the string
 * @return {string} The EAW property of the specified character
 * @example
 * import { getEAW } from "meaw";
 *
 * // Narrow
 * assert(getEAW("A") === "Na");
 * // Wide
 * assert(getEAW("あ") === "W");
 * assert(getEAW("安") === "W");
 * assert(getEAW("🍣") === "W");
 * // Fullwidth
 * assert(getEAW("Ａ") === "F");
 * // Halfwidth
 * assert(getEAW("ｱ") === "H");
 * // Ambiguous
 * assert(getEAW("∀") === "A");
 * assert(getEAW("→") === "A");
 * assert(getEAW("Ω") === "A");
 * assert(getEAW("Я") === "A");
 * // Neutral
 * assert(getEAW("ℵ") === "N");
 *
 * // a position (in code unit) can be specified
 * assert(getEAW("ℵAあＡｱ∀", 2) === "W");
 */
function getEAW(str, at) {
  const codePoint = str.codePointAt(at || 0);
  return codePoint === undefined
    ? undefined
    : _getEAWOfCodePoint(codePoint);
}

const defaultWidthMap = {
  "N" : 1,
  "Na": 1,
  "W" : 2,
  "F" : 2,
  "H" : 1,
  "A" : 1
};

/**
 * Computes width of a string based on the EAW properties of its characters.
 * By default characters with property Wide (W) or Fullwidth (F) are treated as wide (= 2)
 * and the others are as narrow (= 1)
 * @param {string} str A string to compute width
 * @param {Object<string, number> | undefined} [widthMap = undefined]
 *   An object which represents a map from an EAW property to a character width
 * @return {number} The computed width
 * @example
 * import { computeWidth } from "meaw";
 *
 * assert(computeWidth("Aあ🍣Ω") === 6);
 * // custom widths can be specified by an object
 * assert(computeWidth("Aあ🍣Ω", { "A": 2 }) === 7);
 */
function computeWidth(str, widthMap) {
  const map = widthMap
    ? Object.assign({}, defaultWidthMap, widthMap)
    : defaultWidthMap;
  let width = 0;
  for (const char of str) {
    width += map[getEAW(char)];
  }
  return width;
}


//# sourceMappingURL=meaw.es.js.map


/***/ })
/******/ ]);
});