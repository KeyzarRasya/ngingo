package core

type NgingoConfiguration struct {
	Name 	string		`yaml:"name"`
	Server	[]string	`yaml:"server"`
}