package main

import (
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/tarm/serial"
)

const bufSize = 16

type device struct {
	*serial.Port
	*sync.Mutex
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
func (d device) data() ([]int, error) {
	d.Lock()
	defer d.Unlock()

	d.Write([]byte("?\r\n"))

	raw := ""
	buf := make([]byte, bufSize)

	// Read all the data the device has. Data is ended with a DOS-style newline
	for strings.Contains(raw, "\r\n") {
		_, err := d.Read(buf)
		if err != nil {
			return make([]int, 0), errors.Wrap(err, "reading data from device")
		}
		raw += string(buf)
	}

	// Split the string by commas, minus the newline
	dataStr := strings.Split(raw[:len(raw)-2], ",")

	// Convert the split data to ints
	var data []int
	for _, val := range dataStr {
		i, err := strconv.Atoi(val)
		if err != nil {
			return make([]int, 0), errors.Wrap(err, "device returned invalid data")
		}
		data = append(data, i)
	}

	return data, nil
}
