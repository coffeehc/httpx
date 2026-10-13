# httpx
a lightweight HTTP server

## Web/API 响应迁移

新 Web/API 代码统一使用 `github.com/coffeehc/commons/webfacade` 的原生响应、查询与分页合同。`AjaxResponse` 保留原有外壳；新的 `ListData[T]` 在 payload 内支持可选总数和明确的续页/末页游标。

`httpxcommons` 的响应类型与发送 helper 已标记 Deprecated，但为已发布消费者冻结保留。`ListData[T].Total` 仍为 int64，旧源码字面量继续编译；nil List 仍输出 null，Total 为零仍输出。不存在反向依赖或 alias commons；HTTP server、middleware、TLS、protobuf schema 不变。

迁移须显式选择：改用 webfacade，仅在实际计算总数后传入总数指针，未知总数用 nil；游标模式返回有效令牌或 `EndCursor()`。业务成功只看 `success: true`，不要求某个固定成功码。

兼容检查：`go test ./httpxcommons`。
