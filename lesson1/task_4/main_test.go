package main

import (
	"testing"
)

func TestCarStartEngine(t *testing.T) {
	car := &Car{Brand: "Toyota"}

	err := car.StartEngine()
	if err != nil {
		t.Errorf("Ожидался nil, получили: %v", err)
	}
	if !car.GetEngineStatus() {
		t.Errorf("Двигатель должен быть включён, но GetEngineStatus() = false")
	}
}

func TestCarStartEngineAlreadyRunning(t *testing.T) {
	car := &Car{Brand: "Toyota"}
	car.StartEngine()
	err := car.StartEngine()

	if err != ErrEngineAlreadyRunning {
		t.Errorf("Ожидалась ошибка %v, получили: %v", ErrEngineAlreadyRunning, err)
	}
	if !car.GetEngineStatus() {
		t.Errorf("Двигатель должен оставаться включённым, но GetEngineStatus() = false")
	}
}

func TestCarStopEngine(t *testing.T) {
	car := &Car{Brand: "Toyota"}
	car.StartEngine()

	err := car.StopEngine()

	if err != nil {
		t.Errorf("Ожидался nil, получили: %v", err)
	}
	if car.GetEngineStatus() {
		t.Errorf("Двигатель должен быть выключен, но GetEngineStatus() = true")
	}
}

func TestCarStopEngineAlreadyOff(t *testing.T) {
	car := &Car{Brand: "Toyota"}

	err := car.StopEngine()

	if err != ErrEngineOff {
		t.Errorf("Ожидалась ошибка %v, получили: %v", ErrEngineOff, err)
	}
}

func TestCarGetInfo(t *testing.T) {
	car := &Car{Brand: "BMW"}

	info := car.GetInfo()

	expected := "Машина BMW"
	if info != expected {
		t.Errorf("GetInfo() = %s, ожидалось %s", info, expected)
	}
}

func TestCarHonk(t *testing.T) {
	car := &Car{Brand: "Toyota"}

	honk := car.Honk()

	expected := "Beep beep!"
	if honk != expected {
		t.Errorf("Honk() = %s, ожидалось %s", honk, expected)
	}
}

func TestTruckHonk(t *testing.T) {
	truck := &Truck{}

	honk := truck.Honk()

	expected := "Honk Honk!"
	if honk != expected {
		t.Errorf("Honk() = %s, ожидалось %s", honk, expected)
	}
}

func TestTruckGetCargoCapacity(t *testing.T) {
	truck := &Truck{cargoCapacity: 15.5}

	capacity := truck.GetCargoCapacity()

	expected := 15.5
	if capacity != expected {
		t.Errorf("GetCargoCapacity() = %f, ожидалось %f", capacity, expected)
	}
}

func TestTruckInheritsCarMethods(t *testing.T) {
	truck := &Truck{}
	truck.Brand = "Volvo"

	err := truck.StartEngine()
	info := truck.GetInfo()

	if err != nil {
		t.Errorf("StartEngine() вернул ошибку: %v", err)
	}
	if !truck.GetEngineStatus() {
		t.Errorf("Двигатель грузовика должен быть включён")
	}
	expectedInfo := "Машина Volvo"
	if info != expectedInfo {
		t.Errorf("GetInfo() = %s, ожидалось %s", info, expectedInfo)
	}
}

func TestElectricCarStartEngineWithGoodBattery(t *testing.T) {
	electricCar := &ElectricCar{
		Car:          Car{Brand: "Tesla"},
		batteryLevel: 50,
	}

	err := electricCar.StartEngine()

	if err != nil {
		t.Errorf("Ожидался nil, получили: %v", err)
	}
	if !electricCar.GetEngineStatus() {
		t.Errorf("Двигатель должен быть включён, но GetEngineStatus() = false")
	}
}

func TestElectricCarStartEngineWithLowBattery(t *testing.T) {
	electricCar := &ElectricCar{
		Car:          Car{Brand: "Tesla"},
		batteryLevel: 3,
	}

	err := electricCar.StartEngine()

	if err != ErrLowBattery {
		t.Errorf("Ожидалась ошибка %v, получили: %v", ErrLowBattery, err)
	}
	if electricCar.GetEngineStatus() {
		t.Errorf("Двигатель не должен включаться при низком заряде, но GetEngineStatus() = true")
	}
}

func TestElectricCarStartEngineWithBorderBattery(t *testing.T) {
	electricCar := &ElectricCar{
		Car:          Car{Brand: "Tesla"},
		batteryLevel: 5,
	}

	err := electricCar.StartEngine()

	if err != ErrLowBattery {
		t.Errorf("При заряде 5%% ожидалась ошибка %v, получили: %v", ErrLowBattery, err)
	}
}

func TestElectricCarGetBatteryLevel(t *testing.T) {
	electricCar := &ElectricCar{batteryLevel: 75}

	level := electricCar.GetBatteryLevel()

	expected := 75
	if level != expected {
		t.Errorf("GetBatteryLevel() = %d, ожидалось %d", level, expected)
	}
}

func TestElectricCarGetInfo(t *testing.T) {
	electricCar := &ElectricCar{
		Car:          Car{Brand: "Nissan"},
		batteryLevel: 80,
	}

	info := electricCar.GetInfo()

	expected := "Машина Nissanзаряд батареи"
	if info != expected {
		t.Errorf("GetInfo() = %s, ожидалось %s", info, expected)
	}
}

func TestPolymorphism(t *testing.T) {
	vehicles := []Vehicle{
		&Car{Brand: "Toyota"},
		&Truck{Car: Car{Brand: "Volvo"}, cargoCapacity: 10},
		&ElectricCar{Car: Car{Brand: "Tesla"}, batteryLevel: 50},
	}

	for i, v := range vehicles {
		err := v.StartEngine()
		if err != nil {
			t.Errorf("Vehicle[%d] StartEngine() вернул ошибку: %v", i, err)
		}

		info := v.GetInfo()
		if info == "" {
			t.Errorf("Vehicle[%d] GetInfo() вернул пустую строку", i)
		}

		err = v.StopEngine()
		if err != nil {
			t.Errorf("Vehicle[%d] StopEngine() вернул ошибку: %v", i, err)
		}
	}
}

func TestCarGetEngineStatus(t *testing.T) {
	car := &Car{Brand: "Honda"}

	if car.GetEngineStatus() {
		t.Errorf("Новый автомобиль должен иметь выключенный двигатель")
	}

	car.StartEngine()

	if !car.GetEngineStatus() {
		t.Errorf("После StartEngine() двигатель должен быть включён")
	}

	car.StopEngine()

	if car.GetEngineStatus() {
		t.Errorf("После StopEngine() двигатель должен быть выключен")
	}
}
