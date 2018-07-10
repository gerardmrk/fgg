package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

type directive struct {
	Desc string     `toml:"desc"`
	Exec executable `toml:"exec"`
}

type executable struct {
	Cmd  string
	Args []string
	Envs map[string]string
}

func (e *executable) UnmarshalTOML(data interface{}) (err error) {
	switch v := data.(type) {
	case string:
		e.Cmd = fmt.Sprintf("%s -c", os.Getenv("SHELL"))
		e.Args = []string{v}
		break

	case map[string]interface{}:
		if _, ok := v["cmd"]; !ok {
			err = errors.New("invalid exec type")
			break
		}
		e.Cmd = v["cmd"].(string)
		if args, ok := v["args"]; ok {
			e.Args = args.([]string)
		}
		if envs, ok := v["envs"]; ok {
			e.Envs = envs.(map[string]string)
		}
		break

	default:
		fmt.Printf("NOPE: %T\n", v)
		err = errors.New("invalid exec type")
	}
	return
}
