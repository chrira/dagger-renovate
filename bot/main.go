// A generated module for Renovate functions
//
// This module updates Dagger modules:
// Dagger Engine version
// Module dependencies (TODO)
// Container Image tags in module files (TODO).

package main

import (
	"context"
	"dagger/dagger-renovate/internal/dagger"
	"fmt"
)

type DaggerRenovate struct{}

//
// Renovate Dagger Modules on GitHub 
//
// Attention:
// Dagger Engine running this method must be of version latest.
// If the Engine started inside this container has a different version, the outer Engine will be stopped.
//
 func (m *DaggerRenovate) GitHub(
	ctx context.Context,
	// config file
	config *dagger.File,
	// GitHub token any-personal-user-token-for-github-com-for-fetching-changelogs
	githubReadToken *dagger.Secret,
	// GitHub token to do changes
	githubWriteToken *dagger.Secret,
	// Docker Socket
	dockerSock *dagger.Socket,
	// renovate repository configuration
	// +optional
	renovateRepositories string,
	// Docker Engine version
	// +optional
	// +default="24.0"
	version string,
) (string, error) {
	container := m.RenovateContainer(ctx, version).
		WithEnvVariable("RENOVATE_PLATFORM", "github").
		WithEnvVariable("RENOVATE_GIT_AUTHOR", "Renovate Bot <noreply@github.com>").
		WithEnvVariable("LOG_LEVEL", "debug").
		WithSecretVariable("RENOVATE_TOKEN", githubWriteToken).
		WithSecretVariable("GITHUB_COM_TOKEN", githubReadToken).
		WithUnixSocket("/var/run/docker.sock", dockerSock).
		WithFile("/tmp/config.json5", config).
		WithEnvVariable("RENOVATE_CONFIG_FILE", "/tmp/config.json5")

	if renovateRepositories != "" {
		container = container.
			WithEnvVariable("RENOVATE_REPOSITORIES", renovateRepositories).
			WithEnvVariable("RENOVATE_AUTODISCOVER", "false")
	}


	// TODO: cache buster needed?
	// invalidate the cache to never cache the renovate execution
	// WithEnvVariable("CACHEBUSTER", time.Now().String()).

	return container.
		WithExec(
			[]string{"renovate"},
			dagger.ContainerWithExecOpts{
				ExperimentalPrivilegedNesting: true,
				InsecureRootCapabilities: true,
			},
		).
		Stdout(ctx)
}

//
// Build a renovate dind container
//
// Run container locally:
// dagger -m bot call dind-container \
//   export --path=./dind-container.tar
//
// docker load -q -i ./dind-container.tar
// docker tag <container-id> dind-container:latest
//
// sudo docker run \
//   -v /var/run/docker.sock:/var/run/docker.sock \
//   -v $(pwd)/config/:/tmp/config/:ro \
//   --entrypoint /bin/sh \
//   -ti \
//   --env LOG_LEVEL="debug" \
//   --env GITHUB_COM_TOKEN=${MY_GH_PAT} \
//   --env RENOVATE_AUTODISCOVER="false" \
//   --env RENOVATE_CONFIG_FILE="/tmp/config/config.json5" \
//   --env RENOVATE_PLATFORM="github" \
//   --env RENOVATE_GIT_AUTHOR="Renovate Bot <noreply@github.com>" \
//   --env RENOVATE_TOKEN=${MY_GH_WRITE_PAT} \
//   --env RENOVATE_REPOSITORIES="['chrira/dagger-module-helm']" \
//   dind-container:latest
//
//   renovate > /tmp/renovate.log
//
//   TODO:
//   no repository config with:
//   "extends": [
//     "config:recommended"
//   ],
//
func (m *DaggerRenovate) RenovateContainer(
	ctx context.Context,
	// Docker Engine version
	// +optional
	// +default="24.0"
	version string,
) *dagger.Container {
	container := dag.Container().

		// dind container is based on alpine
		From(fmt.Sprintf("index.docker.io/docker:%s-dind", version)).
		WithWorkdir("/usr/src/app").
		WithExec(
			[]string{"sh", "-c", "apk --no-cache add curl"},
		).
		WithExec([]string{"apk","add", "nodejs", "npm"}).
		//pkg install nodejs
		WithExec([]string{"npm", "install", "-g", "renovate", "yarn", "pnpm"}) // TODO take fix renovate version and let it update by mr

	return container
}
