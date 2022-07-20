# trueblocks-giveth

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
