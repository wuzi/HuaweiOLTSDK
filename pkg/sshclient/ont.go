package sshclient

import (
	"strconv"
	"strings"
)

type ONT struct {
	Number             int
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

func ParseUnmanagedONT(output string) ([]ONT, error) {
	var results []ONT

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

		number, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(cleanLines[0], "Number              :")))
		if err != nil {
			return nil, err
		}

		ont := ONT{
			Number:             number,
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