package retry_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/boni-fm/go-libsd3/pkg/retry"
)

func ExampleDo() {
	attempts := 0
	err := retry.Do(context.Background(), retry.DefaultConfig(), func() error {
		attempts++
		if attempts < 2 {
			return errors.New("temporary failure")
		}
		return nil
	})
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleDoWithResult() {
	result, err := retry.DoWithResult(context.Background(), retry.DefaultConfig(), func() (string, error) {
		return "berhasil", nil
	})
	fmt.Println(result, err)
	// Output:
	// berhasil <nil>
}
