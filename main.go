package main

import (
	"context"
	"fmt"

	"github.com/boni-fm/go-libsd3/pkg/db/postgres"
)

func main() {

	connStr, err := postgres.InitConstrByKodeDc(context.Background(), "dcho")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("constr :: " + connStr)
}
