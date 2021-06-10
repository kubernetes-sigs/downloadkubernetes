/*
 * ATTENTION: The "eval" devtool has been used (maybe by default in mode: "development").
 * This devtool is neither made for production nor for readable output files.
 * It uses "eval()" calls to create a separate source file in the browser devtools.
 * If you are trying to read the output file, select a different devtool (https://webpack.js.org/configuration/devtool/)
 * or disable the default devtool with "devtool: false".
 * If you are looking for production-ready output files, see mode: "production" (https://webpack.js.org/configuration/mode/).
 */
/******/ (() => { // webpackBootstrap
/******/ 	var __webpack_modules__ = ({

/***/ "./src/mystyles.scss":
/*!***************************!*\
  !*** ./src/mystyles.scss ***!
  \***************************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

"use strict";
eval("__webpack_require__.r(__webpack_exports__);\n// extracted by mini-css-extract-plugin\n\n\n//# sourceURL=webpack://downloadkubernetes/./src/mystyles.scss?");

/***/ }),

/***/ "./src/index.js":
/*!**********************!*\
  !*** ./src/index.js ***!
  \**********************/
/***/ ((__unused_webpack_module, __unused_webpack_exports, __webpack_require__) => {

eval("/*\nCopyright 2020 The Kubernetes Authors.\n\nLicensed under the Apache License, Version 2.0 (the \"License\");\nyou may not use this file except in compliance with the License.\nYou may obtain a copy of the License at\n\n    http://www.apache.org/licenses/LICENSE-2.0\n\nUnless required by applicable law or agreed to in writing, software\ndistributed under the License is distributed on an \"AS IS\" BASIS,\nWITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\nSee the License for the specific language governing permissions and\nlimitations under the License.\n*/\n\n__webpack_require__(/*! ./mystyles.scss */ \"./src/mystyles.scss\");\n\n[\"os\", \"arch\", \"version\"].map(function(val){\n    eventListener(val);\n});\n\nfunction eventListener(kind) {\n    let buttonGroupQuery = '#' + kind + ' > .button';\n    let buttonGroup = document.querySelectorAll(buttonGroupQuery);\n\n    let rowsQuery = 'tbody tr';\n    let rows = document.querySelectorAll(rowsQuery);\n\n    let hideClass = kind + \"-hide\";\n\n    buttonGroup.forEach(button => {\n        button.addEventListener('click', (evt) => {\n            let buttonData = button.dataset[kind];\n            buttonGroup.forEach(b => {\n                b.classList.remove('is-success');\n            })\n            button.classList.add('is-success');\n            rows.forEach(row => {\n                if (row.classList.contains(buttonData)) {\n                    row.classList.remove(hideClass);\n                } else {\n                    row.classList.add(hideClass);\n                }\n            });\n        });\n    });\n}\n\n// make the click link work\ndocument.querySelectorAll(\".copy\").forEach(link => {\n    link.addEventListener('click', (evt) => {\n        evt.preventDefault();\n\n        let el = document.createElement(\"textarea\");\n        el.value = link.href;\n        el.setAttribute(\"readonly\", \"\");\n        el.style.position = 'absolute';\n        el.style.left = '-9999px';\n        document.body.appendChild(el);\n        el.select();\n        document.execCommand(\"copy\");\n        document.body.removeChild(el);\n        return false;\n    });\n});\n\n\n//# sourceURL=webpack://downloadkubernetes/./src/index.js?");

/***/ })

/******/ 	});
/************************************************************************/
/******/ 	// The module cache
/******/ 	var __webpack_module_cache__ = {};
/******/ 	
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/ 		// Check if module is in cache
/******/ 		var cachedModule = __webpack_module_cache__[moduleId];
/******/ 		if (cachedModule !== undefined) {
/******/ 			return cachedModule.exports;
/******/ 		}
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = __webpack_module_cache__[moduleId] = {
/******/ 			// no module.id needed
/******/ 			// no module.loaded needed
/******/ 			exports: {}
/******/ 		};
/******/ 	
/******/ 		// Execute the module function
/******/ 		__webpack_modules__[moduleId](module, module.exports, __webpack_require__);
/******/ 	
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/ 	
/************************************************************************/
/******/ 	/* webpack/runtime/make namespace object */
/******/ 	(() => {
/******/ 		// define __esModule on exports
/******/ 		__webpack_require__.r = (exports) => {
/******/ 			if(typeof Symbol !== 'undefined' && Symbol.toStringTag) {
/******/ 				Object.defineProperty(exports, Symbol.toStringTag, { value: 'Module' });
/******/ 			}
/******/ 			Object.defineProperty(exports, '__esModule', { value: true });
/******/ 		};
/******/ 	})();
/******/ 	
/************************************************************************/
/******/ 	
/******/ 	// startup
/******/ 	// Load entry module and return exports
/******/ 	// This entry module can't be inlined because the eval devtool is used.
/******/ 	var __webpack_exports__ = __webpack_require__("./src/index.js");
/******/ 	
/******/ })()
;