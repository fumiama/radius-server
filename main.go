package main

import (
	"log"
	"strings"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

func main() {
	db, err := NewUsersDB("radius-server-users.db")
	if err != nil {
		panic(err)
	}

	err = db.Add("default", "test", "test")
	if err != nil {
		panic(err)
	}

	handler := func(w radius.ResponseWriter, r *radius.Request) {
		username := rfc2865.UserName_GetString(r.Packet)
		password := rfc2865.UserPassword_GetString(r.Packet)
		location := "default"
		i := strings.LastIndex(username, "@")
		if i > 0 {
			location = username[i+1:]
			username = username[:i]
		}
		code := radius.CodeAccessReject
		pwd, err := db.Password(location, username)
		if err == nil && password == pwd {
			code = radius.CodeAccessAccept
		}
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
