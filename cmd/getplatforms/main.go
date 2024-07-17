package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/davidh-cyberark/privilegeaccessmanager-sdk-go/pam"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

/*
Create a file, creds.toml with these parameters and fill in your values
idtenanturl = "https://YOUR-TENANT.id.cyberark.cloud"
pcloudurl = "https://YOUR-SUBDOMAIN.privilegecloud.cyberark.cloud"
user = "PAM_SERVICE_ACCOUNT_USER"
pass = "PAM_SERVICE_ACCOUNT_USER password"
*/
func main() {
	k := koanf.New(".")
	err := k.Load(file.Provider("creds.toml"), toml.Parser())
	if err != nil {
		panic(fmt.Sprintf("failed to load creds.toml: %s", err.Error()))
	}

	config := pam.NewConfig(k.String("idtenanturl"), k.String("pcloudurl"), k.String("user"), k.String("pass"))
	client := pam.NewClient(k.String("pcloudurl"), config)
	client.RefreshSession()

	platforms, rc, err := client.GetPlatforms()
	if err != nil {
		panic(fmt.Sprintf("failed to get platforms: (%d) %s", rc, err.Error()))
	}
	ShowPlatforms(platforms)
	platform := AskUserChoosePlatform(platforms)
	ShowPlatform(platform)
}

func AskUserChoosePlatform(platforms pam.GetPlatformsResponse) pam.Platform {
	ShowPlatforms(platforms)

	scanner := bufio.NewScanner(os.Stdin)
	var text string
	fmt.Print("Choose platform, enter record number: ")
	scanner.Scan()
	text = scanner.Text()
	entry, err := strconv.Atoi(text)
	if err != nil {
		panic(fmt.Sprintf("failed to convert user input to integer: %s", err.Error()))
	}
	return platforms.Platforms[entry]
}

func ShowPlatforms(platforms pam.GetPlatformsResponse) {
	for entry := range platforms.Platforms {
		ShowPlatformOption(entry, platforms)
	}
}

func ShowPlatformOption(entry int, platforms pam.GetPlatformsResponse) {
	p := platforms.Platforms[entry]
	fmt.Printf("%d) ", entry)
	ShowPlatform(p)
}
func ShowPlatform(p pam.Platform) {
	fmt.Printf("PLATFORM ID: \"%s\", PLATFORM NAME: \"%s\"\n", p.General.ID, p.General.Name)
	fmt.Println("\tRequired")
	for _, prop := range p.Properties.Required {
		fmt.Printf("\t\t%s\n", prop.Name)
	}
	fmt.Println("\tOptional")
	for _, prop := range p.Properties.Optional {
		fmt.Printf("\t\t%s\n", prop.Name)
	}
}
