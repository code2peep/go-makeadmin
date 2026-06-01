package main

import (
	"fmt"
	"log"
	"os"

	"go-makeadmin/makeadmin/security"
)

func main() {
	plain := os.Getenv("MAKEADMIN_PASSWORD")
	if plain == "" {
		log.Fatal("MAKEADMIN_PASSWORD is required")
	}
	hasher := security.NewBcryptPasswordHasher(0)
	digest, err := hasher.Hash(plain)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(digest.Hash)
}
