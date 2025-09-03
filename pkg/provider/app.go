package provider

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ben11211/pulumi-provider-fal/pkg/fal"
	"github.com/pulumi/pulumi-go-provider/infer"
)

type App struct{}

type GitConfig struct {
	URL                 string  `pulumi:"url"`
	Username            *string `pulumi:"username,optional"`
	Password            *string `pulumi:"password,optional"`
	PrivateKey          *string `pulumi:"privateKey,optional"`
	InsecureHTTPAllowed *bool   `pulumi:"insecureHttpAllowed,optional"`
}

type AppArgs struct {
	Name       string     `pulumi:"name"`
	Entrypoint string     `pulumi:"entrypoint"`
	Strategy   *string    `pulumi:"strategy,optional"`
	AuthMode   *string    `pulumi:"authMode,optional"`
	Git        *GitConfig `pulumi:"git,optional"`
}

type AppState struct {
	AppArgs

	RevisionId string `pulumi:"revisionId"`
	CreatedAt  string `pulumi:"createdAt"`
	UpdatedAt  string `pulumi:"updatedAt"`
}

func (a *App) Annotate(an infer.Annotator) {
	an.Describe(&a, "A fal application deployed from a git repository")
}

func (a *AppArgs) Annotate(an infer.Annotator) {
	an.Describe(&a.Name, "The name of the application")
	an.Describe(&a.Entrypoint, "The entrypoint for the application")
	an.Describe(&a.Strategy, "Deployment strategy (recreate or rolling)")
	an.Describe(&a.AuthMode, "Authentication mode (public, private, or shared)")
	an.Describe(&a.Git, "Git repository configuration")
}

func (a *AppState) Annotate(an infer.Annotator) {
	an.Describe(&a.RevisionId, "The revision ID of the deployed application")
	an.Describe(&a.CreatedAt, "When the application was created")
	an.Describe(&a.UpdatedAt, "When the application was last updated")
}

func (a *App) Create(ctx context.Context, req infer.CreateRequest[AppArgs]) (infer.CreateResponse[AppState], error) {
	if req.DryRun {
		return infer.CreateResponse[AppState]{
			ID: req.Name,
			Output: AppState{
				AppArgs:    req.Inputs,
				RevisionId: "preview",
				CreatedAt:  time.Now().Format(time.RFC3339),
				UpdatedAt:  time.Now().Format(time.RFC3339),
			},
		}, nil
	}

	config := infer.GetConfig[Config](ctx)
	if config.FalKey == "" {
		return infer.CreateResponse[AppState]{}, fmt.Errorf("FAL_KEY must be set")
	}

	client, err := fal.NewClient(config.FalKey)
	if err != nil {
		return infer.CreateResponse[AppState]{}, fmt.Errorf("failed to create fal client: %w", err)
	}
	defer client.Cleanup()

	deployOpts := &fal.DeployOpts{
		Name:       req.Inputs.Name,
		Entrypoint: req.Inputs.Entrypoint,
	}

	if req.Inputs.Strategy != nil {
		deployOpts.Strategy = fal.DeployStrategy(*req.Inputs.Strategy)
	} else {
		deployOpts.Strategy = fal.DeployStrategyRecreate
	}

	if req.Inputs.AuthMode != nil {
		deployOpts.AuthMode = fal.AuthMode(*req.Inputs.AuthMode)
	} else {
		deployOpts.AuthMode = fal.AuthModePrivate
	}

	var gitURL string
	var authOpts *fal.AuthOpts
	if req.Inputs.Git != nil {
		gitURL = req.Inputs.Git.URL
		authOpts = &fal.AuthOpts{}
		if req.Inputs.Git.Username != nil && req.Inputs.Git.Password != nil {
			authOpts.Username = *req.Inputs.Git.Username
			authOpts.Password = *req.Inputs.Git.Password
		}
		if req.Inputs.Git.PrivateKey != nil {
			authOpts.PrivateKey = *req.Inputs.Git.PrivateKey
		}
		if req.Inputs.Git.InsecureHTTPAllowed != nil {
			authOpts.InsecureHTTPAllowed = *req.Inputs.Git.InsecureHTTPAllowed
		}
	}

	result, err := client.Deploy(ctx, gitURL, authOpts, deployOpts)
	if err != nil {
		return infer.CreateResponse[AppState]{}, fmt.Errorf("failed to deploy app: %w", err)
	}

	state := AppState{
		AppArgs:    req.Inputs,
		RevisionId: result.RevisionId,
		CreatedAt:  result.CreatedAt,
		UpdatedAt:  result.UpdatedAt,
	}

	return infer.CreateResponse[AppState]{
		ID:     req.Name,
		Output: state,
	}, nil
}

func (a *App) Read(ctx context.Context, id string, state AppState) (string, AppState, error) {
	config := infer.GetConfig[Config](ctx)
	if config.FalKey == "" {
		return id, state, fmt.Errorf("FAL_KEY must be set")
	}

	client, err := fal.NewClient(config.FalKey)
	if err != nil {
		return id, state, fmt.Errorf("failed to create fal client: %w", err)
	}
	defer client.Cleanup()

	app, err := client.GetApp(ctx, state.Name)
	if err != nil {
		return id, state, fmt.Errorf("failed to get app: %w", err)
	}

	if app == nil {
		return "", AppState{}, os.ErrNotExist
	}

	state.RevisionId = app.Revision
	state.UpdatedAt = app.UpdatedAt

	return id, state, nil
}

func (a *App) Update(ctx context.Context, id string, state AppState, inputs AppArgs) (AppState, error) {
	createReq := infer.CreateRequest[AppArgs]{
		Name:   id,
		Inputs: inputs,
		DryRun: false,
	}

	createResp, err := a.Create(ctx, createReq)
	if err != nil {
		return AppState{}, err
	}

	newState := createResp.Output
	newState.CreatedAt = state.CreatedAt

	return newState, nil
}

func (a *App) Delete(ctx context.Context, id string, state AppState) error {
	config := infer.GetConfig[Config](ctx)
	if config.FalKey == "" {
		return fmt.Errorf("FAL_KEY must be set")
	}

	client, err := fal.NewClient(config.FalKey)
	if err != nil {
		return fmt.Errorf("failed to create fal client: %w", err)
	}
	defer client.Cleanup()

	return client.Delete(ctx, state.Name)
}
