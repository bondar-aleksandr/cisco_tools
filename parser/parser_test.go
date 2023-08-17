package parser

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var testDataDir = "./../test_data"

func getCiscoInterfaceMap(filename string) CiscoInterfaceMap {
	jsonFile, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Cannot open file %s", filename)
	}

	var result CiscoInterfaceMap
	err = json.Unmarshal(jsonFile, &result)
	if err != nil {
		log.Fatalf("Cannot deserialize file %s into JSON", filename)
	}
	return result
}

func Test_parsing(t *testing.T) {

	ios_ifile_router := filepath.Join(testDataDir, "INET-R01.txt")
	ios_ifiile_switch := filepath.Join(testDataDir, "run.txt")
	ios_ifile_routerXR := filepath.Join(testDataDir, "ASR-P.txt")
	nxos_ifile := filepath.Join(testDataDir, "dc0-n9k-d_23.08.txt")

	ios_map_router := getCiscoInterfaceMap(FileExtReplace(ios_ifile_router, "json"))
	ios_map_switch := getCiscoInterfaceMap(FileExtReplace(ios_ifiile_switch, "json"))
	ios_map_routerXR := getCiscoInterfaceMap(FileExtReplace(ios_ifile_routerXR, "json"))
	nxos_map := getCiscoInterfaceMap(FileExtReplace(nxos_ifile, "json"))

	configs := []struct {
		name     string
		ifile    string
		dev_type string
		expected CiscoInterfaceMap
	}{
		{name: "ios-router", ifile: ios_ifile_router, dev_type: "ios", expected: ios_map_router},
		{name: "ios-L3switch", ifile: ios_ifiile_switch, dev_type: "ios", expected: ios_map_switch},
		{name: "ios-XR", ifile: ios_ifile_routerXR, dev_type: "ios", expected: ios_map_routerXR},
		{name: "NXOS", ifile: nxos_ifile, dev_type: "nxos", expected: nxos_map},
	}

	for _, v := range configs {

		ifile := v.ifile
		device := v.dev_type
		target_map := v.expected
		f, err := os.Open(ifile)
		if err != nil {
			t.Errorf("Cannot open configuration file %s because of %q", ifile, err)
		}

		interface_map, _ := Parsing(f, device)
		eq := reflect.DeepEqual(interface_map, target_map)
		if !eq {
			t.Errorf("%s: parsed config doesn't correspond target value", v.name)
		}
	}
}
