package go_mysql_tools

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type DBConfig struct {
	Driver     string
	Host       string
	Port       int
	User       string
	Passwd     string
	Database   string
	Parameters string
}

func ReadConfigFromJsonFile(filename string) (*DBConfig, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	decoder := json.NewDecoder(reader)
	var config DBConfig
	err = decoder.Decode(&config)
	return &config, err
}

func (c *DBConfig) GetConnString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s%s", 
		c.User, c.Passwd, c.Host, c.Port, c.Database, c.Parameters)
}
