# trueblocks-giveth

This repo was used as part of this proposal to the Giveth forum: https://forum.giveth.io/t/proposal-monitoring-tool-to-help-identify-and-mitigate-recirculating-givbacks. The work entailed trying to identify, in an automated way, recirculated Givbacks. More information on this problem is here: https://docs.giveth.io/giveconomy/givbacks and in particular here: https://docs.giveth.io/giveconomy/givbacks/#disqualifying-factors-for-the-givbacks-program.

## Command Line Tool

We built a simple command line tool that currently allows us to download data from the Giveth APIs. While this is not all the data we need for the project, it's a start. We will documented the data structures [here](./data/QUESTIONS.md).

```
Data analysis for Giveth

Usage:
  giveth [command]

Available Commands:
  data        Various routines to download and manipulate the data
  projects    Produces data related to the projects including addresses.tsv
  rounds      Print information about the rounds
  summarize   Summarizes and combines data by type and time period (i.e. rounds)

Flags:
  -r, --round uint   Limits the list of rounds to a single round
  -u, --update       If present, data commands pull data from Giveth's APIs
  -c, --script       If present, data commands generate bash script to query Giveth's APIs
  -s, --sleep uint   Instructs the tool how long to sleep between invocations
  -x, --fmt string   One of [json|csv|txt]
  -v, --verbose      If present, certain commands will display extra data
  -h, --help         help for giveth

Use "giveth [command] --help" for more information about a command.
```

## Building / Running

This tool runs on Linux and Mac. There is no official Windows support. Some users have had success using WSL─you're on your own!

These instructions assume you can navigate the command line and edit configuration files.

0. Install dependencies
    - &#9745; [Install the latest version of Go](https://golang.org/doc/install).
    - &#9745; Install the other dependencies with your command line: `build-essential` `git` `cmake` `ninja` `python` `python-dev` `libcurl3-dev` `clang-format` `jq`.

1. Compile from the codebase
    ```shell
    git clone https://github.com/trueblocks/trueblocks-giveth
    cd trueblocks-giveth
    make
    ```

2. Add `./bin` to your shell `PATH`.

3. Testing
```shell
make test
```
## The API endpoint

```
curl -X 'GET' 'https://givback.develop.giveth.io/purpleList' -H 'accept: application/json'
```

## The Review Process

![image](https://user-images.githubusercontent.com/5417918/180751873-57d86b2c-fde9-4f6d-8e87-5e1795bbe24a.png)

