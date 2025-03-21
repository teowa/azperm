package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:      "azperm",
		Usage:     "Retrieves Azure resource permissions from Terraform configurations (Azure/azapi provider only) or specified resource types. Requires Azure CLI login.",
		UsageText: "azperm [--file-name <file_name>] [--resource-type <resource_type>]",
		Action: func(cCtx *cli.Context) error {
			if len(cCtx.StringSlice("file-name")) > 0 || len(cCtx.StringSlice("resource-type")) > 0 {
				action(cCtx.StringSlice("file-name"), cCtx.StringSlice("resource-type"))
			} else {
				cli.ShowSubcommandHelp(cCtx)
			}
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:  "file-name",
				Usage: "One or more Terraform configuration files to parse, separated by commas(`,`).",
			},

			&cli.StringSliceFlag{
				Name:  "resource-type",
				Usage: "One or more resource types to parse, separated by commas(`,`).",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var azProviderOperationCache map[string]map[string]interface{}

func action(fileNames, resourceTypes []string) {
	azProviderOperationCache := make(map[string]map[string]interface{}, 0)

	result := make(map[string]interface{}, 0)
	parsedResourceTypes := make([]string, 0)

	for _, fn := range fileNames {
		f, err := os.Open(fn)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		stat, err := f.Stat()
		if err != nil {
			log.Fatal(err)
		}

		buf := make([]byte, stat.Size())
		_, err = f.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		hclFile, diags := hclwrite.ParseConfig(buf, "main.tf", hcl.InitialPos)
		if diags != nil {
			log.Fatal(diags)
		}

		body := hclFile.Body()
		for _, block := range body.Blocks() {
			if (block.Type() == "resource" || block.Type() == "data") && strings.Contains(block.Labels()[0], "azapi") {
				if v, ok := block.Body().Attributes()["type"]; ok {
					parsedResourceTypes = append(parsedResourceTypes, strings.Trim(string(v.Expr().BuildTokens(nil).Bytes()), ` "`))
				}
			}
		}
	}

	if len(resourceTypes) > 0 {
		parsedResourceTypes = append(parsedResourceTypes, resourceTypes...)
	}

	for _, typeRaw := range parsedResourceTypes {
		actions := make([]string, 0)
		dataActions := make([]string, 0)
		rp, _, _ := strings.Cut(typeRaw, "/")
		if rp == "" {
			continue
		}

		resourceType, _, _ := strings.Cut(typeRaw, "@")
		if strings.EqualFold(resourceType, "Microsoft.Resources/resourceGroups") {
			resourceType = "Microsoft.Resources/subscriptions/resourceGroups"
		}

		var values map[string]interface{}
		var ok bool
		if values, ok = azProviderOperationCache[rp]; !ok {
			cmd := exec.Command("az", "provider", "operation", "show", "-n", rp)

			output, err := cmd.Output()
			if err != nil {
				log.Fatalf("runing `%s`, err:%+v", cmd.String(), err)
			}

			if err := json.Unmarshal(output, &values); err != nil {
				log.Fatal(err)
			}

			azProviderOperationCache[rp] = values
		}

		for _, opRaw := range values["operations"].([]interface{}) {
			op := opRaw.(map[string]interface{})
			opName := op["name"].(string)
			tempOpName, _ := strings.CutSuffix(opName, "/action")
			if strings.EqualFold(tempOpName, resourceType) {
				if op["isDataAction"].(bool) {
					dataActions = append(dataActions, opName)
				} else {
					actions = append(actions, opName)
				}
			}

			tempOpName = tempOpName[:strings.LastIndex(tempOpName, "/")]
			if strings.EqualFold(tempOpName, resourceType) {
				if op["isDataAction"].(bool) {
					dataActions = append(dataActions, opName)
				} else {
					actions = append(actions, opName)
				}
			}
		}

		for _, v := range values["resourceTypes"].([]interface{}) {
			for _, opRaw := range v.(map[string]interface{})["operations"].([]interface{}) {
				op := opRaw.(map[string]interface{})
				opName := op["name"].(string)
				tempOpName, _ := strings.CutSuffix(opName, "/action")
				if strings.EqualFold(tempOpName, resourceType) {
					if op["isDataAction"].(bool) {
						dataActions = append(dataActions, opName)
					} else {
						actions = append(actions, opName)
					}
				}

				tempOpName = tempOpName[:strings.LastIndex(tempOpName, "/")]
				if strings.EqualFold(tempOpName, resourceType) {
					if op["isDataAction"].(bool) {
						dataActions = append(dataActions, opName)
					} else {
						actions = append(actions, opName)
					}
				}
			}
		}

		result[typeRaw] = map[string]interface{}{
			"action":      actions,
			"dataActions": dataActions,
		}
	}

	output, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(output))
}
