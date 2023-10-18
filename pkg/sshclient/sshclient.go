package sshclient

type ONT struct {
	ID           string
	Frame        string
	Slot         string
	Port         string
	SerialNumber string
	VlanID       string
	Description  string
	ServicePort  string
}

type UnmanagedONT struct {
	Number             string
	FSP                string
	OntSN              string
	Password           string
	Loid               string
	Checkcode          string
	VendorID           string
	OntVersion         string
	OntSoftwareVersion string
	OntEquipmentID     string
	OntCustomizedInfo  string
	OntAutofindTime    string
}
