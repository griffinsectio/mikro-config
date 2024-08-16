package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/go-routeros/routeros/v3"
)

func PrintMenu(menu []string) {
	for i := 0; i < len(menu); i++ {
		fmt.Printf("%d. %s\n", i+1, menu[i])
	}
}

func GetUserInput() int {
	fmt.Print("Select one of the options: ")
	var temp string
	fmt.Scanln(&temp)
	user_choice, err := strconv.Atoi(temp)
	if err != nil {
		fmt.Println("Please enter a valid number!")
		os.Exit(1)
	}
	return user_choice
}

func LoginMikroTik() *routeros.Client {
	var address string
	var username string
	var password string

	fmt.Print("Address of MikroTik and API port (example: 192.168.1.1:8728): ")
	fmt.Scanln(&address)
	fmt.Println()

	fmt.Print("Username to login: ")
	fmt.Scanln(&username)
	fmt.Println()

	fmt.Print("Password for the user: ")
	fmt.Scanln(&password)
	fmt.Println()

	client, err := routeros.Dial(address, username, password)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return client
}

func PrintInterfaces(client *routeros.Client) {
	command := []string{"/interface/print"}
	reply, err := client.RunArgs(command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	reply_list := reply.Re
	// fmt.Println(reply_list)
	for i := 0; i < len(reply_list); i++ {
		curr_interface := reply_list[i].List
		// Print the comment for current interface (if exist)
		for j := 0; j < len(curr_interface); j++ {
			if curr_interface[j].Key == "comment" {
				fmt.Printf(";;;%s\n", curr_interface[j].Value)
			}
		}
		fmt.Println(curr_interface[1].Value)
	}
}

func AddInterfaceComment(client *routeros.Client) {
	PrintInterfaces(client)
	command := "/interface/comment"

	var iface string
	var comment string

	fmt.Print("Interface you want to add comment to: ")
	fmt.Scanln(&iface)
	fmt.Println()
	iface = fmt.Sprintf("=numbers=%s", iface)

	fmt.Println("Comment you want to add: ")
	temp := bufio.NewReader(os.Stdin)
	comment, err := temp.ReadString('\n')

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	comment = fmt.Sprintf("=comment=\"%s\"", comment)

	full_command := []string{command, iface, comment}

	client.RunArgs(full_command)
}

func RemoveInterfaceComment(client *routeros.Client) {
	PrintInterfaces(client)
	command := "/interface/comment"
	var inter_face string

	fmt.Print("Interface you want to add comment to: ")
	fmt.Scanln(&inter_face)
	fmt.Println()

	command = fmt.Sprintf("%s %s comment=\"\"", command, inter_face)
	client.RunArgs([]string{command})

}

func InterfaceConfig(client *routeros.Client) {
	menu := []string{"Print Interfaces", "Add Comment", "Remove Comment", "Add VLAN", "Remove VLAN"}

	PrintMenu(menu)
	user_choice := GetUserInput()

	// TODOs
	switch user_choice {
	case 1:
		// 1. Print interfaces
		PrintInterfaces(client)
	case 2:
		// 2. Add comment to interface
		AddInterfaceComment(client)
	default:
		fmt.Println("Please select the available options!")
		os.Exit(1)
	}

	// 3. Remove comment from interface

	// 4. Add vlan
	// 5. Remove vlan
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
	// 2. Add configuration for each routing protocol options
}

func SystemConfig(client *routeros.Client) {
	// TODOs
	// 1. Print and edit system identity
	// 2. Reboot router
	// 3. Shutdown router
}

func main() {
	fmt.Println("Hello there!")

	client := LoginMikroTik()
	defer client.Close()

	menu := []string{"interface", "ip", "routing", "system"}
	PrintMenu(menu)

	fmt.Print("Select one of the options: ")

	user_choice := GetUserInput()

	switch user_choice {
	case 1:
		InterfaceConfig(client)
	case 2:
		IpConfig(client)
	case 3:
		RoutingConfig(client)
	case 4:
		SystemConfig(client)
	default:
		fmt.Println("Please select the available options!")
		os.Exit(1)
	}
}
