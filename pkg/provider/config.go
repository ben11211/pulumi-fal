package provider

import (
	"context"
	"os"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type Config struct {
	FalKey string `pulumi:"falKey,optional"`
}

func (c *Config) Annotate(a infer.Annotator) {
	a.Describe(&c.FalKey, "The FAL API key for authentication. Can also be set via FAL_KEY environment variable.")
}

func (c *Config) Configure(ctx context.Context) error {
	if c.FalKey == "" {
		c.FalKey = os.Getenv("FAL_KEY")
	}
	return nil
}
