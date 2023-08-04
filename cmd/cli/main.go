package main

import (
	"github.com/bondar-aleksandr/ios-config-parsing/parser"
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	var ifile = flag.String("i", "", "input configuration file to parse data from")
	var ofile = flag.String("o", "", "output csv file, default is input filename with .csv extension")
	var devtype = flag.String("t", "ios", "cisco OS family, possible values are ios, nxos. Default is ios")
	var jsonOut = flag.Bool("j", false, "Whether JSON file needed. Default is false")

	flag.Parse()
	log.Infof("Program started, got the following parameters: input file: %s, output file: %s, device type: %s, JSON output: %v", *ifile, *ofile, *devtype, *jsonOut)

	f, err := os.Open(*ifile)
	if err != nil {
		log.Fatalf("Can not open file %s because of: %q", *ifile, err)
	}
	defer f.Close()

	interface_map := parser.Parsing(f, *devtype)	

	if *ofile == "" {
		*ofile = parser.FileExtReplace(*ifile, "csv")
	}

	parser.ToCSV(interface_map, *ofile)

	if *jsonOut {						// Optional step to store json data needed in testing
		interface_map.ToJSON(*ifile)
	}
}