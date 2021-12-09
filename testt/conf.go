package testt

import (
	"testing"

	"github.com/Weeping-Willow/api-example/config"
)

var Conf *config.Config

func LoadConfig(m *testing.M, path string) {
	conf, _ := config.New(path)
	Conf = conf
}
