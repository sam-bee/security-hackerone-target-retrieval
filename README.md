# HackerOne Target Retrieval

## Introduction

This library is built on a HackerOne API client. It retrieves all the Programmes and Targets that match your
criteria, and puts them in a .csv file. You can use it to assist in target acquisition, or for automated
vulnerability scanning.

You will need your own credentials for the HackerOne researcher API. If you don't have a token already, go to
the [HackerOne settings page](https://hackerone.com/settings/api_token/edit) and generate one. Keep it secret.

## Example Use

You will find other examples in the `./_examples` folder, but here is a quick demo:

```go
	// Output to CSV file
	output := "./targets-1.csv"

	// Add your API creds
	user := "your-hackerone-username-here"
	token := "your-private-token-here"

	// Only show HackerOne Programmes which are open
	programmeIsRelevant := func (prog h1.Programme) bool {
		return prog.SubmissionState == "open"
	}

	// Only show Targets (within a Programme) which are websites where a bug bounty is available
	targetIsRelevant := func (target h1.Target) bool {
		return target.AssetType == "URL" && target.EligibleForBounty
	}

	// Now you have a bug bounty Programme filter and a Target filter
	filter := h1.NewFilter(programmeIsRelevant, targetIsRelevant)

	// Get all the relevant targets from the API
	targetRetriever := h1.NewTargetRetriever(user, token, output, filter)
	targetRetriever.RetrieveTargets()
```

## A Note on Rate Limiting

At time of writing, the API rate limit is 600 queries/minute. This may change. Running the tool several times
in quick succession may stop it from outputting for a short time. Enhance your calm.

## Output

The columns of the output CSV file are these:

```
Programme
PssetIdentifier
AssetType
EligibleForSubmission
EligibleForBounty
```

The author wishes you the best of luck, and happy hacking.
