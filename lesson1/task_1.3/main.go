package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Device interface {
	UpdateOs(version string) error
	GetInfo() string
}
type Smartphone struct {
	OSVersion string
	Model     string
}
type Laptop struct {
	OSVersion string
	Model     string
}
type Smartwatch struct {
	OSVersion string
	Model     string
}

var (
	ErrUnsupported = errors.New("обновление недоступно")
)

func (phone *Smartphone) UpdateOS(version string) error {
	num, err := strconv.Atoi(version)
	if err != nil {
		fmt.Println("Error", err)
		return err
	}
	if num <= 12.0 {
		fmt.Println(ErrUnsupported)
		return ErrUnsupported
	}
	phone.OSVersion = version
	fmt.Println("Model: " + phone.Model + "," + " OC: " + phone.OSVersion)
	return nil
}

func (phone *Smartphone) GetInfo() string {
	str := "Model: " + phone.Model + "," + " OC: " + phone.OSVersion
	fmt.Println(str)
	return str
}

func (lap *Laptop) UpdateOS(version string) error {
	str := "Windows"
	new := strings.Split(version, " ")
	if str == new[0] {
		lap.OSVersion = new[1]
		fmt.Println(version)
	}
	if str != new[0] {
		fmt.Println(ErrUnsupported)
		return ErrUnsupported
	}
	return nil
}

func (lap *Laptop) GetInfo() string {
	str := "Model: " + lap.Model + "," + " OC: " + lap.OSVersion
	fmt.Println(str)
	return str
}

func (whatch *Smartwatch) GetInfo() string {
	str := "Model: " + whatch.Model + "," + " OC: " + whatch.OSVersion
	fmt.Println(str)
	return str
}

func (whatch *Smartwatch) UpdateOS(version string) error {
	new := strings.Split(version, ".")
	if len(new) > 5 {
		fmt.Println(ErrUnsupported)
		return ErrUnsupported
	}
	whatch.OSVersion = version
	fmt.Println("Часы : ", whatch.Model, ",", "Версия: ", whatch.OSVersion)
	return nil
}

func main() {
	tel := Smartphone{
		OSVersion: "11.6",
		Model:     "Iphone",
	}
	tel.GetInfo()
	tel.UpdateOS("5")

	HP := Laptop{
		OSVersion: "7",
		Model:     "Windows",
	}
	Mac := Laptop{
		OSVersion: "17",
		Model:     "Mac",
	}
	HP.UpdateOS("Windows 11")
	Mac.UpdateOS("Mac 11")

	Apple := Smartwatch{
		OSVersion: "3.1.0",
		Model:     "Apple",
	}

	Apple.GetInfo()
	Apple.UpdateOS("5.1.5.3")
}
