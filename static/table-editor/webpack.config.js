// webpack.config.js

const path = require('path');

module.exports = {
  entry: './src/main.js', // 库的入口文件
  output: {
    path: path.resolve(__dirname, 'dist'), // 输出目录
    filename: 'index.js', // 输出文件名称
    library: 'TableEditor', // 库的全局变量名
    libraryTarget: 'umd', // 输出库的目标格式
    umdNamedDefine: true // 对UMD模块进行命名定义
  },
  // 配置其他选项...
};
