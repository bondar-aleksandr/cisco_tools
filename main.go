package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"net"
	"os"
	"regexp"
	log "github.com/sirupsen/logrus"
	// "reflect"
)

type CiscoInterface struct {
	Name string
	Description string
	Ip_addr string
	Subnet string
	Vrf string
}

func (c CiscoInterface) ToSlice() []string {
	return []string{c.Name, c.Description, c.Ip_addr, c.Subnet, c.Vrf}
}

const(
	INTF_REGEXP = `^interface (\S+)`
	DESC_REGEXP = ` description (.*)$`
	IP_REGEXP = ` ip(?:v4)? address (\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}) (\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})(?: secondary)?`
	VRF_REGEXP = ` vrf(?: forwarding)? (\S+)`
)

func main() {
	log.Info("program started")

	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal("Error: ", err)
		}
		defer f.Close()
		log.Infof("got the file: %s", f.Name())
		interface_map := parsing(f)
		// for k,v := range interface_map {							// for debug
		// 	fmt.Printf("interface: %s, content: %+v\n", k, *v)
		// }
		ToCSV(interface_map, "output.csv")

	} else {
		log.Fatal("No filename specified!")
	}
}

func parsing(f *os.File) map[string]*CiscoInterface {

	intf_compiled := regexp.MustCompile(INTF_REGEXP)
	desc_compiled := regexp.MustCompile(DESC_REGEXP)
	ip_compiled := regexp.MustCompile(IP_REGEXP)
	vrf_compiled := regexp.MustCompile(VRF_REGEXP)

	interfaces := map[string]*CiscoInterface{}
	var intf_name string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// line := scanner.Text()	// for debug
		// fmt.Println(line)
		if match, _ := regexp.Match(`^interface `, scanner.Bytes()); match {

			intf_name = intf_compiled.FindStringSubmatch(scanner.Text())[1]
			interfaces[intf_name] = &CiscoInterface{Name: intf_name}

		} else if match, _ := regexp.Match(`^ `, scanner.Bytes()); match && len(interfaces) > 0 {

			if match, _ := regexp.Match(DESC_REGEXP, scanner.Bytes()); match {
				intf_desc := desc_compiled.FindStringSubmatch(scanner.Text())[1]
				interfaces[intf_name].Description = intf_desc

			} else if match, _ := regexp.Match(IP_REGEXP, scanner.Bytes()); match {
				ip_str := ip_compiled.FindStringSubmatch(scanner.Text())[1]
				mask_str := ip_compiled.FindStringSubmatch(scanner.Text())[2]

				ip := net.ParseIP(ip_str).To4()
				mask := net.IPMask(net.ParseIP(mask_str).To4())
				mask_cidr, _ := mask.Size()
				net_addr := ip.Mask(mask)
				ip_cidr := fmt.Sprintf("%s/%v", ip.String(), mask_cidr)
				prefix := fmt.Sprintf("%s/%v", net_addr.String(), mask_cidr)
				interfaces[intf_name].Ip_addr = ip_cidr
				interfaces[intf_name].Subnet = prefix

			} else if match, _ := regexp.Match(VRF_REGEXP, scanner.Bytes()); match {
				vrf := vrf_compiled.FindStringSubmatch(scanner.Text())[1]
				interfaces[intf_name].Vrf = vrf
			}

		} else if match, _ := (regexp.Match(`^!|^interface`, scanner.Bytes())); !match && len(interfaces) > 0 {
			log.Info("parsing finished")
			break
		}
	}
	return interfaces
}

func ToCSV(intf_map map[string]*CiscoInterface, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal("Error in writing csv data to file:", err)
	}
	w := csv.NewWriter(f)
	headers := []string{"name", "description", "ip_addr", "subnet", "vrf"}
	w.Write(headers)
	for _,v := range(intf_map) {
		line := v.ToSlice()
		w.Write(line)
	}
	w.Flush()
	log.Info("Writing to CSV done")
}