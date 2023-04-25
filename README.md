# Twint Zero
Like Twint, but zero fat.

# Installation
1) `git clone https://github.com/twintproject/twint-zero`
2) `cd twint-zero`
3) `go mod init twint-zero`
4) `go mod tidy`

# Usage
- Without compiling: `go run main.go -Query $QUERY -Instance $INSTANCE -Format $FORMAT`
- If you compiled... well at this point you are supposed to know.

## CLI Arguments
1) `$QUERY`: [Here](https://github.com/igorbrigadir/twitter-advanced-search) you go.
2) `$INSTANCE`: [Here](https://github.com/zedeus/nitter/wiki/Instances) you go.
2) `$FORMAT`: "csv" or "json".
4) `$FileName`: csv output file name

# Questions/issues
> Sir, the bill is: five GitHub stars, two forks and one retweet.

That being quoted, feel free to reach out.

# License 
MIT

# Credits
[Francesco Poldi](https://twitter.com/noneprivacy)

[Simon Archer](https://mastodon.social/@archy_bold): JSON formatting and attachments parsing
