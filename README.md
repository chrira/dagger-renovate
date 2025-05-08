# renovate

Update Dagger modules.

## Configure in GitHub

To run renovate in your GitHub repository for your repository use the file `.github/workflows/renovate.yml`

It works unchanged in every GitHub repository when you defined the `DAGGER_CLOUD_TOKEN` secret.

You must configure your GITHUB_TOKEN to have write access to your repository, see [Configuring the default GITHUB_TOKEN permissions](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/enabling-features-for-your-repository/managing-github-actions-settings-for-a-repository#configuring-the-default-github_token-permissions)

Set `Workflow permissions` to `Read and write permissions` and check `Allow GitHub Actions to create and approve pull requests`.

## Local Test

* https://docs.renovatebot.com/self-hosted-configuration/
* https://docs.renovatebot.com/examples/self-hosting/


### Setup

```sh
dagger init --sdk go --source bot
```

docker run --rm -v "/path/to/your/config.js:/usr/src/app/config.js" enovate/renovate:39.173

#### Post upgrade task

* https://docs.renovatebot.com/configuration-options/#postupgradetasks
* https://docs.renovatebot.com/templates/

Each command must match at least one of the patterns defined in `allowedCommands` (a global-only configuration option) in order to be executed.
If the list of allowed tasks is empty then no tasks will be executed.

You can use variable templating in your commands as long as `allowCommandTemplating` is enabled.

```json
{
  "postUpgradeTasks": {
    "commands": ["tslint --fix"],
    "fileFilters": ["yarn.lock", "**/*.js"],
    "executionMode": "update"
  }
}
```

### Run

#### Local

```sh
dagger call run --config ./config/config.json5
```

https://docs.renovatebot.com/modules/platform/local/

```sh
ln -s ../dagger-module-helm/ demo-src

dagger call local --src ./demo-src/

dagger call local --src https://github.com/puzzle/dagger-module-helm.git

dagger call local \
    --config ./config/config.json5 \
    --token env:MY_GH_PAT \
    --src demo-src/
```

#### GitLab

```sh
dagger call git-lab \
    --config ./config/config.json5 \
    --githubToken env:MY_GH_PAT \
    --gitlabToken env:MY_GL_PAT
```

#### GitHub

```sh
dagger call git-hub \
    --config ./config/config.json5 \
    --github-read-token env:MY_GH_PAT \
    --github-write-token env:MY_GH_WRITE_PAT
```
