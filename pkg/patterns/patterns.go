package patterns

import (
	"fmt"
	"math"
	"time"
)

func Retry(fn func() error, maxRetries int, baseDelay int) error {
	for i := 1; i <= maxRetries; i++ {
		err := fn()
		if err == nil {
			return nil
		}
		if i != maxRetries {
			delay := baseDelay * int(math.Pow(2, float64(i)))
			// fmt.Printf("Повтор #%d. Ошибка: %v. Повтор через %d секунд\n", i, err, delay)
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}
	}
	return fmt.Errorf("Не удалось выполнить функцию за %d попыток", maxRetries)
}

func Timeout(fn func() error, timeout int) error {
	ch := make(chan error, 1)
	go func() {
		ch <- fn()
	}()
	select {
	case err := <-ch:
		return err
	case <-time.After(time.Duration(timeout) * time.Millisecond):
		return fmt.Errorf("Функция не смогла завершить сворю работу за %d миллисекунд", timeout)
	}
}

type DeadLetterQueue struct {
	messages chan string
}

func NewDeadLetterQueue() *DeadLetterQueue {
	return &DeadLetterQueue{messages: make(chan string, 1000)}
}

func (dlq *DeadLetterQueue) Add(msg string) {
	dlq.messages <- msg
}

func (dlq *DeadLetterQueue) GetMessages() []string {
	close(dlq.messages)
	var messages []string
	for msg := range dlq.messages {
		messages = append(messages, msg)
	}
	return messages
}

func ProcessWithDLQ(msgs []string, fn func(string) error, dlq *DeadLetterQueue) {
	for _, msg := range msgs {
		err := fn(msg)
		if err != nil {
			dlq.Add(msg)
		}
	}
}

// func errorFunction() error {
// 	num := rand.Intn(100)
// 	if num == 1 {
// 		fmt.Println("Успех!")
// 		return nil
// 	}
// 	fmt.Println("Неудача!")
// 	return fmt.Errorf("Ошибка! num = %d", num)
// }

// func errorFunctionTimeout() error {
// 	time.Sleep(4000 * time.Millisecond)
// 	fmt.Println("Успех!")
// 	return nil
// }
