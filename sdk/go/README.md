# Habarchy SDK — Go

```sh
go get github.com/Esca6585dev/habarchy/sdk/go
```

```go
c := habarchy.New("https://habarchy.example.tm", os.Getenv("HABARCHY_API_KEY"), habarchy.WithSigning())
acc, err := c.SendMessage(ctx, habarchy.SendMessageRequest{
    Channel: "sms", To: habarchy.To("+99365123456"), Template: "otp",
    Data: map[string]any{"code": "4821", "minutes": 5}, IdempotencyKey: "order-77",
})
// acc.ID, acc.Status == "queued"; later: c.GetMessage(ctx, acc.ID)
```

Methods: `SendMessage`, `SendBatch`, `GetMessage`, `CancelMessage`, `GetBatch`, `SendOTP`,
`VerifyOTP`, `RegisterDevice`. Recipients: `habarchy.To("+993…")`,
`habarchy.ToExternal("user-42")`, `habarchy.ToContact(id)`; channel `"auto"` picks the best
channel for a contact.

Errors are `*habarchy.Error{Status, Code, Message, Details}`; `habarchy.IsCode(err, "quota_exceeded")`.

`WithSigning()` adds `X-Timestamp` and `X-Signature = hex(HMAC-SHA256(apiKey,
ts "\n" METHOD "\n" path "\n" hex(sha256(body))))`, required for keys created with
`require_signature`. The clock must be within ±5 minutes of the server.

```sh
go test ./...
```
