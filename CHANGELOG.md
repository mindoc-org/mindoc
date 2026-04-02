# 更新日志 (Changelog)

## [V2.1.1] - 2024-04-02

### 新增 (Added)
- ✨ 目录树默认收起功能
  - 编辑界面和阅读界面的目录树默认全部收起
  - 提供更清爽的界面体验
  
- 💾 目录树状态自动记忆功能
  - 自动保存用户展开/收起的节点
  - 刷新页面后自动恢复上次状态
  - 使用浏览器localStorage存储，不占用服务器资源
  
- 🔄 编辑和阅读界面状态独立管理
  - 编辑界面和阅读界面的展开状态分别存储
  - 使用不同的localStorage键名区分
  
- 🏢 不同项目状态隔离
  - 每个项目的目录树状态独立保存
  - 切换项目时自动加载对应状态

### 优化 (Improved)
- ⚡ 消除页面刷新时的目录树闪烁问题
  - 使用 `visibility: hidden` 隐藏初始状态
  - 在状态恢复完成后再显示目录树
  - 提升视觉体验
  
- 🎯 优化按钮状态同步
  - "展开/收起"按钮状态与实际状态保持一致
  - 区分"全部展开"和"部分展开"两种状态
  
- 🚀 性能优化
  - 使用标志位防止不必要的状态保存
  - 减少localStorage读写次数
  - 优化事件监听逻辑

### 修复 (Fixed)
- 🐛 修复刷新后自动全部展开的问题
  - 移除模板中硬编码的 `menuSetting = "open"`
  - 改为从localStorage读取，默认为收起状态
  
- 🐛 修复自动选中文档时触发状态保存的问题
  - 使用 `isRestoringState` 标志位
  - 在初始化期间禁用状态保存
  
- 🐛 修复按钮状态与实际状态不一致的问题
  - 修正 `menuControl` 初始值逻辑
  - 根据 `menuSetting` 正确设置按钮显示

### 技术细节 (Technical Details)

#### 修改的文件
```
static/js/markdown.js              - 编辑界面目录树逻辑
static/js/kancloud.js              - 阅读界面目录树逻辑
views/document/default_read.tpl    - 阅读界面模板
```

#### 新增函数
- `processTreeData()` - 处理树数据，设置默认收起
- `saveTreeState()` - 保存展开状态到localStorage
- `restoreTreeState()` - 从localStorage恢复展开状态

#### localStorage键名
- `mindoc_tree_state_<项目ID>` - 编辑界面展开状态
- `mindoc_read_tree_state_<项目ID>` - 阅读界面展开状态
- `mindoc_menu_setting_<项目ID>` - 菜单设置（全部展开/收起）

#### 代码统计
- 新增代码：约150行
- 修改代码：约30行
- 新增函数：6个
- 修改文件：3个

### 兼容性 (Compatibility)
- ✅ 向后兼容：不影响现有功能
- ✅ 浏览器支持：所有支持localStorage的现代浏览器
- ✅ 数据迁移：无需数据迁移，自动适配

### 已知问题 (Known Issues)
- 无

### 升级说明 (Upgrade Notes)
1. 备份当前版本
2. 替换修改的文件
3. 重启MinDoc服务
4. 清除浏览器缓存（Ctrl+F5）
5. 测试功能是否正常

### 贡献者 (Contributors)
- [@RogerXie314](https://github.com/RogerXie314) - 功能开发和实现

---

## [V2.1] - 原始版本

基于 [mindoc-org/mindoc](https://github.com/mindoc-org/mindoc) V2.1 正式版

---

## 版本说明

版本号格式：`主版本号.次版本号.修订号`

- **主版本号**：重大功能变更或不兼容的API修改
- **次版本号**：向后兼容的功能性新增
- **修订号**：向后兼容的问题修正

当前版本：**V2.1.1**
- 基于版本：V2.1
- 修订类型：功能优化和问题修复
- 兼容性：完全向后兼容
