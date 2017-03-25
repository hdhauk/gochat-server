package main

import "github.com/hdhauk/gochat-server/auth"

func main() {
	testUsers := make(map[auth.Username]auth.PassSHA512)
	testUsers["test-user"] = "c638833f69bbfb3c267afa0a74434812436b8f08a81fd263c6be6871de4f1265"
	testConfig := auth.TokenServiceConfig{
		Secret:       []byte("test-secret"),
		PreAuthUsers: testUsers,
	}
	s := auth.NewTokenService(testConfig)
	s.Start()
	select {}

}
