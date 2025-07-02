package model

const PATH = "./config.yaml"

type Configuration struct {
	ConfigPath 	string `yaml:"configPath"`;
	Port		string `yaml:"port"`	
}