package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/go-routeros/routeros/v3"
)

func ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func PrintOptions(options []string) {
	for i := 0; i < len(options); i++ {
		fmt.Printf("%d. %s\n", i+1, options[i])
	}
}

func EnterToContinue() {
	fmt.Println("Press 'enter' or 'return' to continue")
	fmt.Scanln()
}

func RunCommand(client *routeros.Client, command []string) {
	reply, err := client.RunArgs(command)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Router reply: %s\n", reply.String())
	}
}

func GetIntUserInput(prompt string) int {
	fmt.Print(prompt)
	var temp string
	fmt.Scanln(&temp)
	user_choice, err := strconv.Atoi(temp)
	if err != nil {
		fmt.Println("Please enter a valid number!")
		os.Exit(1)
	}
	return user_choice
}

func GetStringUserInput(prompt string) string {
	fmt.Print(prompt)
	in := bufio.NewReader(os.Stdin)
	str, err := in.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	str_arr := strings.Split(str, "")
	str_arr = str_arr[:len(str_arr)-1]

	str = strings.Join(str_arr, "")

	return str
}

func LoginMikroTik() *routeros.Client {
	var address string
	port := "8728"

	var username string
	var password string

	fmt.Print("Address of MikroTik and API port (default = 8728): ")
	fmt.Scanln(&address)

	addresses_arr := strings.Split(address, ":")

	ip := net.ParseIP(addresses_arr[0])
	if ip == nil {
		ClearScreen()

		fmt.Println("Please provide a valid IP address!")
		os.Exit(1)
	}

	if len(addresses_arr) == 1 {
		address = address + ":" + port
	}

	fmt.Print("Username to login: ")
	fmt.Scanln(&username)

	fmt.Print("Password for the user: ")
	fmt.Scanln(&password)

	client, err := routeros.Dial(address, username, password)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return client
}

func RetrieveInterfaces(client *routeros.Client) []string {
	ifaces := []string{}
	command := []string{"/interface/print"}
	reply, err := client.RunArgs(command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	reply_list := reply.Re
	for i := 0; i < len(reply_list); i++ {
		curr_iface := reply_list[i].List
		ifaces = append(ifaces, curr_iface[1].Value)
	}

	return ifaces
}

func PrintInterfacesToScreen(client *routeros.Client) {
	command := []string{"/interface/print"}
	reply, err := client.RunArgs(command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	reply_list := reply.Re

	for i := 0; i < len(reply_list); i++ {
		curr_iface := reply_list[i].Map

		if curr_iface["comment"] != "" {
			fmt.Printf(";;; %s\n", curr_iface["comment"])
		}

		fmt.Printf("%d. %s\n", i+1, curr_iface["name"])
	}
}

func PrintInterfaces(client *routeros.Client) {
	ClearScreen()

	PrintInterfacesToScreen(client)

	EnterToContinue()
}

func ChangeInterfaceName(client *routeros.Client) {
	ClearScreen()

	// /interface/set numbers=ether1-lagilagi name=ether1
	main_command := "/interface/set"
	ifaces := RetrieveInterfaces(client)
	var iface string
	var name string

	PrintInterfacesToScreen(client)
	iface_prompt := "Choose the interface you want to change the name: "
	iface = ifaces[GetIntUserInput(iface_prompt)-1]
	iface = fmt.Sprintf("=numbers=%s", iface)

	name_prompt := "The new name: "
	name = GetStringUserInput(name_prompt)
	name = fmt.Sprintf("=name=%s", name)

	full_command := []string{main_command, iface, name}

	RunCommand(client, full_command)

	EnterToContinue()
}

func AddInterfaceComment(client *routeros.Client) {
	ClearScreen()

	main_command := "/interface/set"

	ifaces := RetrieveInterfaces(client)
	var iface string
	var comment string

	PrintInterfacesToScreen(client)
	iface_prompt := "Choose the interface you want to add comment to: "
	iface = ifaces[GetIntUserInput(iface_prompt)-1]
	iface = fmt.Sprintf("=numbers=%s", iface)

	fmt.Println()

	comment_prompt := "Comment you want to add: "
	comment = GetStringUserInput(comment_prompt)

	comment = fmt.Sprintf("=comment=%s", comment)

	full_command := []string{main_command, iface, comment}

	RunCommand(client, full_command)

	EnterToContinue()
}

func RemoveInterfaceComment(client *routeros.Client) {
	ClearScreen()

	var ifaces = RetrieveInterfaces(client)
	main_command := "/interface/set"
	var iface string
	var comment string

	PrintInterfacesToScreen(client)
	iface_prompt := "Interface you want to remove comment from: "
	iface = ifaces[GetIntUserInput(iface_prompt)-1]
	iface = fmt.Sprintf("=numbers=%s", iface)

	comment = "=comment="

	full_command := []string{main_command, iface, comment}

	RunCommand(client, full_command)

	EnterToContinue()
}

func AddVlan(client *routeros.Client) {
	ClearScreen()

	main_command := "/interface/vlan/add"

	ifaces := RetrieveInterfaces(client)
	var vlan_name string
	var vlan_id int
	var vlan_id_str string
	var iface string

	vlan_name = GetStringUserInput("Give the vlan name: ")
	vlan_id = GetIntUserInput("Give the vlan ID: ")
	PrintInterfacesToScreen(client)
	iface = ifaces[GetIntUserInput("Choose the interface: ")-1]

	vlan_name = fmt.Sprintf("=name=%s", vlan_name)
	vlan_id_str = fmt.Sprintf("=vlan-id=%d", vlan_id)
	iface = fmt.Sprintf("=interface=%s", iface)

	full_command := []string{main_command, vlan_name, vlan_id_str, iface}

	RunCommand(client, full_command)

	EnterToContinue()
}

func RemoveVlan(client *routeros.Client) {
	ClearScreen()

	main_command := "/interface/vlan/remove"

	var vlan_name string

	vlan_name = GetStringUserInput("Give the vlan name: ")

	vlan_name = fmt.Sprintf("=numbers=%s", vlan_name)

	full_command := []string{main_command, vlan_name}

	RunCommand(client, full_command)

	EnterToContinue()
}

func InterfaceConfig(client *routeros.Client) {
	ClearScreen()

	menu := []string{
		"Print Interfaces",
		"Change Interface Name",
		"Add/Edit Comment",
		"Remove Comment",
		"Add VLAN",
		"Remove VLAN",
		"Back",
	}

	PrintOptions(menu)
	choose_prompt := "Choose the configuration you want to do: "
	user_choice := GetIntUserInput(choose_prompt)

	// TODOs
	switch user_choice {
	case 1:
		// 1. Print interfaces
		PrintInterfaces(client)
	case 2:
		// 2. Change interface name
		ChangeInterfaceName(client)
	case 3:
		// 3. Add comment to interface
		AddInterfaceComment(client)
	case 4:
		// 4. Remove comment from interface
		RemoveInterfaceComment(client)
	case 5:
		// 5. Add vlan
		AddVlan(client)
	case 6:
		// 6. Remove vlan
		RemoveVlan(client)
	case len(menu):
		return
	default:
		fmt.Println("Please select the available options!")
		os.Exit(1)
	}
}

func IpConfig(client *routeros.Client) {
	// TODOs
	// 1. Print addresses
	// 2. Add address to interface
	// 3. Remove address of an interface
}

func RoutingConfig(client *routeros.Client) {
	// TODOs
	// 1. Print all routing protocol options
	// 2. BGP
	// 3. OSPF
	// 4. RIP
}

func SystemConfig(client *routeros.Client) {
	ClearScreen()

	menu := []string{
		"Print Identity",
		"Set Identity",
		"User Management",
		"Reboot",
		"Shutdown",
		"Back",
	}
	// TODOs

	// 1. Print system identity
	// 2. edit system identity
	// 3. User management
	// 4. Reboot router
	// 5. Shutdown router
}

func PrintVlan(client *routeros.Client) {
	command := []string{"/interface/print"}
	res, err := client.RunArgs(command)
	if err != nil {
		os.Exit(1)
	}

	reply_arr := res.Re

	for i := 0; i < len(reply_arr); i++ {
		fmt.Println(string(i+1) + res.Re[i].Map["name"])
	}
}

func main() {
	ClearScreen()

	is_running := true

	client := LoginMikroTik()
	defer client.Close()

	// PrintVlan(client)
	// os.Exit(0)

	for is_running {
		ClearScreen()

		fmt.Println("Hello there!")

		menu := []string{"interface", "ip", "routing", "system", "exit"}
		PrintOptions(menu)

		choose_prompt := "Select one of the options: "
		user_choice := GetIntUserInput(choose_prompt)

		switch user_choice {
		case 1:
			InterfaceConfig(client)
		case 2:
			IpConfig(client)
		case 3:
			RoutingConfig(client)
		case 4:
			SystemConfig(client)
		case len(menu):
			is_running = false
		default:
			fmt.Println("Please select the available options!")
			os.Exit(1)
		}
	}
}
