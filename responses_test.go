package ngrok

var credsCreateResponse = []byte(`
{
    "acl": [],
    "created_at": "2019-12-10T21:25:10Z",
    "description": "go:test",
    "id": "cr_1UoBXvrhYlYA9um85RKLjAGabBJ",
    "metadata": "go:test",
    "token": "1UoBXvrhYlYA9um85RKLjAGabBJ_26F8RrQveit7vjo1oRddF",
    "uri": "https://api.ngrok.com/credentials/cr_1UoBXvrhYlYA9um85RKLjAGabBJ"
}
`)

var addrsCreateResponse = []byte(`
{
    "addr": "1.tcp.ngrok.io:28712",
    "created_at": "2019-12-10T23:35:14Z",
    "description": "go:test",
    "endpoint_configuration": null,
    "id": "ra_3GCBoX98pv2bBjWqS9AfL",
    "metadata": "go:test",
    "region": "us",
    "uri": "https://api.ngrok.com/reserved_addrs/ra_3GCBoX98pv2bBjWqS9AfL"
}
`)
