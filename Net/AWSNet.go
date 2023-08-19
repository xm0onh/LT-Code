package Net

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func Ec2IpExtractor(region, role, val string) []string {

	cmd := "aws"
	//role:="Role1"
	// args := []string{"ec2", "describe-instances", "--filter", "Name=tag:" + role + ",Values=" + val,
	//  "Name=instance-state-name, Values=running",
	//"--query", "Reservations[].Instances[*].[PrivateIpAddress]"}

	//args := []string{"ec2", "--region", "us-east-1", "describe-instances",
	args := []string{"ec2", "--region", region, "describe-instances",
		"--query", "Reservations[].Instances[*].PrivateIpAddress ", "--filters", "Name=tag:" + role + ",Values=" + val, "Name=instance-state-name, Values=running"}

	//args := []string{"ec2", "describe-instances", "--filter", "Name=tag:" + role + ",Values=" + val,
	//	"--query", "Reservations[].Instances[*].PrivateIpAddress[]"}
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	//sliceIP:=string(out)
	s := string(out)
	str := strings.Trim(s, "\"")
	//fmt.Println("Slice is",str)

	//text:="Ma text is"
	strSlice := StrExtract(str)
	//fmt.Println(" Slice is",strSlice)
	strIP := make([]string, 0, 100)
	for _, v := range strSlice {
		strIP = append(strIP, v[1:len(v)-1])

	}

	//fmt.Println(" Slice is",strSlice[0][1:len(strSlice[0])-1])
	fmt.Println(" Slice is", strIP)
	return strIP

}

func StrExtract(word string) []string {
	r, _ := regexp.Compile(`"[^"]*"`)
	result := r.FindAllString(word, -1)
	//RemoveDuplicates(&result)
	return (result)
}

func EC2IPsForAllRegions(regions []string, role, value string) []string {

	//regions := [...]string{"us-east-1", "ap-northeast-1"}
	//regions := [...]string{"us-east-1"}

	IPs := make([]string, 0, 10)
	for _, v := range regions {
		RegionalIPs := Ec2IpExtractor(v, role, value)
		for _, j := range RegionalIPs {
			IPs = append(IPs, j)

		}

	}
	return IPs
}

func EC2IPsAndIDSForAllRegions(regions []string, role, value string) []string {

	//regions := [...]string{"us-east-1", "ap-northeast-1"}
	//regions := [...]string{"us-east-1"}

	IPsAndIDs := make([]string, 0, 10)
	for _, v := range regions {
		RegionalIPs := Ec2IDandIPExtractor(v, role, value)
		fmt.Println("Regional IPsAndIds are", RegionalIPs)
		for _, j := range RegionalIPs {
			IPsAndIDs = append(IPsAndIDs, j)

		}

	}
	return IPsAndIDs
}

func GetIDs(regions []string) []string {
	//regions := [...]string{"us-east-1"}

	IDs := make([]string, 0, 10)
	fmt.Println("regions are", regions)
	for _, v := range regions {
		fmt.Println("Each regions is", v)
		RegionalIds := Ec2IDExtractor(v)
		for _, j := range RegionalIds {
			IDs = append(IDs, j)

		}

	}
	return IDs
}

func IfIamArequestor(requestorSlice []string, myIP string) bool {
	for _, IP := range requestorSlice {
		if myIP == IP {
			return true
		}
	}
	return false
}

func Ec2IDExtractor(region string) []string {

	cmd := "aws"
	//role:="Role1"

	//aws ec2 describe-instances --filter "Name=tag:Role2,Values=Root-nodes" --query "Reservations[].Instances[*].InstanceId[]"
	//aws ec2 describe-instances --query 'Reservations[*].Instances[*].[InstanceId]' --filters Name=instance-state-name,Values=running
	//    args := []string{"ec2", "--region", region, "describe-instances",
	//          "--query", "Reservations[].Instances[*].InstanceId[]", "--filters", "Name=instance-state-name, Values=running"}

	args := []string{"ec2", "--region", region, "describe-instances",
		"--query", "Reservations[].Instances[*].[InstanceId] ", "--filters", "Name=instance-state-name, Values=running"}
	//"--query", "Reservations[].Instances[*].[InstanceId] ", "--filters", "Name=tag:"+ 'role',  "Values=val" , "Name=instance-state-name, Values=running"}

	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	//sliceIP:=string(out)
	s := string(out)
	str := strings.Trim(s, "\"")
	//	fmt.Println("Slice is",str)

	//text:="Ma text is"
	strSlice := StrExtract(str)
	fmt.Println("ids are", strSlice)
	//if errstr!=nil{
	//      fmt.Println(errstr)
	//}
	//fmt.Println(" Slice is",strSlice)
	strIP := make([]string, 0, 100)
	for _, v := range strSlice {
		strIP = append(strIP, v[1:len(v)-1])

	}

	//fmt.Println(" Slice is",strSlice[0][1:len(strSlice[0])-1])
	//	fmt.Println(" Slice is", strIP)
	return strIP

}

func Ec2RequestorIDExtractor(region, role, val string) []string {

	cmd := "aws"
	//role:="Role1"

	//aws ec2 describe-instances --filter "Name=tag:Role2,Values=Root-nodes" --query "Reservations[].Instances[*].InstanceId[]"
	//aws ec2 describe-instances --query 'Reservations[*].Instances[*].[InstanceId]' --filters Name=instance-state-name,Values=running
	//    args := []string{"ec2", "--region", region, "describe-instances",
	//          "--query", "Reservations[].Instances[*].InstanceId[]", "--filters", "Name=instance-state-name, Values=running"}

	args := []string{"ec2", "--region", region, "describe-instances",
		"--query", "Reservations[].Instances[*].[InstanceId] ", "--filters", "Name=tag:" + role + ",Values=" + val, "Name=instance-state-name, Values=running"}
	//"--query", "Reservations[].Instances[*].[InstanceId] ", "--filters", "Name=tag:"+ 'role',  "Values=val" , "Name=instance-state-name, Values=running"}

	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	//sliceIP:=string(out)
	s := string(out)
	str := strings.Trim(s, "\"")
	//	fmt.Println("Slice is",str)

	//text:="Ma text is"
	strSlice := StrExtract(str)
	fmt.Println("ids are", strSlice)
	//if errstr!=nil{
	//      fmt.Println(errstr)
	//}
	//fmt.Println(" Slice is",strSlice)
	strIP := make([]string, 0, 100)
	for _, v := range strSlice {
		strIP = append(strIP, v[1:len(v)-1])

	}

	//fmt.Println(" Slice is",strSlice[0][1:len(strSlice[0])-1])
	//	fmt.Println(" Slice is", strIP)
	return strIP

}

func Ec2RequestorIDExtractorForAllRegions(regions []string, role, value string) []string {

	//regions := [...]string{"us-east-1", "ap-northeast-1"}
	//regions := [...]string{"us-east-1"}

	IDs := make([]string, 0, 10)
	for _, v := range regions {
		RegionalRequestorID := Ec2RequestorIDExtractor(v, role, value)
		for _, j := range RegionalRequestorID {
			IDs = append(IDs, j)

		}

	}
	return IDs
}

func GetmyID() string {
	cmd := "wget"
	args := []string{"-q", "-O", "-", "http://instance-data/latest/meta-data/instance-id"}
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return string(out)
}

func IPaddress() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you sure that you are connected to the network?")
}

func Ec2IDandIPExtractor(region string, role, val string) []string {

	cmd := "aws"
	//role:="Role1"

	//aws ec2 describe-instances --filter "Name=tag:Role2,Values=Root-nodes" --query "Reservations[].Instances[*].InstanceId[]"
	//aws ec2 describe-instances --query 'Reservations[*].Instances[*].[InstanceId]' --filters Name=instance-state-name,Values=running
	//    args := []string{"ec2", "--region", region, "describe-instances",
	//          "--query", "Reservations[].Instances[*].InstanceId[]", "--filters", "Name=instance-state-name, Values=running"}
	fmt.Println("Just Before Query")
	args := []string{"ec2", "--region", region, "describe-instances",
		"--query", "Reservations[].Instances[*].[InstanceId, PrivateIpAddress] ", "--filters", "Name=tag:" + role + ",Values=" + val, "Name=instance-state-name, Values=running"}
	//"--query", "Reservations[].Instances[*].[InstanceId] ", "--filters", "Name=tag:"+ 'role',  "Values=val" , "Name=instance-state-name, Values=running"}
	fmt.Println("Just after Query")

	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	//sliceIP:=string(out)
	s := string(out)
	str := strings.Trim(s, "\"")
	fmt.Println("Slice is", str)

	//text:="Ma text is"
	strSlice := StrExtract(str)
	//	fmt.Println("ids are",strSlice)
	//if errstr!=nil{
	//      fmt.Println(errstr)
	//}
	//fmt.Println(" Slice is",strSlice)
	strIP := make([]string, 0, 100)
	for _, v := range strSlice {
		strIP = append(strIP, v[1:len(v)-1])

	}

	//fmt.Println(" Slice is",strSlice[0][1:len(strSlice[0])-1])
	//	fmt.Println(" Slice is", strIP)
	return strIP

}

//	ec2 describe-instances --query Reservations[].Instances[*].{InstanceId, PrivateIpAddress} --filters Name=instance-state-name, Values=running
