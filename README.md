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

### reviewers.ignore

The `reviewers.ignore` option is a list of GitHub account names to not
include in the output. This is useful for ignoring bot accounts or
spammers.

```yaml
reviewers:
  ignore:
    - "dependabot[bot]"
```

## Reviewer Statistics

The `reviewers` sub-command generates a report showing the number of
comments made across pull requests in a repository by each reviewer.

```console
$ go run ./main.go  reviewers --org sphinx-contrib --repo datatemplates
Using config file: /Users/dhellmann/.gh-review-stats.yml
...............................................

2/2: janbrohl
	  1: https://github.com/sphinx-contrib/datatemplates/pull/79 [dhellmann] "docs: remove file section from inline example"
	  1: https://github.com/sphinx-contrib/datatemplates/pull/77 [dhellmann] "docs: update use instructions"
2/1: dhellmann
	  2: https://github.com/sphinx-contrib/datatemplates/pull/77 [dhellmann] "docs: update use instructions"
1/1: kevung
	  1: https://github.com/sphinx-contrib/datatemplates/pull/77 [dhellmann] "docs: update use instructions"
```

The report is formatted as

```text
<total comment count>/<total PR count>: <github name>
      <PR comment count>: <PR URL> [<PR author>] "<PR title>"
```

