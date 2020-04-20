package main

import (
	"fmt"
	"log"
	"os"

	"github.com/onelogin/onelogin-go-sdk/pkg/client"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

func main() {
	oneloginClient, err := client.NewClient(&client.APIClientConfig{
		Timeout:      5,
		ClientID:     os.Getenv("OL_CLIENT_ID"),
		ClientSecret: os.Getenv("OL_CLIENT_SECRET"),
		Url:          os.Getenv("OL_ENDPOINT"),
	})
	if err != nil {
		fmt.Println(err)
	}
	resp, auth, err := oneloginClient.Services.AuthV2.Authorize()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
	fmt.Println(*auth.AccessToken)

	handler := func(w radius.ResponseWriter, r *radius.Request) {
		username := rfc2865.UserName_GetString(r.Packet)
		password := rfc2865.UserPassword_GetString(r.Packet)

		var code radius.Code
		// authenticate username and password with ol. make this a goroutine (new thread) call
		// read response from channel when goroutine finishes
		if username == "tim" && password == "12345" {
			code = radius.CodeAccessAccept
		} else {
			code = radius.CodeAccessReject
		}
		// handle challenge here?
		log.Printf("Writing %v to %v", code, r.RemoteAddr)
		w.Write(r.Response(code))
	}

	server := radius.PacketServer{
		Handler:      radius.HandlerFunc(handler),
		SecretSource: radius.StaticSecretSource([]byte(`secret`)),
	}

	log.Printf("Starting server on :1812")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
