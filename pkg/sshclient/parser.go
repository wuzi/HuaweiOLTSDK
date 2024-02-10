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
	ONUNNIPortID                   string `json:"onu_nni_port_id"`
	ModuleType                     string `json:"module_type"`
	ModuleSubType                  string `json:"module_sub_type"`
	UsedType                       string `json:"used_type"`
	EncapsulationType              string `json:"encapsulation_type"`
	OpticalPowerPrecision          string `json:"optical_power_precision"`
	VendorName                     string `json:"vendor_name"`
	VendorRev                      string `json:"vendor_rev"`
	VendorPN                       string `json:"vendor_pn"`
	VendorSN                       string `json:"vendor_sn"`
	DateCode                       string `json:"date_code"`
	RxOpticalPower                 string `json:"rx_optical_power"`
	RxPowerCurrentWarningThreshold string `json:"rx_power_current_warning_threshold"`
	RxPowerCurrentAlarmThreshold   string `json:"rx_power_current_alarm_threshold"`
	TxOpticalPower                 string `json:"tx_optical_power"`
	TxPowerCurrentWarningThreshold string `json:"tx_power_current_warning_threshold"`
	TxPowerCurrentAlarmThreshold   string `json:"tx_power_current_alarm_threshold"`
	LaserBiasCurrent               string `json:"laser_bias_current"`
	TxBiasCurrentWarningThreshold  string `json:"tx_bias_current_warning_threshold"`
	TxBiasCurrentAlarmThreshold    string `json:"tx_bias_current_alarm_threshold"`
	Temperature                    string `json:"temperature"`
	TemperatureWarningThreshold    string `json:"temperature_warning_threshold"`
	TemperatureAlarmThreshold      string `json:"temperature_alarm_threshold"`
	Voltage                        string `json:"voltage"`
	SupplyVoltageWarningThreshold  string `json:"supply_voltage_warning_threshold"`
	SupplyVoltageAlarmThreshold    string `json:"supply_voltage_alarm_threshold"`
	OLTRxONTOpticalPower           string `json:"olt_rx_ont_optical_power"`
	CATVRxOpticalPower             string `json:"catv_rx_optical_power"`
	CATVRxPowerAlarmThreshold      string `json:"catv_rx_power_alarm_threshold"`
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
	FSP              string `json:"fsp"`
	ID               string `json:"id"`
	ControlFlag      string `json:"control_flag"`
	RunState         string `json:"run_state"`
	ConfigState      string `json:"config_state"`
	MatchState       string `json:"match_state"`
	DBAType          string `json:"dba_type"`
	Distance         string `json:"distance"`
	LastDistance     string `json:"last_distance"`
	BatteryState     string `json:"battery_state"`
	AuthenticType    string `json:"authentic_type"`
	SN               string `json:"sn"`
	ManagementMode   string `json:"management_mode"`
	SoftwareWorkMode string `json:"software_work_mode"`
	IsolationState   string `json:"isolation_state"`
	Description      string `json:"description"`
	LatDownCause     string `json:"last_down_cause"`
	LastUpTime       string `json:"last_up_time"`
	LastDownTime     string `json:"last_down_time"`
	OnlineDuration   string `json:"online_duration"`
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

type ServicePort struct {
	Index int
	Vlan  int
}

func ParseServicePorts(output string) ([]ServicePort, error) {
	results := make([]ServicePort, 0)

	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if len(line) > 0 && line[0] >= '0' && line[0] <= '9' {
			parts := strings.Fields(line)
			index, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, err
			}
			vlan, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, err
			}
			results = append(results, ServicePort{Index: index, Vlan: vlan})
		}
	}

	return results, nil
}
