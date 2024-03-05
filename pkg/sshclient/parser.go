package sshclient

import (
	"fmt"
	"strconv"
	"strings"
)

type UnmanagedONT struct {
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

func (o *UnmanagedONT) GetFrameSlotPort() (int, int, int, error) {
	return getFrameSlotPortFromFSP(o.FSP)
}

func ParseUnmanagedONT(output string) ([]UnmanagedONT, error) {
	results := make([]UnmanagedONT, 0)

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

		ont := UnmanagedONT{
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

type OpticalInfo struct {
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

func ParseOpticalInfo(output string) (*OpticalInfo, error) {
	lines := strings.Split(output, "\n")

	err := parseLinesFailure(lines)
	if err != nil {
		return nil, err
	}

	details := &OpticalInfo{}

	fieldMap := map[string]*string{
		"ONU NNI port ID                        : ": &details.ONUNNIPortID,
		"Module type                            : ": &details.ModuleType,
		"Module sub-type                        : ": &details.ModuleSubType,
		"Used type                              : ": &details.UsedType,
		"Encapsulation Type                     : ": &details.EncapsulationType,
		"Optical power precision(dBm)           : ": &details.OpticalPowerPrecision,
		"Vendor name                            : ": &details.VendorName,
		"Vendor rev                             : ": &details.VendorRev,
		"Vendor PN                              : ": &details.VendorPN,
		"Vendor SN                              : ": &details.VendorSN,
		"Date Code                              : ": &details.DateCode,
		"Rx optical power(dBm)                  : ": &details.RxOpticalPower,
		"Rx power current warning threshold(dBm): ": &details.RxPowerCurrentWarningThreshold,
		"Rx power current alarm threshold(dBm)  : ": &details.RxPowerCurrentAlarmThreshold,
		"Tx optical power(dBm)                  : ": &details.TxOpticalPower,
		"Tx power current warning threshold(dBm): ": &details.TxPowerCurrentWarningThreshold,
		"Tx power current alarm threshold(dBm)  : ": &details.TxPowerCurrentAlarmThreshold,
		"Laser bias current(mA)                 : ": &details.LaserBiasCurrent,
		"Tx bias current warning threshold(mA)  : ": &details.TxBiasCurrentWarningThreshold,
		"Tx bias current alarm threshold(mA)    : ": &details.TxBiasCurrentAlarmThreshold,
		"Temperature(C)                         : ": &details.Temperature,
		"Temperature warning threshold(C)       : ": &details.TemperatureWarningThreshold,
		"Temperature alarm threshold(C)         : ": &details.TemperatureAlarmThreshold,
		"Voltage(V)                             : ": &details.Voltage,
		"Supply voltage warning threshold(V)    : ": &details.SupplyVoltageWarningThreshold,
		"Supply voltage alarm threshold(V)      : ": &details.SupplyVoltageAlarmThreshold,
		"OLT Rx ONT optical power(dBm)          : ": &details.OLTRxONTOpticalPower,
		"CATV Rx optical power(dBm)             : ": &details.CATVRxOpticalPower,
		"CATV Rx power alarm threshold(dBm)     : ": &details.CATVRxPowerAlarmThreshold,
	}

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		for prefix, field := range fieldMap {
			if strings.HasPrefix(trimmedLine, prefix) {
				*field = strings.TrimPrefix(trimmedLine, prefix)
				break
			}
		}
	}

	return details, nil
}

type GeneralInfo struct {
	FSP               string `json:"fsp"`
	ID                string `json:"id"`
	ControlFlag       string `json:"control_flag"`
	RunState          string `json:"run_state"`
	ConfigState       string `json:"config_state"`
	MatchState        string `json:"match_state"`
	DBAType           string `json:"dba_type"`
	Distance          string `json:"distance"`
	LastDistance      string `json:"last_distance"`
	BatteryState      string `json:"battery_state"`
	MemoryOccupation  string `json:"memory_occupation"`
	CPUOccupation     string `json:"cpu_occupation"`
	Temperature       string `json:"temperature"`
	AuthenticType     string `json:"authentic_type"`
	SN                string `json:"sn"`
	ManagementMode    string `json:"management_mode"`
	SoftwareWorkMode  string `json:"software_work_mode"`
	IsolationState    string `json:"isolation_state"`
	Description       string `json:"description"`
	LatDownCause      string `json:"last_down_cause"`
	LastUpTime        string `json:"last_up_time"`
	LastDownTime      string `json:"last_down_time"`
	LastDyingGaspTime string `json:"last_dying_gasp_time"`
	OnlineDuration    string `json:"online_duration"`
}

func (o *GeneralInfo) GetFrameSlotPort() (int, int, int, error) {
	return getFrameSlotPortFromFSP(o.FSP)
}

func ParseGeneralInfoBySn(output string) (*GeneralInfo, error) {
	lines := strings.Split(output, "\n")

	ont := &GeneralInfo{}

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "The required ONT does not exist" {
			return nil, NotFoundError{}
		}
		if strings.Contains(trimmedLine, "Parameter error") {
			return nil, InvalidSerialNumberError{}
		}

		fieldMap := map[string]*string{
			"F/S/P                   : ": &ont.FSP,
			"ONT-ID                  : ": &ont.ID,
			"Control flag            : ": &ont.ControlFlag,
			"Run state               : ": &ont.RunState,
			"Config state            : ": &ont.ConfigState,
			"Match state             : ": &ont.MatchState,
			"DBA type                : ": &ont.DBAType,
			"ONT distance(m)         : ": &ont.Distance,
			"ONT last distance(m)    : ": &ont.LastDistance,
			"ONT battery state       : ": &ont.BatteryState,
			"Memory occupation       : ": &ont.MemoryOccupation,
			"CPU occupation          : ": &ont.CPUOccupation,
			"Temperature             : ": &ont.Temperature,
			"Authentic type          : ": &ont.AuthenticType,
			"SN                      : ": &ont.SN,
			"Management mode         : ": &ont.ManagementMode,
			"Software work mode      : ": &ont.SoftwareWorkMode,
			"Isolation state         : ": &ont.IsolationState,
			"Description             : ": &ont.Description,
			"Last down cause         : ": &ont.LatDownCause,
			"Last up time            : ": &ont.LastUpTime,
			"Last down time          : ": &ont.LastDownTime,
			"Last dying gasp time    : ": &ont.LastDyingGaspTime,
			"ONT online duration     : ": &ont.OnlineDuration,
		}

		for prefix, field := range fieldMap {
			if strings.HasPrefix(trimmedLine, prefix) {
				*field = strings.TrimPrefix(trimmedLine, prefix)
				break
			}
		}
	}

	return ont, nil
}

type ServicePort struct {
	Index int `json:"index"`
	Vlan  int `json:"vlan"`
}

func ParseServicePorts(output string) ([]ServicePort, error) {
	results := make([]ServicePort, 0)

	lines := strings.Split(output, "\n")
	if strings.Contains(lines[5], "Failure: No service virtual port can be operated") {
		return results, nil
	}

	err := parseLinesFailure(lines)
	if err != nil {
		return nil, err
	}

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

func getFrameSlotPortFromFSP(fsp string) (int, int, int, error) {
	parts := strings.Split(fsp, "/")
	frame, err := strconv.Atoi(parts[0])
	slot, err := strconv.Atoi(parts[1])
	port, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, 0, 0, err
	}
	return frame, slot, port, nil
}

func parseFailure(line string) error {
	if strings.Contains(line, "Failure: ") {
		return fmt.Errorf(strings.TrimPrefix(strings.TrimSpace(line), "Failure: "))
	}
	return nil
}

func parseLinesFailure(lines []string) error {
	for _, line := range lines {
		err := parseFailure(line)
		if err != nil {
			return err
		}
	}
	return nil
}
