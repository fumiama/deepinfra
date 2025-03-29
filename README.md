# deepinfra
Lightweight OpenAI/OLLaMA/GenAI compatible caller, originally designed for DeepInfra.

## Quick Start
```go
api := NewAPI(OpenAIDeepInfra, "PUT YOUR API KEY HERE")
txt, err := api.Request(model.NewOpenAI(model.ModelDeepDeek, model.SeparatorThink, 0.7, 0.9, 1024).
    System("Be a good assistant.").User("Hello"),
)
if err != nil {
    panic(err)
}
fmt.Println(txt)
// Hello! How can I assist you today?
```
