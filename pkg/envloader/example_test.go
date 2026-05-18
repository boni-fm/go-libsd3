package envloader_test

import (
	"fmt"
	"os"

	"github.com/boni-fm/go-libsd3/pkg/envloader"
)

func ExampleLoad() {
	f, _ := os.CreateTemp("", "example*.env")
	f.WriteString("APP_NAME=myapp\nAPP_PORT=8080\n")
	f.Close()
	defer os.Remove(f.Name())

	os.Unsetenv("APP_NAME")
	os.Unsetenv("APP_PORT")

	if err := envloader.Load(f.Name()); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(os.Getenv("APP_NAME"))
	fmt.Println(os.Getenv("APP_PORT"))
	// Output:
	// myapp
	// 8080
}
