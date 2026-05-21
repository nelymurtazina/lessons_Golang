package main

import (
	"errors"
	"fmt"
	"math/rand"
)

var (
	ErrInvalidAmount       = errors.New("некорректная сумма платежа")
	ErrProviderUnavailable = errors.New("провайдер недоступен")
)

type PaymentProcessor interface {
	ProcessPayment(amount float64) error
}

type Sberbank struct {
	APIKey string
}
type Tbank struct {
	APIKey string
}
type Alfabank struct {
	APIKey string
}

func (bank Sberbank) ProcessPayment(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	shans := rand.Intn(100)
	if shans < 25 {
		return ErrProviderUnavailable
	}
	if amount > 0 {
		return nil
	}
	return nil
}
func (bank Tbank) ProcessPayment(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	shans := rand.Intn(2)
	if shans < 25 {
		return ErrProviderUnavailable
	}
	if amount > 0 {
		return nil
	}
	return nil
}
func (bank Alfabank) ProcessPayment(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	shans := rand.Intn(100)
	if shans < 25 {
		return ErrProviderUnavailable
	}
	if amount > 0 {
		return nil
	}
	return nil
}

func main() {
	sber := Sberbank{
		APIKey: "sber",
	}
	tbank := Tbank{
		APIKey: "tbank",
	}
	alfa := Alfabank{
		APIKey: "alfa",
	}

	one := sber.ProcessPayment(100)
	fmt.Println(one)
	two := tbank.ProcessPayment(-100)
	fmt.Println(two)
	tree := alfa.ProcessPayment(0)
	fmt.Println(tree)

}
