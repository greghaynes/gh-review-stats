# gh-review-stats - GitHub Review Stats

A command line tool for examining review statistics for GitHub repositories.

## Configuration

The default configuration file is `~/.gh-review-stats.yml`. The
`config-gen` command can be used to create a skeleton configuration
file, or to update the file if new configuration options are added.

```console
$ gh-review-stats config-gen
Using config file: /Users/doc/.gh-review-stats.yml
wrote "/Users/doc/.gh-review-stats.yml"
```

### github.token

The `github.token` is a [personal access
token](https://github.com/settings/tokens) used to access the GitHub
API to increase the hourly API call limit.
