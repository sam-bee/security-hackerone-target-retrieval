# HackerOne Target Retrieval

This is a limited use of an API client for retrieving particular information from the HackerOne researcher API.

In particular, in produces a list of all web applications currently in scope in any programme in HackerOne.

You will need a HackerOne researcher API token from their website to use this tool.

## Building

Type:

```
make build
```

## Testing

```
make test
```

## Configuring

You will need `HACKERONEAPICLIENT_USERNAME` and other settings on your machine. Either add these as environment
variables, or create a copy of `./config/.env.example` called `.env`, then pass the config file location to the
programme using the `--config` flag. See the example file for all necessary configuration.

## Running

```
cd bin
hackerone-api-client webapptargets -o ../data/test-output.csv -c ../config/.env
```
