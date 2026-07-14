package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/EchoSell/keyapi-cli/internal/openapi"
)

func runSpec(args []string) error {
	if len(args) == 0 || isHelp(args[0]) {
		printSpecHelp()
		return nil
	}

	switch args[0] {
	case "list-platforms":
		if len(args) == 2 && isHelp(args[1]) {
			printSpecPlatformsHelp()
			return nil
		}
		platforms, err := openapi.LoadPlatforms()
		if err != nil {
			return err
		}
		return printJSON(platforms)
	case "list-endpoints":
		if len(args) >= 2 && isHelp(args[1]) {
			printSpecListHelp()
			return nil
		}
		if len(args) < 2 {
			return errors.New("usage: keyapi spec list-endpoints <platform...|all>")
		}
		grouped, err := openapi.LoadEndpointsForPlatforms(args[1:])
		if err != nil {
			return err
		}
		view := map[string][]endpointSummary{}
		for platform, endpoints := range grouped {
			for _, endpoint := range endpoints {
				view[platform] = append(view[platform], endpointSummary{
					Method:      endpoint.Method,
					Path:        endpoint.Path,
					Summary:     endpoint.Summary,
					Description: endpoint.Description,
				})
			}
		}
		return printJSON(view)
	case "get-endpoint":
		if len(args) >= 2 && isHelp(args[1]) {
			printSpecShowHelp()
			return nil
		}
		return runEndpointSchema(args[1:], "keyapi spec get-endpoint <endpoint-or-path> [--format json|pretty]")
	default:
		return fmt.Errorf("unknown spec subcommand %q", args[0])
	}
}

func runSchema(args []string) error {
	if len(args) == 0 || isHelp(args[0]) {
		printSchemaHelp()
		return nil
	}
	return runEndpointSchema(args, "keyapi schema <endpoint-or-path> [--format json|pretty]")
}

func runEndpointSchema(args []string, usage string) error {
	endpointName, format, err := parseTargetAndFormatDefault(args, usage, "json")
	if err != nil {
		return err
	}
	registry, err := openapi.LoadRegistryForPlatform(platformFromEndpoint(endpointName))
	if err != nil {
		return err
	}
	endpoint, ok := registry.Find(endpointName)
	if !ok {
		return fmt.Errorf("endpoint %q not found; run keyapi spec list-endpoints <platform> to inspect available endpoints", endpointName)
	}
	if format == "json" {
		return printJSON(endpoint)
	}
	if format != "pretty" {
		return fmt.Errorf("unsupported format %q; use pretty or json", format)
	}
	return openapi.PrintEndpointDetail(os.Stdout, endpoint)
}
