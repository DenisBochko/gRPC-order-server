package patterns

import "fmt"

func Retry(fn func() error, maxRetries int, baseDelay int) error {
	x := fn()
	return x
}

func errorFunction() error {
	return fmt.Errorf("Ошибка!")
}