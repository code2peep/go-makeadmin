# P3 Module Apply Result Empty State

更新时间：2026-06-03

## 目标

P3.24 收敛模块 manifest apply 结果视图的空态文案。

本阶段只调整前端展示，不改变后端，不新增接口，不写数据库。

## 页面空态

`module-manifest-apply-result.vue` 新增空态：

- 权限编码为空时显示 `无权限编码`。
- 没有执行前后快照时显示 `无快照`。
- 没有门禁检查项时显示 `无检查项`。

## 当前边界

P3.24 不改变：

- 按钮状态规则。
- 错误归一化规则。
- 审计预览构造规则。
- 后端接口。
- 数据库写入。
- 菜单权限 SQL。

## 验收

P3.24 需要通过：

```bash
cd admin && npm run type-check
scripts/check-module-tools-no-db.sh
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
