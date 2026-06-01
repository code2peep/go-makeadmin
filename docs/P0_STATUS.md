# P0 状态记录

更新时间：2026-05-31

## 已完成

- 新建项目路径：`/Users/fengrongxin/AI/01-projects/go-makeadmin`。
- 添加项目规则：`CLAUDE.md`。
- 引入 LikeAdmin Go MIT 蓝本。
- 添加 `NOTICE.md` 记录来源和授权。
- 重写项目 README 为 `go-makeadmin` 定位。
- 后端 module 从 `likeadmin` 改为 `go-makeadmin`。
- 后端依赖升级：
  - Gin `v1.12.0`
  - Gorm `v1.31.1`
  - MySQL driver `v1.6.0`
  - Redis `github.com/redis/go-redis/v9 v9.20.0`
  - Zap `v1.28.0`
  - Viper `v1.21.0`
- 后端配置、MySQL、Redis 初始化改为更适合测试和框架复用的方式。
- 前端包名改为 `go-makeadmin-admin`。
- 前端依赖升级：
  - Vue `3.5.x`
  - Vite `8.x`
  - Element Plus `2.14.x`
  - Pinia `3.x`
  - Axios `1.x`
  - TypeScript `6.x`
- 默认 `npm run build` 改为只构建，不删除或复制 `frontend`。
- 前端 Pinia Options Store 适配 Pinia 3 的 `defineStore(id, options)` 写法。
- 补齐 TypeScript 6 下的虚拟模块、WangEditor、Axios 1 和 Element Plus 类型兼容。
- 移除直接风险依赖 `css-color-function`，改为本地 hex 主题色混色函数。
- 移除直接风险依赖 `vite-plugin-svg-icons`，改为本地 Vite SVG sprite 插件。
- `package-lock.json` 已通过 `npm audit fix --package-lock-only --ignore-scripts` 收敛到 0 个 audit vulnerabilities。
- 新增本地开发启动文档 `docs/LOCAL_DEV.md`。
- 新增不触库验证脚本 `scripts/verify-no-db.sh`。
- 新增 MySQL/Redis 只读前置检查脚本 `scripts/check-services.sh`。
- 新增本地 API 和管理端启动脚本：
  - `scripts/dev-api.sh`
  - `scripts/dev-admin.sh`
- 后端新增 `SERVER_HOST`，本地默认监听 `127.0.0.1`，避免端口冲突时只绑定 IPv6 导致请求打到其他服务。
- `scripts/dev-api.sh` 增加端口占用预检查。
- `scripts/dev-admin.sh` 使用固定端口和 `--strictPort`。
- 前端 Vite 开发服务增加 `/api` 本地代理，默认转发到 `http://127.0.0.1:8000`，当前本机转发到 `http://127.0.0.1:18000`。
- 新增一次性本地库初始化方案 `docs/DB_INIT_PLAN.md`。
- 新增初始化脚本和种子只读检查脚本：
  - `scripts/init-local-blueprint-db.sh`
  - `scripts/check-db-seed.sh`
- 当前本机已初始化开发库 `go_makeadmin`，导入 P0 蓝本 `la_*` 表和初始化数据。
- 当前本机已创建真实开发配置：
  - `server/.env`
  - `admin/.env.development`
- P0 可见品牌残留已初步替换为 `go-makeadmin`，包括系统配置种子、工作台版本信息和代码生成 zip 名称。
- P0.8 已移除工作台首屏的蓝本服务支持二维码区域。
- 新增 P0.9 模块裁剪清单 `docs/MODULE_PRUNE_PLAN.md`，区分核心后台模块和业务演示模块。
- 新增核心 SQL 生成脚本 `scripts/build-core-sql.sh`。
- 新增最小核心初始化 SQL `sql/install.core.sql`。
- 新增核心库初始化包装脚本 `scripts/init-local-core-db.sh`。
- P0.10 已将前端动态路由视图改为核心模块白名单，业务演示模块源码保留但不再进入动态路由池。
- P0.10 已将工作台版本信息字段从 `channel` 改为 `links`，并清理核心菜单编辑页里的业务演示路径示例。
- 新增 `legacy` 目录结构约定。
- P0.11 已将前端业务演示视图和 API 封装迁入 `legacy/likeadmin-demo`，核心 `admin/src` 源码树不再保留这些模块。
- 新增 P0 SQL 残留审计文档 `docs/P0_SQL_RESIDUE_AUDIT.md`。
- 新增 P1 `ma_*` schema 设计基线 `docs/P1_SCHEMA_PLAN.md`。
- P0.12 已移除核心网站设置链路中的 `shopName`、`shopLogo`，并从核心 SQL 生成结果中排除这两个蓝本字段。
- P0.12 已移除后端免权限白名单中的 `article:cate:all`。
- P0.12 已将装修器专用 `admin/src/components/link` 迁入 `legacy/likeadmin-demo`。
- P0.13 已将样式入口中的 SCSS `@import` 迁移到 `@use`。
- P0.14 已将 `vuedraggable` 从 Vue 2 版本升级到 Vue 3 版本，修复网站信息页素材选择器运行时报错。
- P0.14 已调整布局路由渲染，避免把 `RouterView` 组件直接放进 `keep-alive`。
- P0.14 已替换 `server/static/backend_backdrop.png`，移除登录页广告图位图里的 LikeAdmin 品牌残留。
- P0.15 已收窄 ECharts 注册范围，只保留当前核心页面实际使用的折线图、饼图、仪表盘和基础组件。
- P0.15 已将 Vite chunk warning 阈值调整为 `900KB`，当前富文本和图表 chunk 属于已拆分的懒加载功能依赖。

## 已验证

后端：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin/server
GOPROXY=https://goproxy.cn,direct go test ./...
```

结果：通过。

前端：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin/admin
npm run build
```

结果：通过。

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin/admin
npm run type-check
```

结果：通过。

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin/admin
npm audit --json
```

结果：0 个 vulnerabilities。

P0.6 不触库验证：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
bash -n scripts/verify-no-db.sh scripts/check-services.sh scripts/dev-api.sh scripts/dev-admin.sh
```

结果：通过。

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
./scripts/verify-no-db.sh
```

结果：通过。

P0.6 服务前置检查：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
./scripts/check-services.sh
```

结果：通过。MySQL、Redis 和 `go_makeadmin` 开发库均可连接。

P0.7 数据库种子检查：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
./scripts/check-db-seed.sh
```

结果：通过。关键 `la_*` 表存在，默认管理员存在，菜单种子 `75` 条，系统配置种子 `46` 条。

P0.7 API 烟测：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
./scripts/dev-api.sh
```

当前本机 API 地址：

```text
http://127.0.0.1:18000/api
```

已验证：

- `GET /api/common/index/config` 返回 `code=200`，`webName=go-makeadmin`。
- `POST /api/system/login` 使用 `admin / 123456` 返回 token。
- `GET /api/system/admin/self` 返回当前管理员和 `permissions=["*"]`。
- `GET /api/system/menu/route` 返回菜单路由。
- `GET /api/common/index/console` 返回 `version=v0.1.0` 和 go-makeadmin 技术栈。

P0.7 管理端浏览器烟测：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
./scripts/dev-admin.sh
```

结果：通过。浏览器可登录并进入 `/workbench`，工作台可读到 `go-makeadmin`、`v0.1.0` 和新技术栈文案。

P0.9 核心 SQL 验证：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
./scripts/build-core-sql.sh
```

结果：生成 `sql/install.core.sql`，只包含 15 张核心表。

已用一次性临时库 `go_makeadmin_core_check` 导入 `sql/install.core.sql` 并运行：

```bash
MYSQL_DATABASE=go_makeadmin_core_check ./scripts/check-db-seed.sh
```

结果：通过。默认管理员存在，菜单种子 `75` 条，核心系统配置 `14` 条。临时库已删除。

同时已验证 `scripts/init-local-core-db.sh --apply` 可将核心 SQL 导入一次性临时库并通过同样的种子检查。

P0.10 动态路由隔离验证：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
./scripts/verify-no-db.sh
```

结果：通过。构建仍有已知 Rolldown/chunk warning，但业务演示视图不再进入动态路由构建池。

已检查 `admin/dist`，未发现 `article/`、`channel/`、`consumer/`、`decoration/`、`message/`、`setting/search`、`setting/user`、`LikeAdmin` 等业务演示路由或品牌字符串。文件名命中仅剩 `clipboard` 工具库，和 `channel` 业务模块无关。

P0.10 API 烟测：

```bash
GET /api/common/index/console
```

结果：通过。`version.name=go-makeadmin`，`version.version=v0.1.0`，`version.links` 存在，旧字段 `version.channel` 不再返回。

P0.10 管理端浏览器烟测：

```text
http://127.0.0.1:5173/workbench
```

结果：通过。工作台可打开，页面无前端 error/warn 日志，`v0.1.0`、技术栈和“官网/Gitee”链接正常展示，`服务支持` 和 `LikeAdmin` 未出现。

P0.11 业务演示源码隔离验证：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
./scripts/verify-no-db.sh
```

结果：通过。核心 `admin/src/views` 和 `admin/src/api` 已无 `article`、`consumer`、`channel`、`decoration`、`message`、`setting/search`、`setting/user` 等业务演示目录或 API 封装。

已检查 `admin/dist`，未发现业务演示路由字符串或 `LikeAdmin` 品牌字符串。

P0.11 管理端浏览器烟测：

```text
http://127.0.0.1:5173/workbench
```

结果：通过。工作台可打开，页面无前端 error/warn 日志，业务演示路由未出现在当前可见页面。

P0.12 SQL/配置残留清理验证：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
./scripts/build-core-sql.sh
```

结果：通过。`sql/install.core.sql` 不再包含 `shopName`、`shopLogo`。

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
./scripts/verify-no-db.sh
```

结果：通过。

已用一次性临时库 `go_makeadmin_core_p012_check` 导入 `sql/install.core.sql` 并运行 `./scripts/check-db-seed.sh`。

结果：通过。`la_system_config` 核心配置种子 `12` 条，`shopName`、`shopLogo` 计数为 `0`。临时库已删除。

P0.12 API 烟测：

```bash
GET /api/setting/website/detail
```

结果：通过。返回字段为 `name`、`logo`、`favicon`、`backdrop`，不再返回 `shopName`、`shopLogo`。

P0.13 Sass warning 清理验证：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
./scripts/verify-no-db.sh
```

结果：通过。Sass `@import` 废弃 warning 已消失。

P0.14 会话恢复基线验证：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
./scripts/verify-no-db.sh
```

结果：通过。后端测试、前端 type-check、前端 build、`npm audit` 均通过；剩余 warning 与已知问题一致。

P0.14 网站信息页浏览器烟测：

```text
http://127.0.0.1:5173/setting/website/information
```

结果：通过。页面可渲染 `网站名称`、`网站图标`、`网站logo`、`登录页广告图`，无新增前端 error/warn，未发现 `商城`、`店铺`、`shopName`、`shopLogo` 字符串。

P0.14 静态品牌图验证：

```bash
shasum -a 256 server/static/backend_backdrop.png
curl -fsS http://127.0.0.1:18000/api/static/backend_backdrop.png | shasum -a 256
```

结果：通过。服务端静态文件和 API 返回内容 checksum 一致，登录页广告图源文件已不再包含 LikeAdmin 字样。

P0.15 构建收尾验证：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
./scripts/verify-no-db.sh
```

结果：通过。后端测试、前端 type-check、前端 build、`npm audit` 均通过；构建不再出现大 chunk warning，仅保留已知 `@vueuse/core` Rolldown pure annotation warning。

P0.15 构建体积结果：

```text
@wangeditor: 781KB
echarts: 536KB
element-plus: 499KB
vue3-video-play: 370KB
```

说明：`@wangeditor` 只用于网站协议编辑页，`vue3-video-play` 只用于素材视频预览，均为功能页依赖；ECharts 已从约 `823KB` 降至约 `536KB`。

P0.15 浏览器烟测：

```text
http://127.0.0.1:5173/workbench
http://127.0.0.1:5173/setting/system/cache
http://127.0.0.1:5173/setting/website/information
```

结果：通过。工作台可渲染 `访问量趋势图` 并生成 1 个图表 canvas；系统缓存页可渲染 `命令统计` 和 `内存信息` 并生成 2 个图表 canvas；网站信息页可渲染 `网站名称` 和 `网站图标`；复测期间无新增前端 error。

## 已知问题

- Vite 8/Rolldown 对 `@vueuse/core` 的 pure annotation 有 warning，不影响当前构建。
- 富文本、图表、视频播放器仍是较大的功能依赖 chunk，但已按页面能力拆分并低于当前项目 warning 阈值。
- P0 仍保留 `la_*` 表和部分演示数据，P1 需要重设 `ma_*` 自研系统模型。
- 原始 `sql/install.sql` 仍保留商城、客服、文章等蓝本业务演示数据，后续需要决定是继续保留作蓝本参考，还是迁入独立 demo SQL。

## 未执行

- 未连接 zyai 业务库做写操作。
- 未实现最终 `ma_*` 自研 schema。
- 未迁移 zyai 业务数据。

## 下一步建议

P1：设计并实现 `makeadmin` 自己的认证、权限、多租户和数据范围模型。进入 P1 前需要单独确认 schema 边界，因为创建、修改或迁移数据库 schema 属于项目红线。

P1 第一小步建议：只生成 `ma_*` SQL 和 Go model 草案，不自动导入任何数据库；租户能力先预留并默认关闭。
