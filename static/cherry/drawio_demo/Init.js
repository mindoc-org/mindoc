// urlParams is null when used for embedding
window.urlParams = window.urlParams || {};

// Public global variables
window.DOM_PURIFY_CONFIG = window.DOM_PURIFY_CONFIG ||
    {ADD_TAGS: ['use'], FORBID_TAGS: ['form'],
    ALLOWED_URI_REGEXP: /^((?!javascript:).)*$/i,
    ADD_ATTR: ['target', 'content']};
// Public global variables
window.MAX_REQUEST_SIZE = window.MAX_REQUEST_SIZE  || 10485760;
window.MAX_AREA = window.MAX_AREA || 15000 * 15000;

// URLs for save and export
window.EXPORT_URL = window.EXPORT_URL || './drawio_demo';
window.SAVE_URL = window.SAVE_URL || './drawio_demo';
window.PROXY_URL = window.PROXY_URL || null;
window.OPEN_URL = window.OPEN_URL || './drawio_demo';
window.RESOURCES_PATH = window.RESOURCES_PATH || './drawio_demo/resources';
window.STENCIL_PATH = window.STENCIL_PATH || './drawio_demo/image/stencils';
window.IMAGE_PATH = window.IMAGE_PATH || './drawio_demo/image';
window.STYLE_PATH = window.STYLE_PATH || './drawio_demo/src/css';
window.CSS_PATH = window.CSS_PATH || './drawio_demo/src/css';
window.OPEN_FORM = window.OPEN_FORM ||  './drawio_demo';

window.mxBasePath = window.mxBasePath || './drawio_demo/src';
window.mxLanguage = window.mxLanguage || window.RESOURCES_PATH + '/zh';
window.mxLanguages = window.mxLanguages || ['zh'];
