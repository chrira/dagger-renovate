// A generated module for Renovate functions
//
// This module updates Dagger modules.
// Dagger version
// Module dependencies
// Container Image tags in module files.

package main

import (
	"context"
	"dagger/renovate/internal/dagger"
	"time"
)

type Renovate struct{}

func (m *Renovate) Run(
	// config file
	config *dagger.File,
) *dagger.Container {
	return dag.Container().
		From("docker.io/renovate/renovate:39.173").
		WithFile("/usr/src/app/config.json", config).
		//WithEnvVariable("LOG_LEVEL", "debug").
		WithExec([]string{"renovate", "user/repo"})
}

func (m *Renovate) Local(
	ctx context.Context,
	// config file
	config *dagger.File,
	// repo path
	src *dagger.Directory,
	// GitHub token
	token *dagger.Secret,
) (string, error) {
	return m.renovateContainer(ctx, config, token).
		//WithMountedDirectory("/usr/src/app/", src, dagger.ContainerWithMountedDirectoryOpts{Owner: "1001"}).
		//WithMountedDirectory("/usr/src/app/", src, dagger.ContainerWithMountedDirectoryOpts{Owner: "ubuntu"}).
		WithMountedDirectory("/usr/src/app/", src).
		Terminal().
		WithExec([]string{
			"renovate",
				"--allowed-commands=\"[\"^dagger\\s*.*?develop$\"]\"", // for post upgrade task
				"--allow-command-templating=true", // for post upgrade task
				"--platform=local",
			}).
		Stdout(ctx)
}

func (m *Renovate) GitLab(
	ctx context.Context,
	// config file
	config *dagger.File,
	// docker socket -v /var/run/docker.sock:/var/run/docker.sock \
	//dockerSocket *dagger.Directory,
	// GitHub token any-personal-user-token-for-github-com-for-fetching-changelogs
	githubToken *dagger.Secret,
	// GitLab token your-github-enterprise-renovate-user-token
	gitlabToken *dagger.Secret,
) (string, error) {
	return m.renovateContainer(ctx, config, githubToken).
		// # You can set RENOVATE_AUTODISCOVER to true to run Renovate on all repos you have push access to
		// RENOVATE_AUTODISCOVER: 'false'
		WithEnvVariable("RENOVATE_AUTODISCOVER", "false").
		WithEnvVariable("RENOVATE_PLATFORM", "gitlab").
		WithEnvVariable("RENOVATE_ENDPOINT", "https://gitlab.com/api/v4/").
		WithEnvVariable("RENOVATE_GIT_AUTHOR", "Renovate Bot <noreply@gitlab.com>").
		//WithEnvVariable("RENOVATE_REPOSITORIES", "['puzzle/dagger-module-helm']"). // TODO not hardcoded
		WithSecretVariable("RENOVATE_TOKEN", gitlabToken).
		//WithMountedDirectory("/var/run/docker.sock", dockerSocket, dagger.ContainerWithMountedDirectoryOpts{Owner: "1001"}).
		/*
		WithExec(
			[]string{"renovate"},
			dagger.ContainerWithExecOpts{
				ExperimentalPrivilegedNesting: true,
				InsecureRootCapabilities: true,
			},
		).*/
		WithExec(
			[]string{"docker", "version"},
			dagger.ContainerWithExecOpts{
				ExperimentalPrivilegedNesting: true,
				InsecureRootCapabilities: true,
			},
		).
		Stdout(ctx)
}


func (m *Renovate) GitHub(
	ctx context.Context,
	// config file
	config *dagger.File,
	// GitHub token any-personal-user-token-for-github-com-for-fetching-changelogs
	githubReadToken *dagger.Secret,
	// GitHub token to do changes
	githubWriteToken *dagger.Secret,
	// renovate repository configuration
	// +optional
	renovateRepositories string,
) (string, error) {
	container := m.renovateContainer(ctx, config, githubReadToken).
		WithEnvVariable("RENOVATE_AUTODISCOVER", "false").
		WithEnvVariable("RENOVATE_PLATFORM", "github").
		WithEnvVariable("RENOVATE_GIT_AUTHOR", "Renovate Bot <noreply@github.com>").
		WithSecretVariable("RENOVATE_TOKEN", githubWriteToken)

	if renovateRepositories != "" {
		container = container.WithEnvVariable("RENOVATE_REPOSITORIES", renovateRepositories)
	}

	return container.
		/*
		WithExec(
			[]string{"docker", "version"},
			dagger.ContainerWithExecOpts{
				ExperimentalPrivilegedNesting: true,
				InsecureRootCapabilities: true,
			},
		).
		*/
		// Terminal().
		WithExec(
			[]string{"renovate"},
			dagger.ContainerWithExecOpts{
				ExperimentalPrivilegedNesting: true,
				InsecureRootCapabilities: true,
			},
		).
		Stdout(ctx)
}

func (m *Renovate) Hack(
	ctx context.Context,
	// config file
	config *dagger.File,
	// GitHub token any-personal-user-token-for-github-com-for-fetching-changelogs
	githubToken *dagger.Secret,
) (string, error) {
	return m.renovateContainer(ctx, config, githubToken).
		WithExec(
			[]string{"docker", "version"},
			dagger.ContainerWithExecOpts{
				ExperimentalPrivilegedNesting: true,
				InsecureRootCapabilities: true,
			},
		).
		Stdout(ctx)
}


func (m *Renovate) renovateContainer(
	ctx context.Context,
	// config file
	config *dagger.File,
	// GitHub token
	// +optional
	token *dagger.Secret,
) *dagger.Container {
	// Start an ephemeral dockerd
	dockerd := dag.Docker().Engine(dagger.DockerEngineOpts{
		Persist: false,
	})

	container := dag.Container().
		// official renovate image
		From("docker.io/renovate/renovate:39.173").
		//From("registry.access.redhat.com/ubi9/nodejs-20:9.5-1739783265").

		WithServiceBinding("dockerd", dockerd).
		WithEnvVariable("DOCKER_HOST", "tcp://dockerd:2375").

		//WithEnvVariable("GOPROXY", "direct").
		//WithEnvVariable("GOPROXY", "https://goproxy.io").

		//WithEnvVariable("_EXPERIMENTAL_DAGGER_RUNNER_HOST", "tcp://dockerd:2375").

		/*
		From("registry.access.redhat.com/ubi9/openjdk-17:1.18-1").
		WithUser("root").
		//WithExec([]string{"dnf", "module", "enable", "nodejs:20"}).
		WithExec([]string{"microdnf","install", "-y", "gzip"}).
		WithExec([]string{"microdnf","install", "-y", "nodejs"}).

		WithExec([]string{"dnf", "install", "-y", "dnf-plugins-core"}).
		WithExec([]string{"dnf-3", "config-manager", "--add-repo", "https://download.docker.com/linux/fedora/docker-ce.repo"}).
		WithExec([]string{"microdnf", "install", "-y", "docker-ce"}).
		*/

		//WithUser("root").
		//WithExec([]string{"npm", "install", "-g", "renovate", "yarn", "pnpm"}). // TODO take fix renovate version and let it update by mr
		
		/*
		// install docker
		WithExec([]string{"yum", "install", "-y", "yum-utils", "device-mapper-persistent-data", "lvm2"}).
		WithExec([]string{"yum-config-manager", "--add-repo", "https://download.docker.com/linux/centos/docker-ce.repo"}).
		WithExec([]string{"yum", "install", "-y", "docker-ce"}).
		// docker ohne sudo
		WithExec([]string{"usermod", "-aG", "docker", "default"}).
		*/

		WithUser("root").
		WithExec([]string{"mkdir", "-p", "/src/helm"}).
		WithExec([]string{"touch", "/schema.json"}).
		WithUser("ubuntu").

		//WithUser("1001").
		WithWorkdir("/usr/src/app/").
		WithFile("/tmp/config.json5", config).
		WithEnvVariable("RENOVATE_CONFIG_FILE", "/tmp/config.json5").
		/*WithUser("root").
		WithExec([]string{"chown", "-R", "1001", "/usr/src/app/"}). // make workdir owned by container user
		WithExec([]string{"chgrp", "-R", "0", "/usr/src/app/"}). // make workdir group modifiable
		WithExec([]string{"chmod", "-R", "g=u", "/usr/src/app/"}). // make workdir group modifiable
		WithExec([]string{"chown", "-R", "1001", "/tmp/"}). // make tmp dir owned by container user
		WithExec([]string{"chgrp", "-R", "0", "/tmp/"}). // make tmp dir group modifiable
		WithExec([]string{"chmod", "-R", "g=u", "/tmp/"}). // make tmp dir group modifiable
		WithExec([]string{"mkdir", "/src/"}). // prepare /src dir
		WithExec([]string{"chown", "-R", "1001", "/src/"}). // make src dir owned by container user
		WithExec([]string{"chgrp", "-R", "0", "/src/"}). // make src dir group modifiable
		WithExec([]string{"chmod", "-R", "g=u", "/src/"}). // make src dir group modifiable
		WithUser("1001").
		*/


		WithUser("root").
		WithExec([]string{"chown", "-R", "ubuntu", "/usr/src/app/"}). // make workdir owned by container user
		WithExec([]string{"chgrp", "-R", "0", "/usr/src/app/"}). // make workdir group modifiable
		WithExec([]string{"chmod", "-R", "g=u", "/usr/src/app/"}). // make workdir group modifiable
		WithExec([]string{"chown", "-R", "ubuntu", "/tmp/"}). // make tmp dir owned by container user
		WithExec([]string{"chgrp", "-R", "0", "/tmp/"}). // make tmp dir group modifiable
		WithExec([]string{"chmod", "-R", "g=u", "/tmp/"}). // make tmp dir group modifiable
		WithExec([]string{"mkdir", "-p", "/src/"}). // prepare /src dir
		WithExec([]string{"chown", "-R", "ubuntu", "/src/"}). // make src dir owned by container user
		WithExec([]string{"chgrp", "-R", "0", "/src/"}). // make src dir group modifiable
		WithExec([]string{"chmod", "-R", "g=u", "/src/"}). // make src dir group modifiable
		WithUser("ubuntu").


		WithEnvVariable("LOG_LEVEL", "debug")

	if (token != nil) {
		 // optional, not used when running on github.com
		container = container.WithSecretVariable("GITHUB_COM_TOKEN", token)
	}

	// invalidate the cache to never cache the renovate execution
	container = container.WithEnvVariable("CACHEBUSTER", time.Now().String())

	return container
}
