package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
)

type dataHandler struct {
	d device
}

type dataResp struct {
	Accel []int `json:"accel"`
	Gyro  []int `json:"gyro"`
}

func (h *dataHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	accelData := make([]int, 100)
	for i := 0; i < 100; i++ {
		accelData[i] = rand.Intn(128)
	}
	gyroData := make([]int, 100)
	for i := 0; i < 100; i++ {
		gyroData[i] = rand.Intn(128)
	}

	resp := dataResp{
		Accel: accelData,
		Gyro:  gyroData,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "could not send data", 500)
	}

	w.Write(data)
}
