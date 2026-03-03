export default class MermaidCodeEngine {
    static TYPE: string;
    static install(cherryOptions: any, ...args: any[]): void;
    constructor(mermaidOptions?: {});
    mermaidAPIRefs: any;
    options: {
        theme: string;
        altFontFamily: string;
        fontFamily: string;
        themeCSS: string;
        flowchart: {
            useMaxWidth: boolean;
        };
        sequence: {
            useMaxWidth: boolean;
        };
        startOnLoad: boolean;
        logLevel: number;
    };
    dom: any;
    mermaidCanvas: any;
    mountMermaidCanvas($engine: any): void;
    /**
     * 转换svg为img，如果出错则直出svg
     * @param {string} svgCode
     * @param {string} graphId
     * @returns {string}
     */
    convertMermaidSvgToImg(svgCode: string, graphId: string): string;
    render(src: any, sign: any, $engine: any): boolean;
}
