# 项目测试与初始化上下文提醒（context_reminder.md）

## 1. gRPC 连接统一规范

请统一使用新版 API：
```go
func NewClient(addr string) (*grpc.ClientConn, error) {
    return grpc.NewClient(
        addr,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
}
```
> 不要再用 `grpc.Dial`！

---

## 2. 测试配置初始化方法

必须存在如下方法，并能正确加载 config/config-develop.yaml：
```go
func initTestConfig() error {
    // ...（加载配置，使用 viper 解析到 global.ServerConfig）
}
```

---

## 3. TestMain 规范

- `TestMain` 只允许在 `init_test.go` 中唯一声明。
- 负责所有测试初始化（日志、配置、数据库、gRPC连接等）。
- 其他测试文件无需再声明 `TestMain`。

---

## 4. 全局变量复用

- `goodsClient`、`conn` 只在 `init_test.go` 中声明和初始化。
- 其他测试文件直接复用，无需重复声明。

---

## 5. 变更同步

如需修改初始化逻辑，务必同步更新本文件和所有相关测试文件。

---

## 6. AI/团队协作提醒

- 每次有测试相关问题或需求时，请主动关联本文件内容。
- AI 回答时请严格遵循本文件约定，避免重复犯错。

---

> 本文件为项目测试和初始化相关的通用约定，适用于所有团队成员和AI助手。 