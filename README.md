# gh-review-stats - GitHub Review Stats

A command line tool for examining review statistics for GitHub repositories.

## Installing

1. Download a pre-built binary for your platform from [the
   releases](https://github.com/dhellmann/gh-review-stats/releases)
   page on GitHub.
2. Unpack the archive.
3. Copy the binary to a directory in your `$PATH` (for example,
   `~/bin`).

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
$ gh-review-stats reviewers --org sphinx-contrib --repo datatemplates
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

## Pull Request Statistics

The `pull-requests` sub-command produces a CSV report with details of
pull requests that can be imported into other data analysis tools for
processing.

```console
$ gh-review-stats pull-requests -o metal3-io -r metal3-docs
Using config file: /Users/dhellmann/.gh-review-stats.yml
..............................................................................................................................................................
ID,Title,State,Author,URL,Created,Closed,Days to Merge
179,Add andfasano as approver,merged,hardys,https://github.com/metal3-io/metal3-docs/pull/179,2021-04-26,2021-04-26,0
174,Update inspection API proposal status,merged,fmuyassarov,https://github.com/metal3-io/metal3-docs/pull/174,2021-04-01,2021-04-01,0
172,add feruzjon muyassarov as approver,merged,dhellmann,https://github.com/metal3-io/metal3-docs/pull/172,2021-03-19,2021-03-19,0
169,Add update strategy to Metal3DataTemplate,merged,kashifest,https://github.com/metal3-io/metal3-docs/pull/169,2021-03-16,2021-04-01,16
166,Update disabling automated cleaning proposal,merged,fmuyassarov,https://github.com/metal3-io/metal3-docs/pull/166,2021-03-04,2021-03-17,13
164,Add explicit reboot mode options,merged,rdoxenham,https://github.com/metal3-io/metal3-docs/pull/164,2021-02-10,2021-02-24,13
163,Presentations framework proposal with a sample presentation,merged,hroyrh,https://github.com/metal3-io/metal3-docs/pull/163,2021-02-08,2021-04-28,79
162,✨ Proposal: node reuse,merged,furkatgofurov7,https://github.com/metal3-io/metal3-docs/pull/162,2021-02-04,2021-03-10,34
161,design: support automatic secure boot,merged,dtantsur,https://github.com/metal3-io/metal3-docs/pull/161,2021-02-02,2021-02-12,10
155,Add proposal for supporting external introspection,merged,hardys,https://github.com/metal3-io/metal3-docs/pull/155,2021-01-06,2021-03-19,71
152,Proposal for new parameters: Disk and NIC in HWCC,merged,Ashughorla,https://github.com/metal3-io/metal3-docs/pull/152,2020-12-16,2021-03-09,83
149,Add design proposal for label sync mechanism between BMHs and K Nodes,merged,Arvinderpal,https://github.com/metal3-io/metal3-docs/pull/149,2020-11-09,2021-02-17,100
147,Support for new parameters in HWCC.,merged,Ashughorla,https://github.com/metal3-io/metal3-docs/pull/147,2020-10-30,2021-03-09,130
138,design: add sub-states,merged,dtantsur,https://github.com/metal3-io/metal3-docs/pull/138,2020-09-21,2021-02-10,141
```

