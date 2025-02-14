# deepinfra
Call OpenAI compatible APIs, originally designed for DeepInfra.

## Quick Start
```go
api := NewAPI(APIDeepInfra, "PUT YOUR API KEY HERE")
txt, err := api.Request(model.NewDeepSeek(0.7, 0.9, 1024).
    System("Be a good assistant.").User("Hello"),
)
if err != nil {
    panic(err)
}
fmt.Println(txt)
// Hello! How can I assist you today?
```

## Custom Call
```go
api := NewAPI(APIDeepInfra, "PUT YOUR API KEY HERE")
txt, err := api.Request(model.NewCustom("fumiama/ninus", "", 0.7, 0.9, 1024).
    System("你正在QQ群与用户聊天，用户发送了消息。按自己的心情简短思考后，条理清晰地回应**一句话**，禁止回应多句。").
    User("总不能什么都查吧").User("后面DOGE就成恶龙了 很常见的场景"),
)
if err != nil {
    panic(err)
}
fmt.Println(txt)
// 要不我给你查一下？
```
