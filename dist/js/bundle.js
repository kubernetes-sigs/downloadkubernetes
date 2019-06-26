/******/ (function(modules) { // webpackBootstrap
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
/******/ 			Object.defineProperty(exports, name, { enumerable: true, get: getter });
/******/ 		}
/******/ 	};
/******/
/******/ 	// define __esModule on exports
/******/ 	__webpack_require__.r = function(exports) {
/******/ 		if(typeof Symbol !== 'undefined' && Symbol.toStringTag) {
/******/ 			Object.defineProperty(exports, Symbol.toStringTag, { value: 'Module' });
/******/ 		}
/******/ 		Object.defineProperty(exports, '__esModule', { value: true });
/******/ 	};
/******/
/******/ 	// create a fake namespace object
/******/ 	// mode & 1: value is a module id, require it
/******/ 	// mode & 2: merge all properties of value into the ns
/******/ 	// mode & 4: return value when already ns object
/******/ 	// mode & 8|1: behave like require
/******/ 	__webpack_require__.t = function(value, mode) {
/******/ 		if(mode & 1) value = __webpack_require__(value);
/******/ 		if(mode & 8) return value;
/******/ 		if((mode & 4) && typeof value === 'object' && value && value.__esModule) return value;
/******/ 		var ns = Object.create(null);
/******/ 		__webpack_require__.r(ns);
/******/ 		Object.defineProperty(ns, 'default', { enumerable: true, value: value });
/******/ 		if(mode & 2 && typeof value != 'string') for(var key in value) __webpack_require__.d(ns, key, function(key) { return value[key]; }.bind(null, key));
/******/ 		return ns;
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
/******/
/******/ 	// Load entry module and return exports
/******/ 	return __webpack_require__(__webpack_require__.s = "./src/index.js");
/******/ })
/************************************************************************/
/******/ ({

/***/ "./src/index.js":
/*!**********************!*\
  !*** ./src/index.js ***!
  \**********************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

eval("__webpack_require__(/*! ./mystyles.scss */ \"./src/mystyles.scss\");\n\n[\"os\", \"arch\", \"version\"].map(function(val){\n    eventListener(val);\n});\n\nfunction eventListener(kind) {\n    document.getElementById(kind).addEventListener(\"click\", function(evt){\n        let hideClass = kind + \"-hide\";\n        let k = evt.target.dataset[kind];\n        let notSelector = \"tbody tr:not(.\"+k+\")\";\n        let selector = \"tbody tr.\"+k;\n\n        var rows;\n\n        // if this is the actively selected one then clear stuff\n        if (evt.srcElement.classList.contains(\"is-success\")) {\n            evt.srcElement.classList.remove(\"is-success\")\n            rows = document.querySelectorAll(notSelector);\n            for (let i=0; i < rows.length; i++) {\n                rows.item(i).classList.remove(hideClass)\n            }\n            return\n        }\n\n        for (let i=0; i < this.children.length; i++) {\n            if (this.children.item(i).dataset[kind] === evt.srcElement.dataset[kind]) {\n                evt.srcElement.classList.add(\"is-success\");\n                continue\n            }\n            this.children.item(i).classList.remove(\"is-success\");\n        }\n        rows = document.querySelectorAll(notSelector);\n        for (let i=0; i < rows.length; i++) {\n            rows.item(i).classList.add(hideClass)\n        }\n        rows = document.querySelectorAll(selector);\n        for (let i=0; i < rows.length; i++) {\n            rows.item(i).classList.remove(hideClass)\n        }\n    });\n}\n\nlet links = document.getElementsByClassName(\"copy\");\nfor (let i=0; i < links.length; i++) {\n    links.item(i).addEventListener('click', function(evt) {\n        evt.preventDefault();\n\n        let el = document.createElement(\"textarea\");\n        el.value = this.href;\n        el.setAttribute(\"readonly\", \"\");\n        el.style.position = 'absolute';\n        el.style.left = '-9999px';\n        document.body.appendChild(el);\n        el.select();\n        document.execCommand(\"copy\");\n        document.body.removeChild(el);\n    });\n}\n\n\n\n//# sourceURL=webpack:///./src/index.js?");

/***/ }),

/***/ "./src/mystyles.scss":
/*!***************************!*\
  !*** ./src/mystyles.scss ***!
  \***************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

eval("// extracted by mini-css-extract-plugin\n\n//# sourceURL=webpack:///./src/mystyles.scss?");

/***/ })

/******/ });