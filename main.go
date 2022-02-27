package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/xlzd/gotp"
)

func main() {
	formatJSON := false
	flag.BoolVar(&formatJSON, "json", false, "output in json")
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "need 1pif file path")
		os.Exit(1)
	}
	path := flag.Arg(0)
	if err := run(path, formatJSON); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(path string, formatJSON bool) error {
	otps, err := loadOTPs(path)
	if err != nil {
		return err
	}
	for _, otp := range otps {
		u, err := url.Parse(otp.URI)
		if err != nil {
			return err
		}
		secret := u.Query().Get("secret")
		t := gotp.NewDefaultTOTP(strings.ToUpper(secret))
		number, expiration := t.NowWithExpiration()
		otp.Number = number
		otp.Expiration = expiration
	}
	if formatJSON {
		b, err := json.MarshalIndent(otps, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	} else {
		for _, otp := range otps {
			fmt.Printf("* %s (%s)\n  %s\n", otp.Title, otp.Username, otp.Number)
		}
	}
	return nil

}

type OTP struct {
	Title      string
	Username   string
	URI        string
	Number     string
	Expiration int64
}

func loadOTPs(path string) ([]*OTP, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var otps []*OTP
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "{") {
			continue
		}
		var item struct {
			Title          string
			SecureContents struct {
				Fields []struct {
					Name  string
					Value string
				}
				Sections []struct {
					Fields []struct {
						V string
					}
				}
			}
		}
		if err := json.Unmarshal([]byte(line), &item); err != nil {
			continue
		}
		var uri string
		for _, section := range item.SecureContents.Sections {
			for _, field := range section.Fields {
				if strings.HasPrefix(field.V, "otpauth://totp/") {
					uri = field.V
				}
			}
		}
		if uri == "" {
			continue
		}

		var username string
		for _, field := range item.SecureContents.Fields {
			if field.Name == "username" {
				username = field.Value
			}
		}
		otps = append(otps, &OTP{
			Title:    item.Title,
			Username: username,
			URI:      uri,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	sort.Slice(otps, func(i, j int) bool {
		ti := otps[i].Title
		tj := otps[j].Title
		if ti != tj {
			return ti < tj
		}
		ui := otps[i].Username
		uj := otps[j].Username
		return ui < uj
	})
	return otps, nil
}
