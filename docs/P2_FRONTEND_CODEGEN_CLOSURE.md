# P2 Frontend Codegen Closure

更新时间：2026-06-02

## 目标

P2.10 把代码生成器的前端输出纳入可验证闭环，确认生成的 API 和 Vue 页面模板符合当前 `admin` 工程约定。

## 当前落地

- 新增 `scripts/check-codegen-frontend.sh`。
- 新增 env-gated 测试 `TestGeneratedCrudFrontendCodeTypeChecks`。
- 测试默认跳过，只有 `MAKEADMIN_CODEGEN_FRONTEND_CHECK=1` 时执行。
- 测试会渲染 `vue/api.ts.tpl`、`vue/index.vue.tpl`、`vue/edit.vue.tpl`。
- 测试临时写入：
  - `admin/src/api/article.ts`
  - `admin/src/views/article/index.vue`
  - `admin/src/views/article/edit.vue`
- 测试执行 `npm run type-check`，验证生成代码与当前 admin TypeScript/Vue 约定兼容。
- 测试结束后清理临时生成文件和目录。

## 不在 P2.10 做

- 不提交生成后的 `admin/src/api/article.ts`。
- 不提交生成后的 `admin/src/views/article/`。
- 不注册前端菜单。
- 不写入权限种子。
- 不改 admin 全局自动导入配置。

## 验证

P2.10 需要通过：

```bash
./scripts/check-codegen-frontend.sh
cd admin && npm run type-check
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
