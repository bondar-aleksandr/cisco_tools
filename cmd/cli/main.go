package main

import (
	"flag"
	"log"
	"os"
	"time"
	"github.com/bondar-aleksandr/cisco_parser"
)

var (
	InfoLogger  *log.Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger *log.Logger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger *log.Logger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func main() {

	start := time.Now()
	InfoLogger.Println("Starting...")

	var iFileName = flag.String("i", "", "input configuration file to parse data from")
	var oFileName = flag.String("o", "", "output csv file, default is input filename with .csv extension")
	var devtype = flag.String("t", "ios", "cisco OS family, possible values are ios, nxos. Default is ios")
	var jsonOut = flag.Bool("j", false, "Whether JSON file needed. Default is false")

	if len(os.Args) < 2 {
		ErrorLogger.Fatalf("No input data provided, use -h flag for help. Exiting...")
	}
	flag.Parse()

	InfoLogger.Printf("Program started, got the following parameters:\ninput file: %s\noutput file: %s\ndevice type: %s\nJSON output: %v\n", *iFileName, *oFileName, *devtype, *jsonOut)

	iFile, err := os.Open(*iFileName)
	if err != nil {
		ErrorLogger.Fatalf("Can not open file %q because of: %q", *iFileName, err)
	}
	defer iFile.Close()

	interface_map, err := cisco_parser.ParseInterfaces(iFile, *devtype)
	if err != nil {
		ErrorLogger.Fatal(err)
	}	

	if *oFileName == "" {
		*oFileName = cisco_parser.FileExtReplace(*iFileName, "csv")
	}

	oFile, err := os.Create(*oFileName)
	if err != nil {
		ErrorLogger.Fatalf("Error in writing csv data to file %q because of: %q", oFile.Name(), err)
	}
	defer oFile.Close()
	interface_map.ToCSV(oFile)
	InfoLogger.Printf("Saved to %q\n", oFile.Name())

	if *jsonOut {		
		jsonFileName := cisco_parser.FileExtReplace(*oFileName, "json")
		jsonFile, err := os.Create(jsonFileName)
		if err != nil {
			ErrorLogger.Fatalf("Error in writing json data to file %q because of: %q", jsonFile.Name(), err)
		}
		defer jsonFile.Close()
		interface_map.ToJSON(jsonFile)
		InfoLogger.Printf("Saved to %q\n", jsonFileName)
	}
	InfoLogger.Printf("Finished! Time taken: %s\n", time.Since(start))
}