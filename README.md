# httpx
a lightweight HTTP server

## 安全依赖与 TLS

Fiber 更新至 `v3.4.0`，连带使用 fasthttp `v1.72.0`、compress `v1.19.0` 和 gofiber/utils/v2 `v2.1.1`；x/text 更新至 `v0.41.0`。模块最低 Go 版本仍为 `1.25.0`。这不是全部依赖告警清零声明，消费者仍须核对自己的最终模块图及构建 Go 补丁版本。

`StartWithCertificate` 现在显式提供 `TLSConfig` 和传入的证书，最低 TLS 版本为 1.2。旧实现只提供 `TLSConfigFunc`，没有启用 TLS，回环回归确认旧版与升级后依赖均会退化为普通 HTTP；本次修复该既存问题。测试使用内存生成的回环证书并正常验证证书，要求 HTTPS 握手成功且明文 HTTP 不能正常访问该端口。不读取或保存真实证书。普通 `Start` 仍启动 HTTP。

## Web/API 响应迁移

新 Web/API 代码统一使用 `github.com/coffeehc/commons/webfacade` 的原生响应、查询与分页合同。`AjaxResponse` 保留原有外壳；新的 `ListData[T]` 在 payload 内支持可选总数和明确的续页/末页游标。

`httpxcommons` 中已迁移的响应类型、发送 helper 和响应 protobuf schema 已删除，不保留 alias 或 wrapper。该包只保留未迁移的 `LoadCertificate` 证书加载功能。httpx 不反向依赖 commons。

本次是 Go 源码 API 的破坏性更新：旧 `httpxcommons.AjaxResponse`、`BaseResponse`、`ListData`、`PBResponse` 及 `Send*` 调用统一迁到 `github.com/coffeehc/commons/webfacade`。仅在实际计算总数后传入总数指针，未知总数用 nil；游标模式返回有效令牌或 `EndCursor()`。业务成功只看 `success: true`，不要求某个固定成功码。既有 JSON 外壳和有效 PB 字节由 webfacade 的固定 fixture 验证。

检查：`go test ./...`、`go test -race ./...`、`go vet ./...`。响应与分页合同检查在 commons 的 `go test ./webfacade` 中执行。
