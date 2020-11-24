package main

import (
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

		// If argument is passed use selected function
		if os.Args != nil && len(os.Args) > 1 {
			switch os.Args[1] {
			case "init":
				if os.MkdirAll(path.Join(user.HomeDir, ".gauth"), 0700); err != nil {
					log.Fatalln("Error when creating folder : ", path.Join(user.HomeDir, ".gauth"), err)
				}

				_, err := os.OpenFile(path.Join(user.HomeDir, ".gauth/gauth.csv"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
				if err != nil {
					log.Fatalln("Error when creating gauth.csv file : ", path.Join(user.HomeDir, ".gauth/gauth.csv"), err)
				}

				log.Println("Config file created at : ", path.Join(user.HomeDir, ".gauth/gauth.csv"))
				log.Fatalln("You can now modify your config file and add your auth information (1 per line)")
				break
			case "help":
				log.Println("Welcome to gauth (modified by martin.raynov)")
				log.Println("============================================")
				log.Println("1. Use 'gauth init' command to to init the application. This command will :")
				log.Println("   -> create the directory :", path.Join(user.HomeDir, ".gauth"))
				log.Println("   -> create the config file inside : gauth.csv")
				log.Println("2. Modify your config file ", path.Join(user.HomeDir, ".gauth/gauth.csv"), "and add your auth information (1 per line)")
				log.Println("   Example : Github:234567qrstuvwxyz")
				log.Println("3. Start your application 'gauth' and enjoy !")
				fmt.Scanln()
				os.Exit(1)

			default:
				log.Fatalln("Unknown argument :", os.Args[1])
			}
		}

		cfgPath = path.Join(user.HomeDir, ".gauth/gauth.csv")
	}

	if _, err := os.Stat(cfgPath); err != nil {
		log.Println("gauth.csv doesn't exist : ", err)
		log.Println("Use 'gauth init' command to init the file !")
		log.Println("Press the Enter Key to stop anytime")
		fmt.Scanln()
		os.Exit(1)
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
			c := exec.Command("cls")
			c.Stdout = os.Stdout
			c.Run()
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
			time.Sleep(100 * time.Millisecond)
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
