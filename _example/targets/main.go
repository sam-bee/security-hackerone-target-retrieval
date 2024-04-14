package main

import (
	h1 "github.com/sam-bee/security-hackerone-target-retrieval"
)

func main () {
	example1()
	example2()
}

func example1() {

	output := "./targets-1.csv"
	user := "your-hackerone-username-here"
	token := "your-private-token-here"

	// Don't filter out anything
	filter := h1.NullFilter()

	filter := h1.NewFilter(programmeIsRelevant, targetIsRelevant)

	targetRetriever := h1.NewTargetRetriever(user, token, output, filter)
	targetRetriever.RetrieveTargets()
}

func example2() {

	output := "./targets-2.csv"
	user := "your-hackerone-username-here"
	token := "your-private-token-here"

	// Only show HackerOne Programmes which are open
	programmeIsRelevant := func (prog h1.Programme) bool {
		return prog.SubmissionState == "open"
	}

	// Only show Targets (within a Programme) which are websites where a bug bounty is available
	targetIsRelevant := func (target h1.Target) bool {
		return target.AssetType == "URL" && target.EligibleForSubmission && target.EligibleForBounty
	}

	filter := h1.NewFilter(programmeIsRelevant, targetIsRelevant)

	targetRetriever := h1.NewTargetRetriever(user, token, output, filter)
	targetRetriever.RetrieveTargets()
}
