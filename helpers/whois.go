package helpers

import (
	"log"
	"os/exec"
	"strings"
)

func AnalizeWhois(ipAddress string) (string, string) {
	whoisCountry := ""
	whoisOwner := ""
	out, err := exec.Command("sh", "-c", "whois "+ipAddress+" | grep -E \"OrgName|Country\"").Output()

	if err != nil {
		log.Println(err)
	} else {
		for _, line := range strings.Split(string(out), "\n") {
			if strings.HasPrefix(line, "Country:") {
				whoisCountry = strings.TrimPrefix(line, "Country:")
				whoisCountry = strings.TrimSpace(whoisCountry)
			}

			if strings.HasPrefix(line, "OrgName:") {
				whoisOwner = strings.TrimPrefix(line, "OrgName:")
				whoisOwner = strings.TrimSpace(whoisOwner)
			}
		}
	}

	return whoisCountry, whoisOwner
}
