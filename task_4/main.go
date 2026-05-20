package main

import (
	"errors"
	"fmt"
)

var (
	ErrEngineAlreadyRunning = errors.New("двигатель уже работает")
	ErrEngineOff            = errors.New("двигатель не запущен")
	ErrLowBattery           = errors.New("низкий заряд батареи")
)

type Vehicle interface {
	StartEngine() error
	StopEngine() error
	GetInfo() string
}

type Car struct{
	Brand string
	engineOn bool
}

type Truck struct{
	Car
	cargoCapacity float64
}

type ElectricCar struct{
	Car
	batteryLevel int
}

func (c Car) Honk() string{
	fmt.Println("Beep beep!")
	return "Beep beep!"
}

func (t Truck) Honk() string{
	fmt.Println("Honk Honk!")
	return "Honk Honk!"
}

func (c *Car) StartEngine() error{
	if c.engineOn {
		return ErrEngineAlreadyRunning
	}
	c.engineOn=true
	fmt.Println("Двигатель запущен")

	return nil
}

func (c *Car) StopEngine() error{
	if !c.engineOn{
		return ErrEngineOff
	}
	c.engineOn = false
	fmt.Println("Двигатель Остановлен")
	return nil
}

func (c *Car) GetInfo() string{
	str := "Машина "+c.Brand
	return str
}

func (t Truck) GetInfo() string{
	str := "Машина "+t.Brand
	return str
}

func (e ElectricCar) GetInfo() string{
	fmt.Println("Заряд батареи", e.batteryLevel)
	str := "Машина "+ e.Brand 
	return str + "заряд батареи"
}

func (e *ElectricCar) StartEngine() error{
	if e.batteryLevel <= 5 {
		return ErrLowBattery
	}

	return e.Car.StartEngine()
}

func (c *Car) GetEngineStatus() bool {
    return c.engineOn
}

func (e *ElectricCar) GetBatteryLevel() int {
    return e.batteryLevel
}

func (t *Truck) GetCargoCapacity() float64{
	return t.cargoCapacity
}