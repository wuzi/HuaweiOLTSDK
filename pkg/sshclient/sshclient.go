package sshclient

import (
	"strconv"
	"strings"
)

type ONT struct {
	ID           int    `json:"id"`
	Frame        int    `json:"frame"`
	Slot         int    `json:"slot"`
	Port         int    `json:"port"`
	SerialNumber string `json:"serial_number"`
	VlanID       int    `json:"vlan_id"`
	Description  string `json:"description"`
	ServicePort  int    `json:"service_port"`
}

type ONTDetail struct {
	Number             string `json:"number"`
	FSP                string `json:"fsp"`
	OntSN              string `json:"serial_number"`
	Password           string `json:"password"`
	Loid               string `json:"lo_id"`
	Checkcode          string `json:"check_code"`
	VendorID           string `json:"vendor_id"`
	OntVersion         string `json:"version"`
	OntSoftwareVersion string `json:"software_version"`
	OntEquipmentID     string `json:"equipment_id"`
	OntCustomizedInfo  string `json:"customized_info"`
	OntAutofindTime    string `json:"auto_find_time"`
}

func (o *ONTDetail) GetFrameSlotPort() (*int, *int, *int) {
	parts := strings.Split(o.FSP, "/")
	frame, err := strconv.Atoi(parts[0])
	slot, err := strconv.Atoi(parts[1])
	port, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, nil, nil
	}
	return &frame, &slot, &port
}
