package sshclient

import (
	"strconv"
	"strings"
)

func ParseUnmanagedONT(output string) ([]ONTDetail, error) {
	results := make([]ONTDetail, 0)

	sections := strings.Split(output, "   ----------------------------------------------------------------------------")

	for _, section := range sections[1:] {
		lines := strings.Split(section, "\n")

		var cleanLines []string
		for _, line := range lines {
			if line = strings.TrimSpace(line); line != "" && !strings.Contains(line, "---- More (") {
				cleanLines = append(cleanLines, line)
			}
		}

		if len(cleanLines) < 11 {
			continue
		}

		ont := ONTDetail{
			Number:             strings.TrimSpace(strings.TrimPrefix(cleanLines[0], "Number              :")),
			FSP:                strings.TrimSpace(strings.TrimPrefix(cleanLines[1], "F/S/P               :")),
			OntSN:              strings.TrimSpace(strings.TrimPrefix(cleanLines[2], "Ont SN              :")),
			Password:           strings.TrimSpace(strings.TrimPrefix(cleanLines[3], "Password            :")),
			Loid:               strings.TrimSpace(strings.TrimPrefix(cleanLines[4], "Loid                :")),
			Checkcode:          strings.TrimSpace(strings.TrimPrefix(cleanLines[5], "Checkcode           :")),
			VendorID:           strings.TrimSpace(strings.TrimPrefix(cleanLines[6], "VendorID            :")),
			OntVersion:         strings.TrimSpace(strings.TrimPrefix(cleanLines[7], "Ont Version         :")),
			OntSoftwareVersion: strings.TrimSpace(strings.TrimPrefix(cleanLines[8], "Ont SoftwareVersion :")),
			OntEquipmentID:     strings.TrimSpace(strings.TrimPrefix(cleanLines[9], "Ont EquipmentID     :")),
			OntCustomizedInfo:  strings.TrimSpace(strings.TrimPrefix(cleanLines[10], "Ont Customized Info :")),
			OntAutofindTime:    strings.TrimSpace(strings.TrimPrefix(cleanLines[11], "Ont autofind time   :")),
		}

		results = append(results, ont)
	}

	return results, nil
}

type OnuOpticalInfo struct {
	ONUNNIPortID                   string
	ModuleType                     string
	ModuleSubType                  string
	UsedType                       string
	EncapsulationType              string
	OpticalPowerPrecision          string
	VendorName                     string
	VendorRev                      string
	VendorPN                       string
	VendorSN                       string
	DateCode                       string
	RxOpticalPower                 string
	RxPowerCurrentWarningThreshold string
	RxPowerCurrentAlarmThreshold   string
	TxOpticalPower                 string
	TxPowerCurrentWarningThreshold string
	TxPowerCurrentAlarmThreshold   string
	LaserBiasCurrent               string
	TxBiasCurrentWarningThreshold  string
	TxBiasCurrentAlarmThreshold    string
	Temperature                    string
	TemperatureWarningThreshold    string
	TemperatureAlarmThreshold      string
	Voltage                        string
	SupplyVoltageWarningThreshold  string
	SupplyVoltageAlarmThreshold    string
	OLTRxONTOpticalPower           string
	CATVRxOpticalPower             string
	CATVRxPowerAlarmThreshold      string
}

func ParseOnuOpticalInfo(output string) *OnuOpticalInfo {
	lines := strings.Split(output, "\n")

	details := &OnuOpticalInfo{
		ONUNNIPortID:                   strings.TrimPrefix(strings.TrimSpace(lines[2]), "ONU NNI port ID                        : "),
		ModuleType:                     strings.TrimPrefix(strings.TrimSpace(lines[3]), "Module type                            : "),
		ModuleSubType:                  strings.TrimPrefix(strings.TrimSpace(lines[4]), "Module sub-type                        : "),
		UsedType:                       strings.TrimPrefix(strings.TrimSpace(lines[5]), "Used type                              : "),
		EncapsulationType:              strings.TrimPrefix(strings.TrimSpace(lines[6]), "Encapsulation Type                     : "),
		OpticalPowerPrecision:          strings.TrimPrefix(strings.TrimSpace(lines[7]), "Optical power precision(dBm)           : "),
		VendorName:                     strings.TrimPrefix(strings.TrimSpace(lines[8]), "Vendor name                            : "),
		VendorRev:                      strings.TrimPrefix(strings.TrimSpace(lines[9]), "Vendor rev                             : "),
		VendorPN:                       strings.TrimPrefix(strings.TrimSpace(lines[10]), "Vendor PN                              : "),
		VendorSN:                       strings.TrimPrefix(strings.TrimSpace(lines[11]), "Vendor SN                              : "),
		DateCode:                       strings.TrimPrefix(strings.TrimSpace(lines[12]), "Date Code                              : "),
		RxOpticalPower:                 strings.TrimPrefix(strings.TrimSpace(lines[13]), "Rx optical power(dBm)                  : "),
		RxPowerCurrentWarningThreshold: strings.TrimPrefix(strings.TrimSpace(lines[14]), "Rx power current warning threshold(dBm): "),
		RxPowerCurrentAlarmThreshold:   strings.TrimPrefix(strings.TrimSpace(lines[15]), "Rx power current alarm threshold(dBm)  : "),
		TxOpticalPower:                 strings.TrimPrefix(strings.TrimSpace(lines[16]), "Tx optical power(dBm)                  : "),
		TxPowerCurrentWarningThreshold: strings.TrimPrefix(strings.TrimSpace(lines[17]), "Tx power current warning threshold(dBm): "),
		TxPowerCurrentAlarmThreshold:   strings.TrimPrefix(strings.TrimSpace(lines[18]), "Tx power current alarm threshold(dBm)  : "),
		LaserBiasCurrent:               strings.TrimPrefix(strings.TrimSpace(lines[19]), "Laser bias current(mA)                 : "),
		TxBiasCurrentWarningThreshold:  strings.TrimPrefix(strings.TrimSpace(lines[20]), "Tx bias current warning threshold(mA)  : "),
		TxBiasCurrentAlarmThreshold:    strings.TrimPrefix(strings.TrimSpace(lines[21]), "Tx bias current alarm threshold(mA)    : "),
		Temperature:                    strings.TrimPrefix(strings.TrimSpace(lines[22]), "Temperature(C)                         : "),
		TemperatureWarningThreshold:    strings.TrimPrefix(strings.TrimSpace(lines[23]), "Temperature warning threshold(C)       : "),
		TemperatureAlarmThreshold:      strings.TrimPrefix(strings.TrimSpace(lines[24]), "Temperature alarm threshold(C)         : "),
		Voltage:                        strings.TrimPrefix(strings.TrimSpace(lines[25]), "Voltage(V)                             : "),
		SupplyVoltageWarningThreshold:  strings.TrimPrefix(strings.TrimSpace(lines[26]), "Supply voltage warning threshold(V)    : "),
		SupplyVoltageAlarmThreshold:    strings.TrimPrefix(strings.TrimSpace(lines[27]), "Supply voltage alarm threshold(V)      : "),
		OLTRxONTOpticalPower:           strings.TrimPrefix(strings.TrimSpace(lines[28]), "OLT Rx ONT optical power(dBm)          : "),
		CATVRxOpticalPower:             strings.TrimPrefix(strings.TrimSpace(lines[29]), "CATV Rx optical power(dBm)             : "),
		CATVRxPowerAlarmThreshold:      strings.TrimPrefix(strings.TrimSpace(lines[30]), "CATV Rx power alarm threshold(dBm)     : "),
	}

	return details
}

type OntGeneralInfo struct {
	FSP              string
	ID               string
	ControlFlag      string
	RunState         string
	ConfigState      string
	MatchState       string
	DBAType          string
	Distance         string
	LastDistance     string
	BatteryState     string
	AuthenticType    string
	SN               string
	ManagementMode   string
	SoftwareWorkMode string
	IsolationState   string
	Description      string
	LatDownCause     string
	LastUpTime       string
	LastDownTime     string
	OnlineDuration   string
}

func (o *OntGeneralInfo) GetFrameSlotPort() (*int, *int, *int) {
	parts := strings.Split(o.FSP, "/")
	frame, err := strconv.Atoi(parts[0])
	slot, err := strconv.Atoi(parts[1])
	port, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, nil, nil
	}
	return &frame, &slot, &port
}

func ParseOntInfoBySn(output string) *OntGeneralInfo {
	lines := strings.Split(output, "\n")

	if strings.TrimSpace(lines[1]) == "The required ONT does not exist" {
		return nil
	}

	ont := &OntGeneralInfo{
		FSP:              strings.TrimPrefix(strings.TrimSpace(lines[2]), "F/S/P                   : "),
		ID:               strings.TrimPrefix(strings.TrimSpace(lines[3]), "ONT-ID                  : "),
		ControlFlag:      strings.TrimPrefix(strings.TrimSpace(lines[4]), "Control flag            : "),
		RunState:         strings.TrimPrefix(strings.TrimSpace(lines[5]), "Run state               : "),
		ConfigState:      strings.TrimPrefix(strings.TrimSpace(lines[6]), "Config state            : "),
		MatchState:       strings.TrimPrefix(strings.TrimSpace(lines[7]), "Match state             : "),
		DBAType:          strings.TrimPrefix(strings.TrimSpace(lines[8]), "DBA type                : "),
		Distance:         strings.TrimPrefix(strings.TrimSpace(lines[9]), "ONT distance(m)         : "),
		LastDistance:     strings.TrimPrefix(strings.TrimSpace(lines[10]), "ONT last distance(m)    : "),
		BatteryState:     strings.TrimPrefix(strings.TrimSpace(lines[11]), "ONT battery state       : "),
		AuthenticType:    strings.TrimPrefix(strings.TrimSpace(lines[15]), "Authentic type          : "),
		SN:               strings.TrimPrefix(strings.TrimSpace(lines[16]), "SN                      : "),
		ManagementMode:   strings.TrimPrefix(strings.TrimSpace(lines[17]), "Management mode         : "),
		SoftwareWorkMode: strings.TrimPrefix(strings.TrimSpace(lines[18]), "Software work mode      : "),
		IsolationState:   strings.TrimPrefix(strings.TrimSpace(lines[19]), "Isolation state         : "),
		Description:      strings.TrimPrefix(strings.TrimSpace(lines[21]), "Description             : "),
		LatDownCause:     strings.TrimPrefix(strings.TrimSpace(lines[22]), "Last down cause         : "),
		LastUpTime:       strings.TrimPrefix(strings.TrimSpace(lines[23]), "Last up time            : "),
		LastDownTime:     strings.TrimPrefix(strings.TrimSpace(lines[24]), "Last down time          : "),
		OnlineDuration:   strings.TrimPrefix(strings.TrimSpace(lines[25]), "ONT online duration     : "),
	}

	return ont
}
