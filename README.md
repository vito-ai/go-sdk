# RTZR STT SDK for golang

Go packages for [RTZR STT API](https://developers.rtzr.ai/docs/) service.

# Authorization 

By default, each API will use [RTZR CLIENT SECRET](https://developers.rtzr.ai/docs/authentications) for authorization. 

To make authorize, you should export variables,
``` bash
export RTZR_CLIENT_ID="YOUR_CLIENT_ID"
export RTZR_CLIENT_SECRET="YOUR_CLIENT_SECRET"
```

This will allow your application to run without requiring configuration
``` go
client, err := speech.NewRestClient(nil)
```

or you can manually authorize in your code,
``` go
client, err := speech.NewRestClient(&option.ClientOption{
    ClientId : "YOUR_CLIENT_ID",
    ClientSecret : "YOUR_CLIENT_SECRET",
})
```

# Examples

you can see examples of using RTZR STT SDK.
[rtzr-go-tutorial](https://github.com/vito-ai/go-tutorial)
