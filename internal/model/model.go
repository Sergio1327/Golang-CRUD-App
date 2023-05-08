package model

import "fmt"

type CV struct {
	Name   string `json:"name"`
	Job    string `json:"job"`
	Salary string `json:"salary"`
}

func (c *CV) Info() string {
	str := fmt.Sprintf("Name:%s | Job:%s | salary:%s", c.Name, c.Job, c.Salary)
	return str
}
