# LikeAdmin Demo Modules

This package keeps LikeAdmin business demo modules outside the core `go-makeadmin` runtime source tree.

Moved modules:

- `admin/src/views/article`
- `admin/src/views/consumer`
- `admin/src/views/channel`
- `admin/src/views/decoration`
- `admin/src/views/message`
- `admin/src/views/setting/search`
- `admin/src/views/setting/user`
- `admin/src/api/article.ts`
- `admin/src/api/consumer.ts`
- `admin/src/api/channel`
- `admin/src/api/decoration.ts`
- `admin/src/api/message.ts`
- `admin/src/api/setting/search.ts`
- `admin/src/api/setting/user.ts`
- `admin/src/components/link`

These files are reference material only. They are not part of the P0 core admin framework, are not loaded by dynamic routes, and should not be imported by production code.
