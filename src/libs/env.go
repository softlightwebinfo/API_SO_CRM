package libs

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
	"os"
	"so-crm/src/models"
)

type Env struct {
	ConfigFile string
}

func (s *Env) processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func (s *Env) ReadFile(cfg *models.Config) {
	f, err := os.Open(s.ConfigFile)
	if err != nil {
		s.processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		s.processError(err)
	}
}

func (s *Env) ReadEnv(cfg *models.Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		s.processError(err)
	}
}
