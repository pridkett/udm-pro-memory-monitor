/**
 * a simple program to log into a UDM Pro and reset
 * the system if the memory is too low.
 */

package main

import (
	"flag"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/melbahja/goph"
)

var memoryCmd = "cat /proc/meminfo"
var restartCmd = "unifi-os restart"
var memoryField = "MemAvailable"

func main() {
	udmProHost := flag.String("host", "192.168.1.1", "Hostname or IP address of UDM Pro")
	udmProUser := flag.String("user", "root", "Username to connect to UDM Pro")
	sshKey := flag.String("keyfile", "", "SSH key file for connection")
	sshKeyPass := flag.String("keypass", "", "Password for SSH key if needed")
	sshUseAgent := flag.Bool("agent", false, "Use SSH agent for authentication")
	udmProPassword := flag.String("password", "", "Password for UDM Pro user account")
	minMemoryAvailable := flag.Int64("memavailable", 200000, "minimum memory available in KB")
	runAsDaemon := flag.Bool("daemon", false, "run the application in a loop")
	daemonDelay := flag.Int("delay", 600, "delay between successive daemon calls")
	flag.Parse()

	if *udmProPassword == "" && *sshKey == "" && *sshUseAgent == false {
		log.Fatal("one of -keyfile, -password, or -agent must be given")
	}
	if *udmProPassword != "" && *sshKey != "" {
		log.Fatal("only one of -keyfile or -password may be specified")
	}

	log.Printf("Connecting to %s as user %s", *udmProHost, *udmProUser)

	var auth goph.Auth
	var err error

	if *sshUseAgent {
		auth, err = goph.UseAgent()
		if err != nil {
			log.Fatal(err)
		}
	} else if *sshKey != "" {
		auth, err = goph.Key(*sshKey, *sshKeyPass)

		if err != nil {
			log.Fatal(err)
		}
	} else {
		auth = goph.Password(*udmProPassword)
	}

	for true {
		client, err := goph.New(*udmProUser, *udmProHost, auth)

		if err != nil {
			log.Fatal(err)
		}

		out, err := client.Run(memoryCmd)

		if err != nil {
			log.Fatal(err)
		}

		md := parseMeminfo(string(out))
		log.Printf("current available memory: %d", md[memoryField])
		if md[memoryField] < *minMemoryAvailable {
			log.Printf("current available memory below threshold of %d", *minMemoryAvailable)
			out, err := client.Run(restartCmd)
			if err != nil {
				log.Fatalf("Error running command \"%s\": %s", restartCmd, err)
			}
			log.Printf(string(out))
		}
		client.Close()

		if !*runAsDaemon {
			break
		}

		time.Sleep(time.Duration(*daemonDelay) * time.Second)
	}
}

func parseMeminfo(meminfo string) map[string]int64 {
	rexp := regexp.MustCompile(`(?P<memtype>.+):\s+(?P<usage>\d+)`)
	lines := strings.Split(meminfo, "\n")

	dict := map[string]int64{}
	for i := 0; i < len(lines); i++ {
		matches := rexp.FindStringSubmatch(lines[i])
		if len(matches) > 1 {
			memUsage, _ := strconv.ParseInt(matches[rexp.SubexpIndex("usage")], 10, 0)
			dict[matches[rexp.SubexpIndex("memtype")]] = memUsage
		}
	}
	return dict
}
