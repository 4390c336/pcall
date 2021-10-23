package pcall

import (
	"log"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() { plugin.Register("pcall", setup) }

func setup(c *caddy.Controller) error {
	commandPath, err := parse(c)

	//parsing err
	if err != nil {
		return err
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Pcall{Next: next, CommandPath: commandPath}
	})

	return nil
}

func parse(c *caddy.Controller) (string, error) {
	var commandPath string

	for c.Next() {

		for c.NextBlock() {
			if c.Val() != "run" {
				return "", plugin.Error("pcall", c.Err("only `run` operation is supported"))
			}
			commandPath = c.RemainingArgs()[0]
		}
	}
	log.Println("commandPath:", commandPath)

	return commandPath, nil
}
