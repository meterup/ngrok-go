package ngrok_test

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	ngrok "github.com/meterup/ngrok-go"
)

func Example() {
	client := ngrok.New(ngrok.BaseURL, os.Getenv("NGROK_API_TOKEN"))
	data := url.Values{}
	data.Set("description", "ngrok-go")
	cred, err := client.Creds.Create(context.TODO(), data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("created credential: %#v\n", cred)
}
