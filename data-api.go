package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type dataHandler struct {
	dev device
}

type dataResp struct {
	Accel []int `json:"accel"`
	Gyro  []int `json:"gyro"`
}

func (h *dataHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	accel, gyro, err := h.dev.data()
	if err != nil {
		http.Error(w, "could not get data", 500)
		log.Println(err)
		return
	}

	resp := dataResp{
		Accel: xyzAverage(accel),
		Gyro:  xyzAverage(gyro),
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "could not send data", 500)
		log.Println(err)
		return
	}

	w.Write(data)
}

func xyzAverage(data []xyz) []int {
	averages := make([]int, len(data))

	for i, pnt := range data {
		averages[i] = (pnt.x + pnt.y + pnt.z) / 3
	}

	return averages
}
