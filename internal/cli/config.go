package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/EchoSell/keyapi-cli/internal/config"
)

func runConfig(args []string) error {
	if len(args) == 0 || isHelp(args[0]) {
		printConfigHelp()
		return nil
	}

	switch args[0] {
	case "show", "get":
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if args[0] == "get" && len(args) == 2 {
			switch args[1] {
			case "token", "api-key":
				fmt.Println(mask(cfg.APIKey))
				return nil
			case "auth-scheme":
				fmt.Println(cfg.AuthScheme)
				return nil
			case "language", "lang":
				fmt.Println(cfg.Language)
				return nil
			default:
				return fmt.Errorf("unknown config key %q", args[1])
			}
		}
		cfg.APIKey = mask(cfg.APIKey)
		return printJSON(configView{
			APIKey:     cfg.APIKey,
			AuthScheme: cfg.AuthScheme,
			Language:   cfg.Language,
		})
	case "set":
		if len(args) >= 2 && isHelp(args[1]) {
			printConfigHelp()
			return nil
		}
		if len(args) != 3 {
			return errors.New("usage: keyapi config set <token|api-key|auth-scheme|language> <value>")
		}
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		switch args[1] {
		case "token", "api-key":
			cfg.APIKey = args[2]
		case "auth-scheme":
			cfg.AuthScheme = args[2]
		case "language", "lang":
			language := strings.ToLower(strings.TrimSpace(args[2]))
			if language != "en" && language != "zh" {
				return errors.New("language must be en or zh")
			}
			cfg.Language = language
		default:
			return fmt.Errorf("unknown config key %q", args[1])
		}
		if err := config.Save(cfg); err != nil {
			return err
		}
		fmt.Println("saved")
		return nil
	case "clear":
		if len(args) != 2 {
			return errors.New("usage: keyapi config clear <token|api-key>")
		}
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		switch args[1] {
		case "token", "api-key":
			cfg.APIKey = ""
		default:
			return fmt.Errorf("unknown config key %q", args[1])
		}
		if err := config.Save(cfg); err != nil {
			return err
		}
		fmt.Println("cleared")
		return nil
	default:
		return fmt.Errorf("unknown config subcommand %q", args[0])
	}
}
