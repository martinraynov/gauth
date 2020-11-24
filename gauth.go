package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/pcarrier/gauth/gauth"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	cfgPath := os.Getenv("GAUTH_CONFIG")
	if cfgPath == "" {
		user, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}

		cfgPath = path.Join(user.HomeDir, ".gauth/gauth.csv")

		if _, err := os.Stat(cfgPath); err != nil {
			log.Println("gauth.csv doesn't exist : ", err)
			log.Println("Do you wans to init the gauth.csv config file ? [y/N]")
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				if scanner.Text() == "y" {
					errCreate := os.MkdirAll(path.Join(user.HomeDir, ".gauth"), 0700)
					if errCreate != nil {
						log.Fatalln("Error when creating folder : ", path.Join(user.HomeDir, ".gauth"), errCreate)
					}

					_, errCreate = os.OpenFile(path.Join(user.HomeDir, ".gauth/gauth.csv"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
					if errCreate != nil {
						log.Fatalln("Error when creating gauth.csv file : ", path.Join(user.HomeDir, ".gauth/gauth.csv"), errCreate)
					}

					log.Println("Config file created at : ", path.Join(user.HomeDir, ".gauth/gauth.csv"))
					log.Println("You can now modify your config file and add your auth information (1 per line)")

					break
				} else {
					os.Exit(1)
				}
			}

			fmt.Scanln()
			os.Exit(1)
		}
	}

	cfgContent, err := gauth.LoadConfigFile(cfgPath, getPassword)
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	cfgReader := csv.NewReader(bytes.NewReader(cfgContent))
	// Unix-style tabular
	cfgReader.Comma = ':'

	cfg, err := cfgReader.ReadAll()
	if err != nil {
		log.Fatalf("Decoding CSV: %v", err)
	}

	go func() {
		for {
			// Clear terminal For windows
			c := exec.Command("cmd", "/c", "cls")
			c.Stdout = os.Stdout
			c.Run()

			// Clear terminal For linux
			fmt.Print("\033[H\033[2J")

			currentTS, progress := gauth.IndexNow()

			tw := tabwriter.NewWriter(os.Stdout, 0, 8, 1, ' ', 0)
			fmt.Fprintln(tw, "\tprev\tcurr\tnext")
			for _, record := range cfg {
				name, secret := record[0], record[1]
				prev, curr, next, err := gauth.Codes(secret, currentTS)
				if err != nil {
					log.Fatalf("Generating codes: %v", err)
				}
				fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", name, prev, curr, next)
			}
			tw.Flush()
			fmt.Printf("[%-29s]\n\n\n", strings.Repeat("=", progress))

			log.Println("Press the Enter Key to stop anytime")
			time.Sleep(200 * time.Millisecond)
		}
	}()

	fmt.Scanln()
	os.Exit(1)
}

func getPassword() ([]byte, error) {
	fmt.Printf("Encryption password: ")
	defer fmt.Println()
	return terminal.ReadPassword(int(syscall.Stdin))
}
