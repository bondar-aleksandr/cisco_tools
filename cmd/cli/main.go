package main

import (
	"github.com/bondar-aleksandr/cisco_parser"
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	var iFileName = flag.String("i", "", "input configuration file to parse data from")
	var oFileName = flag.String("o", "", "output csv file, default is input filename with .csv extension")
	var devtype = flag.String("t", "ios", "cisco OS family, possible values are ios, nxos. Default is ios")
	var jsonOut = flag.Bool("j", false, "Whether JSON file needed. Default is false")

	flag.Parse()
	log.Infof(`Program started, got the following parameters: input file: %s, output file: %s,
			device type: %s, JSON output: %v`, *iFileName, *oFileName, *devtype, *jsonOut)

	iFile, err := os.Open(*iFileName)
	if err != nil {
		log.Fatalf("Can not open file %s because of: %q", iFile.Name(), err)
	}
	defer iFile.Close()

	interface_map, err := cisco_parser.ParseInterfaces(iFile, *devtype)
	if err != nil {
		log.Fatal(err)
	}	

	if *oFileName == "" {
		*oFileName = cisco_parser.FileExtReplace(*iFileName, "csv")
	}

	oFile, err := os.Create(*oFileName)
	if err != nil {
		log.Fatalf("Error in writing csv data to file %s because of: %q", oFile.Name(), err)
	}
	defer oFile.Close()
	interface_map.ToCSV(oFile)
	log.Infof("Saved to %s", oFile.Name())

	if *jsonOut {		
		jsonFileName := cisco_parser.FileExtReplace(*oFileName, "json")
		jsonFile, err := os.Create(jsonFileName)
		if err != nil {
			log.Fatalf("Error in writing json data to file %s because of: %q", jsonFile.Name(), err)
		}
		defer jsonFile.Close()
		interface_map.ToJSON(jsonFile)
		log.Infof("Saved to %s", jsonFileName)
	}
}