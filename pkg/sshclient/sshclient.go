package sshclient

import "strings"

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

func (o *ONTDetail) GetFrameSlotPort() (string, string, string) {
	parts := strings.Split(o.FSP, "/")
	return parts[0], parts[1], parts[2]
}
