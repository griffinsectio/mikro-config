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

func PrintAddressesToScreen(client *routeros.Client) {
	main_command := []string{"/ip/address/print"}
	RunCommand(client, main_command)

	reply, err := client.RunArgs(main_command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	reply_list := reply.Re

	for i := 0; i < len(reply_list); i++ {
		curr_iface := reply_list[i].Map
		iface_name := curr_iface["interface"]
		network := curr_iface["network"]
		ip_address := curr_iface["address"]

		fmt.Printf("%d. %s %s %s\n", i, iface_name, network, ip_address)
	}
}

func PrintAddresses(client *routeros.Client) {
	ClearScreen()

	PrintAddressesToScreen(client)

	EnterToContinue()
}

func AddAdress(client *routeros.Client) {
	ifaces := RetrieveInterfaces(client)
	main_command := "/ip/address/add"
	var ip_address string
	var iface string

	// TODO:
	// Validate ip_address and the subnet

	ip_address_prompt := "Provide the ip address and the prefix: "
	ip_address = GetStringUserInput(ip_address_prompt)
	ip_address = fmt.Sprintf("=address=%s", ip_address)

	PrintInterfacesToScreen(client)
	iface_prompt := "Select interface you want to add the address to: "
	iface = ifaces[GetIntUserInput(iface_prompt)-1]
	iface = fmt.Sprintf("=interface=%s", iface)

	full_command := []string{main_command, ip_address, iface}

	RunCommand(client, full_command)

	EnterToContinue()
}

func RemoveAddress(client *routeros.Client) {
	ClearScreen()

	main_command := "/ip/address/remove"
	var ip_address_num string

	PrintAddressesToScreen(client)
	ip_address_num_prompt := "Select the ip address number to delete: "
	ip_address_num = GetStringUserInput(ip_address_num_prompt)

	ip_address_num = fmt.Sprintf("=numbers=%s", ip_address_num)

	full_command := []string{main_command, ip_address_num}

	RunCommand(client, full_command)

	EnterToContinue()
}

func IpConfig(client *routeros.Client) {
	ClearScreen()

	menu := []string{
		"Print Addresses",
		"Add Address",
		"Remove Address",
		"Back",
	}

	PrintOptions(menu)
	choose_prompt := "Choose the configuration you want to do: "
	user_choice := GetIntUserInput(choose_prompt)

	// TODOs
	switch user_choice {
	case 1:
		// 1. Print addresses
		PrintAddresses(client)
	case 2:
		// 2. Add address to interface
		AddAdress(client)
	case 3:
		// 3. Remove address of an interface
		RemoveAddress(client)
	case len(menu):
		return
	default:
		fmt.Println("Please select the available options!")
		os.Exit(1)

	}
}

func RoutingConfig(client *routeros.Client) {
	// TODOs
	// 1. Print all routing protocol options
	// 2. BGP
	// 3. OSPF
	// 4. RIP
}

func PrintIdentity(client *routeros.Client) {
	ClearScreen()
	main_command := "/system/identity/print"

	full_command := []string{main_command}

	RunCommand(client, full_command)

	EnterToContinue()
}

func SetIdentity(client *routeros.Client) {
	ClearScreen()
	// /system/identity/set name=daisy
	main_command := "/system/identity/set"
	var name string

	name_prompt := "Give the new identity: "
	name = GetStringUserInput(name_prompt)
	name = fmt.Sprintf("=name=%s", name)

	full_command := []string{main_command, name}

	RunCommand(client, full_command)

	EnterToContinue()
}

func RetrieveUsersData(client *routeros.Client) map[string][]string {
	// /user/print
	main_command := "/user/print"
	full_command := []string{main_command}

	reply, err := client.RunArgs(full_command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	reply_list := reply.Re

	users_data := make(map[string][]string)
	var users_arr []string
	var groups_arr []string

	for _, element := range reply_list {
		element_map := element.Map
		users_arr = append(users_arr, element_map["name"])
		groups_arr = append(users_arr, element_map["group"])
	}

	users_data["users"] = users_arr
	users_data["groups"] = groups_arr

	return users_data

}

func RetrieveGroups(client *routeros.Client) []string {
	var groups []string

	main_command := "/user/group/print"
	full_command := []string{main_command}

	reply, err := client.RunArgs(full_command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	reply_list := reply.Re

	for _, element := range reply_list {
		element_map := element.Map
		groups = append(groups, element_map["name"])
	}

	return groups
}

func PrintUsersToScreen(client *routeros.Client) {
	main_command := "/user/print"
	full_command := []string{main_command}

	reply, err := client.RunArgs(full_command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	reply_list := reply.Re

	for i, element := range reply_list {
		element_map := element.Map
		fmt.Printf("%d. %s %s\n", i+1, element_map["name"], element_map["group"])
	}
}

func PrintUsers(client *routeros.Client) {
	ClearScreen()

	PrintUsersToScreen(client)

	EnterToContinue()
}

func AddUser(client *routeros.Client) {
	ClearScreen()

	main_command := "/user/add"
	var name string
	groups := RetrieveGroups(client)
	var group string
	var password string

	name_prompt := "Give the name for the new user: "
	name = GetStringUserInput(name_prompt)
	name = fmt.Sprintf("=name=%s", name)

	PrintOptions(groups)
	group_prompt := "Group for the new user: "
	group = groups[GetIntUserInput(group_prompt)-1]
	group = fmt.Sprintf("=group=%s", group)

	password_prompt := "Password for the new user: "
	password = GetStringUserInput(password_prompt)
	password = fmt.Sprintf("=password=%s", password)

	full_command := []string{main_command, name, group, password}

	RunCommand(client, full_command)

	EnterToContinue()
}

func RemoveUser(client *routeros.Client) {
	// /user/remove numbers=
	ClearScreen()
	main_command := "/user/remove"
	users := RetrieveUsersData(client)["users"]
	var name string

	name_prompt := "The the user you want to remove: "
	name = users[GetIntUserInput(name_prompt)-1]
	name = fmt.Sprintf("=numbers=%s", name)

	full_command := []string{main_command, name}

	RunCommand(client, full_command)

	EnterToContinue()
}

// Apparently MikroTik doesn't support editing the user name
func EditUserName(client *routeros.Client) {
	// /user/set numbers= name=
	ClearScreen()
	remove_user := "/user/remove"
	add_user := "/user/add"

	var username string
	var new_username string

	users := RetrieveUsersData(client)["users"]
	username_prompt := "The user you to change the username: "
	username = users[GetIntUserInput(username_prompt)-1]
	username = fmt.Sprintf("=numbers=%s", username)

	new_username_prompt := "The username you want to change it to: "
	new_username = GetStringUserInput(new_username_prompt)
	new_username = fmt.Sprintf("==%s", new_username)

	// Remove user
	remove_command := []string{remove_user, username}
	RunCommand(client, remove_command)

	// Add user
	add_command := []string{add_user, new_username}
	RunCommand(client, add_command)

	EnterToContinue()
}

func EditUserGroup(client *routeros.Client) {
	// /user/set =numbers= =group=
	ClearScreen()

	main_command := "/user/set"
	users := RetrieveUsersData(client)["users"]
	groups := RetrieveGroups(client)
	var user string
	var group string

	PrintUsersToScreen(client)
	user_prompt := "The user you want to edit: "
	user = users[GetIntUserInput(user_prompt)-1]
	user = fmt.Sprintf("=numbers=%s", user)

	PrintOptions(groups)
	group_prompt := "The group you want to give for the user: "
	group = groups[GetIntUserInput(group_prompt)-1]
	group = fmt.Sprintf("=group=%s", group)

	full_command := []string{main_command, user, group}

	RunCommand(client, full_command)

	EnterToContinue()
}

func EditUserPassword(client *routeros.Client) {
	// /user/set =numbers= =password=
	ClearScreen()

	main_command := "/user/set"
	users := RetrieveUsersData(client)["users"]

	var user string
	var password string

	PrintUsersToScreen(client)
	user_prompt := "The user you want to change the password"
	user = users[GetIntUserInput(user_prompt)-1]
	user = fmt.Sprintf("=numbers=%s", user)

	password_prompt := "New password: "
	password = GetStringUserInput(password_prompt)
	password = fmt.Sprintf("=password=%s", password)

	full_command := []string{main_command, user, password}

	RunCommand(client, full_command)

	EnterToContinue()
}

func EditUser(client *routeros.Client) {
	ClearScreen()

	menu := []string{
		"Edit user name",
		"Edit user group",
		"Edit user password",
		"Back",
	}

	PrintOptions(menu)
	user_choice_prompt := "Choose the configuration you want to do: "
	user_choice := GetIntUserInput(user_choice_prompt)

	switch user_choice {
	case 1:
		fmt.Println("Work in progress")
		EnterToContinue()
	case 2:
		EditUserGroup(client)
	case 3:
		EditUserPassword(client)
	case len(menu):
		return
	}
	// TODOS
	// 1. Edit user name
	// 2. Change user group
	// 3. Edit user password
}

func UserManagement(client *routeros.Client) {
	ClearScreen()

	menu := []string{
		"Print Users",
		"Create User",
		"Delete User",
		"Edit User",
		"Back",
	}

	PrintOptions(menu)
	choose_prompt := "Choose the configuration you want to do: "
	user_choice := GetIntUserInput(choose_prompt)

	switch user_choice {
	case 1:
		PrintUsers(client)
	case 2:
		AddUser(client)
	case 3:
		RemoveUser(client)
	case 4:
		EditUser(client)
	case len(menu):
		return

	}
	// TODOS
	// 1. Print users
	// 2. Create user
	// 3. Delete user
	// 4. Edit user
}

func Reboot(client *routeros.Client) {
	ClearScreen()

	main_command := "/system/reboot"

	full_command := []string{main_command}

	RunCommand(client, full_command)

	EnterToContinue()
}

func Shutdown(client *routeros.Client) {
	ClearScreen()

	main_command := "/system/shutdown"

	full_command := []string{main_command}

	RunCommand(client, full_command)

	EnterToContinue()

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

	PrintOptions(menu)
	choose_prompt := "Choose the configuration you want to do: "
	user_choice := GetIntUserInput(choose_prompt)

	switch user_choice {
	case 1:
		PrintIdentity(client)
	case 2:
		SetIdentity(client)
	case 3:
		UserManagement(client)
	case 4:
		Reboot(client)
	case 5:
		Shutdown(client)
	case len(menu):
		return
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
