package sshclient

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
