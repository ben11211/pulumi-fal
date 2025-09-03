package fal

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type DeployStrategy string
type AuthMode string

const (
	DeployStrategyRecreate DeployStrategy = "recreate"
	DeployStrategyRolling  DeployStrategy = "rolling"

	AuthModePublic  AuthMode = "public"
	AuthModePrivate AuthMode = "private"
	AuthModeShared  AuthMode = "shared"
)

type Client struct {
	key string
	dir string
}

type AuthOpts struct {
	Username            string
	Password            string
	PrivateKey          string
	InsecureHTTPAllowed bool
}

type DeployOpts struct {
	Name       string
	Entrypoint string
	Strategy   DeployStrategy
	AuthMode   AuthMode
}

type DeployResult struct {
	RevisionId string
	CreatedAt  string
	UpdatedAt  string
}

type App struct {
	Alias     string `json:"alias"`
	Revision  string `json:"revision"`
	AuthMode  string `json:"auth_mode"`
	UpdatedAt string `json:"updated_at"`
}

func NewClient(key string) (*Client, error) {
	if key == "" {
		return nil, fmt.Errorf("FAL_KEY is required")
	}

	tmpDir, err := os.MkdirTemp("", "pulumi-fal-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	return &Client{
		key: key,
		dir: tmpDir,
	}, nil
}

func (c *Client) Cleanup() {
	if c.dir != "" {
		os.RemoveAll(c.dir)
	}
}

func (c *Client) Deploy(ctx context.Context, gitURL string, authOpts *AuthOpts, opts *DeployOpts) (*DeployResult, error) {
	if gitURL == "" {
		return nil, fmt.Errorf("git URL is required")
	}

	repoDir := filepath.Join(c.dir, "repo")

	if err := c.cloneRepo(ctx, gitURL, repoDir, authOpts); err != nil {
		return nil, fmt.Errorf("failed to clone repository: %w", err)
	}

	args := []string{"deploy"}
	if opts.Strategy != "" {
		args = append(args, "--strategy", string(opts.Strategy))
	}
	if opts.AuthMode != "" {
		args = append(args, "--auth", string(opts.AuthMode))
	}
	if opts.Entrypoint != "" {
		args = append(args, opts.Entrypoint)
	}

	output, err := c.runFalCommand(ctx, repoDir, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy: %w", err)
	}

	now := time.Now().Format(time.RFC3339)
	return &DeployResult{
		RevisionId: c.extractRevisionFromOutput(output),
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func (c *Client) GetApp(ctx context.Context, name string) (*App, error) {
	output, err := c.runFalCommand(ctx, c.dir, "apps", "list", "--json")
	if err != nil {
		return nil, fmt.Errorf("failed to list apps: %w", err)
	}

	var apps []App
	if err := json.Unmarshal([]byte(output), &apps); err != nil {
		return nil, fmt.Errorf("failed to parse apps list: %w", err)
	}

	for _, app := range apps {
		if app.Alias == name {
			return &app, nil
		}
	}

	return nil, nil
}

func (c *Client) Delete(ctx context.Context, name string) error {
	_, err := c.runFalCommand(ctx, c.dir, "apps", "delete", name)
	if err != nil {
		return fmt.Errorf("failed to delete app: %w", err)
	}
	return nil
}

func (c *Client) cloneRepo(ctx context.Context, gitURL, repoDir string, authOpts *AuthOpts) error {
	args := []string{"clone"}

	if authOpts != nil {
		if authOpts.Username != "" && authOpts.Password != "" {
			cloneURL := strings.Replace(gitURL, "://", "://"+authOpts.Username+":"+authOpts.Password+"@", 1)
			args = append(args, cloneURL, repoDir)
		} else {
			args = append(args, gitURL, repoDir)
		}
	} else {
		args = append(args, gitURL, repoDir)
	}

	cmd := exec.CommandContext(ctx, "git", args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}

	return nil
}

func (c *Client) runFalCommand(ctx context.Context, workDir string, args ...string) (string, error) {
	keyCmd := exec.CommandContext(ctx, "fal", "profile", "key", "set", c.key)
	keyCmd.Dir = workDir
	_, keyErr := keyCmd.CombinedOutput()

	if keyErr != nil {
		return "", fmt.Errorf("failed to set fal profile key: %w", keyErr)
	}

	cmd := exec.CommandContext(ctx, "fal", args...)
	cmd.Dir = workDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("fal command failed: %w\nOutput: %s", err, string(output))
	}

	return string(output), nil
}

func (c *Client) extractRevisionFromOutput(output string) string {
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "revision") || strings.Contains(line, "Revision") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				return parts[len(parts)-1]
			}
		}
	}
	return "unknown"
}
