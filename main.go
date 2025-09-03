package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ben11211/pulumi-provider-fal/pkg/provider"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

var version = "0.0.1"

func main() {
	// Handle version command
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println(version)
		return
	}

	p.RunProvider(context.Background(), "fal", version,
		infer.Provider(infer.Options{
			Resources: []infer.InferredResource{
				infer.Resource(&provider.App{}),
			},
			ModuleMap: map[tokens.ModuleName]tokens.ModuleName{
				"provider": "index",
			},
			Config: infer.Config(&provider.Config{}),
		}))
}
