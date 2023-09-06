# Twint Zero
Like Twint, but zero fat.

# First things first

Users are invited **to not** scrape public instances, that will cause a bad experience for some users. Instead, you are invited to setup your own custom Nitter instance.
Thank you, and enjoy!

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
2) `$INSTANCE`: [Setup your own](https://github.com/zedeus/nitter/#installation).
2) `$FORMAT`: "csv" or "json".

# Questions/issues
> Sir, the bill is: five GitHub stars, two forks and one retweet.

That being quoted, feel free to reach out.

# License 
MIT

# Credits
[Francesco Poldi](https://twitter.com/noneprivacy)

[Simon Archer](https://mastodon.social/@archy_bold): JSON formatting and attachments parsing

[Julian](https://github.com/juste97): quoted tweet and its metadata fields
