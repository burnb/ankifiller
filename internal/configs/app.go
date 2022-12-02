package configs

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	Name                  = "ankifiller"
	defaultConfigFilePath = "config.yml"
)

type Anki struct {
	Deck                    string  `yaml:"deck"`
	ImageDefinitionField    *string `yaml:"image_definition_field"`
	ImageField              *string `yaml:"image_field"`
	PhonemicDefinitionField *string `yaml:"phonemic_definition_field"`
	PhonemicField           *string `yaml:"phonemic_field"`
}

type Phonemic struct {
	Locale string `yaml:"locale"`
	System string `yaml:"system"`
}

func (c *Phonemic) GetLocale() string {
	return c.Locale
}

func (c *Phonemic) GetSystem() string {
	return c.System
}

type GoogleCustomSearch struct {
	APIKey string  `yaml:"api_key"`
	Cx     string  `yaml:"cx"`
	Gl     *string `yaml:"gl"`
}

func (c *GoogleCustomSearch) GetAPIKey() string {
	return c.APIKey
}

func (c *GoogleCustomSearch) GetCx() string {
	return c.Cx
}

func (c *GoogleCustomSearch) GetGl() *string {
	return c.Gl
}

type App struct {
	ConfigFilePath     string
	Anki               Anki                `yaml:"anki"`
	Phonemic           *Phonemic           `yaml:"phonemic"`
	GoogleCustomSearch *GoogleCustomSearch `yaml:"google_custom_search"`
	flagSet            *flag.FlagSet
}

func (c *App) Prepare() (err error) {
	c.flagSet = flag.NewFlagSet(Name, flag.ExitOnError)
	c.flagSet.Usage = c.printUsage

	help := c.flagSet.Bool("h", false, "Show this help message and exit")
	c.flagSet.StringVar(&c.ConfigFilePath, "c", defaultConfigFilePath, "Config file path")
	if err = c.flagSet.Parse(os.Args[1:]); err != nil {
		return err
	}

	if *help {
		c.flagSet.Usage()
		os.Exit(1)
	}

	file, err := os.ReadFile(c.ConfigFilePath)
	if err != nil {
		return fmt.Errorf("can't read config file: %w", err)
	}

	if err = yaml.Unmarshal(file, c); err != nil {
		return fmt.Errorf("can't parse config file: %w", err)
	}

	return nil
}

func (c *App) printUsage() {
	fmt.Println("\nUsage: ankifiller [OPTIONS]\nFill your language cards with auto-generated data\n\nOptions:")
	c.flagSet.PrintDefaults()
}
