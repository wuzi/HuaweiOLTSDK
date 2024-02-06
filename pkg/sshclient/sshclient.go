package sshclient

type ONT struct {
	ID           string `json:"id"`
	Frame        string `json:"frame"`
	Slot         string `json:"slot"`
	Port         string `json:"port"`
	SerialNumber string `json:"serial_number"`
	VlanID       string `json:"vlan_id"`
	Description  string `json:"description"`
	ServicePort  string `json:"service_port"`
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
