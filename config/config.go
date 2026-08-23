// пакет для работы с конфигурацией программы в которой указываются настройки программы,
// а так же чтение и запись настроек

package config

import (
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v3"
)


type Config struct {
	LogSettings LogSettings `yaml:"log_settings"`
}

type LogSettings struct {
	Directory 	string `yaml:"directory"`
	Level 		string `yaml:"level"`
}

func Load(configName string) *Config {

	data, err := ioutil.ReadFile(configName)
	if err != nil {

		fmt.Println("Ошибка загрузки настроек")

		cnf := createDefault()
		return &cnf
	}

	var cnf Config
	err = yaml.Unmarshal(data, &cnf)
	if err != nil {
		fmt.Println("Ошибка анмаршала yaml файла")
		cnf = createDefault()
		return &cnf
	}

	return &cnf
}

func createDefault() Config {
	return Config {
		LogSettings: LogSettings {
			Directory: "logs",
			Level: "debug",
		},
	}
}


