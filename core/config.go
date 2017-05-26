package core

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type jsonConfig struct {
	Rule        *JsonRule        `json:"rule"`
	GroupKey    string           `json:"group"`
	Calculation *JsonCalculation `json:"calculation"`
}

type JsonRule struct {
	Type  string      `json:"type"`
	Rules []*JsonRule `json:"rules"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type JsonCalculation struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type Config struct {
	json *jsonConfig
}

func (c *Config) GroupBy() string {
	if c.json.GroupKey == "" {
		return "app_id_code" // default
	}
	return c.json.GroupKey
}

func (c *Config) Rules() *JsonRule {
	return c.json.Rule
}

func (c *Config) Calculation() *JsonCalculation {
	return c.json.Calculation
}

func ParseConfig(reader io.Reader) (*Config, error) {
	c := &jsonConfig{}

	d, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(d, &c)

	return &Config{json: c}, err
}
