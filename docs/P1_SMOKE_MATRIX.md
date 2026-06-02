# P1 Smoke Matrix

## 目标

P1.21 建立一套面向 P1 独立库的 HTTP smoke 验证，覆盖已切换到 `ma_*` 的核心后台链路。该 smoke 只应运行在本地一次性验证库和本地上传目录上。

## 脚本

```bash
python3 scripts/p1-smoke.py --print-matrix
```

写操作默认被保护。真正执行前需要先启动 API，并显式允许写入：

```bash
P1_SMOKE_ALLOW_WRITE=1 \
P1_SMOKE_BASE_URL=http://127.0.0.1:8000/api \
P1_SMOKE_ADMIN_PASSWORD='your-local-admin-password' \
python3 scripts/p1-smoke.py
```

## 运行前提

- 已用 `scripts/init-p1-db.sh` 初始化本地 disposable P1 数据库。
- API 服务指向该 disposable P1 数据库。
- `admin` 账号密码通过 `P1_SMOKE_ADMIN_PASSWORD` 或 `ADMIN_PASSWORD` 传入。
- 不要对共享库、生产库或保留数据的开发库运行。

## 覆盖矩阵

| 模块 | 接口 | 类型 | 断言 |
| --- | --- | --- | --- |
| auth | `POST /system/login` | read | 登录返回 token |
| auth | `JWT claims` | read | 登录 token 包含 `sid/adminId/tenantId/iat/exp/iss` |
| auth | `X-Tenant-ID mismatch` | read | token 请求拒绝不匹配租户头 |
| auth | `GET /system/tenant/list` | read | 租户列表返回默认租户 |
| auth | `POST /system/tenant/switch` | write | 切换到默认租户返回新 token |
| auth | `POST /system/logout` | write | JWT session state 可删除 |
| auth | `GET /system/admin/self` | read | token 能解析当前管理员 |
| auth | `GET /system/menu/route` | read | token 能解析菜单路由 |
| common | `GET /common/index/config` | read | 公共配置从 `ma_setting` 返回 |
| common | `GET /common/index/console` | read | 控制台信息从 `ma_setting` 返回 |
| log | `GET /system/log/login` | read | 登录日志可查询 |
| role | `POST /system/role/add` | write | 角色可新增 |
| log | `GET /system/log/operate` | read | 写操作后操作日志可查询 |
| role | `POST /system/role/edit` | write | 角色可编辑 |
| role | `POST /system/role/del` | write | 管理员清理后角色可删除 |
| admin | `POST /system/admin/add` | write | 管理员可新增 |
| admin | `POST /system/admin/edit` | write | 管理员可编辑 |
| admin | `POST /system/admin/disable` | write | 管理员状态可切换 |
| admin | `POST /system/admin/del` | write | 管理员可删除 |
| menu | `POST /system/menu/add` | write | 菜单按钮可新增 |
| menu | `POST /system/menu/edit` | write | 菜单按钮可编辑 |
| menu | `POST /system/menu/del` | write | 菜单按钮可删除 |
| dict | `POST /setting/dict/type/add` | write | 字典类型可新增 |
| dict | `POST /setting/dict/data/add` | write | 字典项可新增 |
| dict | `POST /setting/dict/data/edit` | write | 字典项可编辑 |
| dict | `POST /setting/dict/data/del` | write | 字典项可删除 |
| dict | `POST /setting/dict/type/del` | write | 字典类型可删除 |
| file | `POST /common/album/cateAdd` | write | 文件分类可新增 |
| file | `POST /common/upload/image` | write | 图片上传可写入 `ma_file` 元数据 |
| file | `POST /common/album/albumDel` | write | 上传文件元数据可删除 |
| file | `POST /common/album/cateDel` | write | 文件分类可删除 |
| codegen | `POST /gen/importTable` | write | 表元数据可导入 |
| codegen | `POST /gen/syncTable` | write | 表元数据可同步 |
| codegen | `GET /gen/previewCode` | read | 表元数据可渲染预览 |
| codegen | `GET /gen/downloadCode` | read | 表元数据可渲染 zip |
| codegen | `POST /gen/delTable` | write | 表元数据可删除 |

## 清理策略

- 脚本会按反序清理已创建的管理员、角色、菜单、字典、文件元数据、文件分类和代码生成配置。
- 上传 smoke 会写入一个 1x1 PNG；脚本会清理 `ma_file` 元数据，但物理文件可能仍留在本地上传目录。
- 如果中途失败，脚本仍会尽最大努力执行已注册的清理步骤。

## 不覆盖范围

- 不覆盖旧 `la_*` 兜底链路。
- 不验证生产部署、远程存储驱动、真实业务表迁移。
- 不修改数据库 schema。
