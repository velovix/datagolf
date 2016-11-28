package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/tarm/serial"
)

const bufSize = 1

type device struct {
	*serial.Port
	*sync.Mutex
}

type xyz struct {
	x int
	y int
	z int
}

// newDevice creates a new device object that is connected to the given serial
// device name.
func newDevice(serialName string) (device, error) {
	config := &serial.Config{
		Name: serialName,
		Baud: 9600}
	port, err := serial.OpenPort(config)
	if err != nil {
		return device{}, errors.Wrap(err, "while connecting to device")
	}

	return device{
		Port:  port,
		Mutex: &sync.Mutex{}}, nil
}

// data fetches data from the device and returns it.
func (d device) data() ([]xyz, []xyz, error) {
	d.Lock()
	defer d.Unlock()

	fmt.Println("about to start read")

	d.Write([]byte("?"))

	fmt.Println("Starting accel")
	accel, err := d.readDataLine()
	if err != nil {
		return make([]xyz, 0), make([]xyz, 0), err
	}
	fmt.Println("Starting gyro")
	gyro, err := d.readDataLine()
	if err != nil {
		return make([]xyz, 0), make([]xyz, 0), err
	}

	fmt.Println(accel)
	fmt.Println(gyro)

	return accel, gyro, nil
}

func (d device) readDataLine() ([]xyz, error) {
	raw := ""
	buf := make([]byte, bufSize)

	// Read all the data the device has. Data is ended with a DOS-style newline
	for !strings.Contains(raw, "\n") {
		_, err := d.Read(buf)
		if err == io.EOF {
			fmt.Println("Got EOF")
			break
		}
		if err != nil && err != io.EOF {
			return make([]xyz, 0), errors.Wrap(err, "reading data from device")
		}
		raw += string(buf)
	}

	fmt.Println("Finished with", strings.TrimSpace(raw))

	// Split the string by commas, minus whitespaces
	dataStr := strings.Split(raw, ",")

	// Convert the split data to ints
	var flatData []int
	for _, val := range dataStr {
		i, err := strconv.Atoi(strings.TrimSpace(val))
		if err != nil {
			return make([]xyz, 0), errors.Wrap(err, "device returned invalid data")
		}
		flatData = append(flatData, i)
	}

	// Compartmentalize the ints into "points"
	var data []xyz
	for i := 0; i < len(flatData); i += 3 {
		data = append(data, xyz{
			x: flatData[i],
			y: flatData[i+1],
			z: flatData[i+2]})
	}

	return data, nil
}
