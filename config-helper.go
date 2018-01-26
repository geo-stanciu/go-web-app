package main

import (
	"encoding/xml"
	"os"
)

// Configuration - config helper
type Configuration struct {
	XMLName       xml.Name `xml:"config"`
	General       ConfigurationGeneral
	Database      ConfigurationDatabase
	PasswordRules ConfigurationPassword
}

// ConfigurationGeneral - general config
type ConfigurationGeneral struct {
	XMLName  xml.Name `xml:"general"`
	Port     string   `xml:"port"`
	Timezone string   `xml:"timezone"`
	IsHTTPS  bool     `xml:"use-https"`
}

// ConfigurationDatabase - database config
type ConfigurationDatabase struct {
	XMLName xml.Name `xml:"database"`
	DbType  string   `xml:"db-type"`
	DbURL   string   `xml:"db-url"`
}

// ConfigurationPassword - password config
type ConfigurationPassword struct {
	XMLName                   xml.Name `xml:"password-rules"`
	ChangeInterval            int      `xml:"change-interval,attr"`
	PasswordFailInterval      int      `xml:"password-fail-interval,attr"`
	MaxAllowedFailedAtmpts    int      `xml:"max-allowed-failed-atmpts,attr"`
	NotRepeatLastXPasswords   int      `xml:"not-repeat-last-x-passwords,attr"`
	MinCharacters             int      `xml:"min-characters,attr"`
	MinLetters                int      `xml:"min-letters,attr"`
	MinCapitals               int      `xml:"min-capitals,attr"`
	MinDigits                 int      `xml:"min-digits,attr"`
	MinNonAlphaNumerics       int      `xml:"min-non-alpha-numerics,attr"`
	AllowRepetitiveCharacters bool     `xml:"allow-repetitive-characters,attr"`
	CanContainUsername        bool     `xml:"can-contain-username,attr"`
}

// ReadFromFile - read config from file
func (c *Configuration) ReadFromFile(cfgFile string) error {
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		return err
	}

	file, err := os.Open(cfgFile)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)

	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "config" {
				decoder.DecodeElement(c, &se)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
