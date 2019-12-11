# ngrok-go

This client simplifies interaction with the ngrok API.

Example usage:

```go
import ngrok "github.com/meterup/ngrok-go"

client := ngrok.New(ngrok.BaseURL, os.Getenv("NGROK_API_TOKEN"))
data := url.Values{}
data.Set("description", "ngrok-go")
cred, err := client.Creds.Create(context.TODO(), data)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("created credential: %#v\n", cred)
```

The client does not implement the entire API, but it should be easy to add new
endpoints following the existing pattern.
