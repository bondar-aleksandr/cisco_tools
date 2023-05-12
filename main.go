package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"sort"
	"strings"
	log "github.com/sirupsen/logrus"
	"reflect"
)

type CiscoInterface struct {
	Name string
	Description string
	Ip_addr string
	Subnet string
	Vrf string
	ACLin string
	ACLout string
}

func (c CiscoInterface) ToSlice() []string {
	return []string{c.Name, c.Description, c.Ip_addr, c.Subnet, c.Vrf, c.ACLin, c.ACLout}
}

type CiscoInterfaceMap map[string]*CiscoInterface

func (c CiscoInterfaceMap) GetSortedKeys() []string {
	keys := make([]string,0)
	for k := range c {
		keys = append(keys,k)
	}
	sort.Strings(keys)
	return keys
}

func (c CiscoInterfaceMap) GetFields() []string {
	fields := reflect.VisibleFields(reflect.TypeOf(CiscoInterface{}))
	result := []string{}
	for _,v := range(fields) {
		result = append(result, v.Name)
	}
	return result
}

const(
	INTF_REGEXP = `^interface (\S+)`
	DESC_REGEXP = ` {1,2}description (.*)$`
	IP_REGEXP = ` {1,2}ip(?:v4)? address (\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}) (\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})(?: secondary)?`
	VRF_REGEXP = ` {1,2}vrf(?: forwarding| member)? (\S+)`
	ACLIN_REGEXP = ` {1,2}access-group (\S+) in`
	ACLOUT_REGEXP = ` {1,2}access-group (\S+) out`
)

var (
	intf_compiled = regexp.MustCompile(INTF_REGEXP)
	desc_compiled = regexp.MustCompile(DESC_REGEXP)
	ip_compiled = regexp.MustCompile(IP_REGEXP)
	vrf_compiled = regexp.MustCompile(VRF_REGEXP)
	aclin_compiled = regexp.MustCompile(ACLIN_REGEXP)
	aclout_compiled = regexp.MustCompile(ACLOUT_REGEXP)
)

func main() {
	var ifile = flag.String("i", "", "input configuration file to parse data from")
	var ofile = flag.String("o", "", "output csv file")
	var devtype = flag.String("t", "ios", "cisco OS family, possible values are ios, nxos. Default is ios")

	flag.Parse()
	log.Infof("Program started, got the following parameters: input file: %s, output file: %s, device type: %s", *ifile, *ofile, *devtype)

	f, err := os.Open(*ifile)
	if err != nil {
		log.Fatalf("Can not open file %s because of: %q", *ifile, err)
	}
	defer f.Close()
	interface_map := parsing(f, *devtype)
	// for k,v := range(interface_map) {
	// 	fmt.Printf("%s: %+v\n", k,v)
	// }
	ToCSV(interface_map, *ofile)
}

func getIP(s string, d string) (ip_addr, subnet string) {
	
	if strings.Contains(s, "dhcp") {
		return "dhcp", "dhcp"
	}
	
	if d == "ios" {

		ip_str := ip_compiled.FindStringSubmatch(s)[1]
		mask_str := ip_compiled.FindStringSubmatch(s)[2]

		ip := net.ParseIP(ip_str).To4()
		mask := net.IPMask(net.ParseIP(mask_str).To4())
		mask_cidr, _ := mask.Size()
		net_addr := ip.Mask(mask)
		ip_cidr := fmt.Sprintf("%s/%v", ip.String(), mask_cidr)
		prefix := fmt.Sprintf("%s/%v", net_addr.String(), mask_cidr)

		return ip_cidr, prefix

	} else if d == "nxos" {
		ip_str := regexp.MustCompile(` {2}ip address (\S+)`).FindStringSubmatch(s)[1]
		_, prefix, _ := net.ParseCIDR(ip_str)
		return ip_str, prefix.String()
	}
	return
}

func parsing(f *os.File, d string) CiscoInterfaceMap {

	interfaces := CiscoInterfaceMap{}
	var intf_name string

	line_separator := "!"
	line_ident := " "

	if d == "nxos" {
		line_separator = ""
		line_ident = "  "
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		line := scanner.Text()
		fmt.Println(line)

		if strings.HasPrefix(scanner.Text(),`interface `) {

			intf_name = intf_compiled.FindStringSubmatch(scanner.Text())[1]
			interfaces[intf_name] = &CiscoInterface{Name: intf_name}

		} else if strings.HasPrefix(scanner.Text(), line_ident) && len(interfaces) > 0 {

			if match, _ := regexp.Match(DESC_REGEXP, scanner.Bytes()); match {
				intf_desc := desc_compiled.FindStringSubmatch(scanner.Text())[1]
				interfaces[intf_name].Description = intf_desc

			} else if strings.HasPrefix(scanner.Text(), fmt.Sprintf("%sip address", line_ident)) || strings.HasPrefix(scanner.Text(), fmt.Sprintf("%sipv4 address", line_ident)) {
			
					ip_cidr, prefix := getIP(scanner.Text(), d)
					interfaces[intf_name].Ip_addr = ip_cidr
					interfaces[intf_name].Subnet = prefix

			} else if match, _ := regexp.Match(VRF_REGEXP, scanner.Bytes()); match {
				vrf := vrf_compiled.FindStringSubmatch(scanner.Text())[1]
				interfaces[intf_name].Vrf = vrf

			} else if match, _ := regexp.Match(ACLIN_REGEXP, scanner.Bytes()); match {
				aclin := aclin_compiled.FindStringSubmatch(scanner.Text())[1]
				interfaces[intf_name].ACLin = aclin

			} else if match, _ := regexp.Match(ACLOUT_REGEXP, scanner.Bytes()); match {
				aclout := aclout_compiled.FindStringSubmatch(scanner.Text())[1]
				interfaces[intf_name].ACLout = aclout
			}

		} else if !(strings.HasPrefix(scanner.Text(), line_separator) || strings.HasPrefix(scanner.Text(), `interface`)) && len(interfaces) > 0 {
			break
		}
	}
	log.Infof("parsing finished, got %v interfaces", len(interfaces))
	return interfaces
}

func ToCSV(intf_map CiscoInterfaceMap, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal("Error in writing csv data to file:", err)
	}
	w := csv.NewWriter(f)
	headers := intf_map.GetFields()
	w.Write(headers)

	for _,v := range intf_map.GetSortedKeys() {
		line := intf_map[v].ToSlice()
		w.Write(line)
	}
	w.Flush()
	log.Infof("Writing CSV to %s done", filename)
}