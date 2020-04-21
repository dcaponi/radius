package main

import (
	"fmt"
	"log"
	"os"

	"github.com/onelogin/onelogin-go-sdk/pkg/client"
	"github.com/onelogin/onelogin-go-sdk/pkg/models"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
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

	handler := func(w radius.ResponseWriter, r *radius.Request) {
		username := rfc2865.UserName_GetString(r.Packet)
		password, er := rfc2865.UserPassword_LookupString(r.Packet)
		if er != nil {
			fmt.Println("ERROR")
			fmt.Println(er)
		}
		fmt.Printf("UN %s PW %s:\n", username, password)

		attrs := r.Packet.Attributes
		for k, v := range attrs {
			fmt.Printf("K: %d, V: %s\n", k, string(v[0]))
		}

		subdomain := os.Getenv("OL_SUBDOMAIN")

		request := &models.SessionLoginTokenRequest{
			UsernameOrEmail: oltypes.String(username),
			Password:        oltypes.String(password),
			Subdomain:       oltypes.String(subdomain),
		}

		var code radius.Code

		resp, _, err := oneloginClient.Services.SessionLoginTokensV1.CreateSessionLoginToken(request)
		if err != nil {
			log.Println(err)
			code = radius.CodeAccessReject
		}

		if resp.StatusCode == 200 {
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
		SecretSource: radius.StaticSecretSource([]byte(`testing123`)),
	}

	log.Printf("Starting server on :1812")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
