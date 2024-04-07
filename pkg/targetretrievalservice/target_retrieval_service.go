package targetretrievalservice

import (
	"context"
	"fmt"
	"hackeroneapiclient/pkg/csvfiles"
	"io"
	"sync"

	"github.com/liamg/hackerone"
	"github.com/liamg/hackerone/pkg/api"
)

type programme struct {
	handle          string
	submissionState string
}

type target struct {
	programme             programme
	assetIdentifier       string
	assetType             string
	eligibleForSubmission bool
	eligibleForBounty     bool
}

func SearchForWebApps(o csvfiles.OutputDestination, username string, token string, stdErr io.Writer, stdOut io.Writer) {

	h1 := hackerone.New(username, token)

	programmes := getProgrammes(h1)
	programmes = filterRelevantProgrammes(programmes)

	o.Open()
	defer o.Close()

	wg := sync.WaitGroup{}

	for i := 0; i < len(programmes); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			programme := programmes[i]
			fmt.Printf("Getting structured scopes for programme %s\n", programme.handle)
			targets := getTargetsForProgramme(h1, programme)
			targets = filterRelevantTargets(targets)
			writeTargetsToCsv(o, targets)
		}()
	}
	wg.Wait()
}

func getProgrammes(h1 *hackerone.API) []programme {

	programmes := []programme{}

	pageOptions := &api.PageOptions{
		PageNumber: 1,
		PageSize:   100,
	}

	nextPageNumber := 1
	var programmesFull []api.Program

	for pageNo := 1; nextPageNumber > 0; pageNo++ {
		pageOptions.PageNumber = pageNo
		programmesFull, nextPageNumber, _ = h1.Hackers.GetPrograms(context.TODO(), pageOptions)

		for _, programmeFull := range programmesFull {
			programme := programme{
				handle:          programmeFull.Attributes.Handle,
				submissionState: programmeFull.Attributes.SubmissionState,
			}
			programmes = append(programmes, programme)
			fmt.Printf("Discovered programme %s\n", programme.handle)
		}
	}

	return programmes
}

func filterRelevantProgrammes(all []programme) []programme {
	relevantProgrammes := []programme{}
	for _, programme := range all {
		if programme.submissionState == "open" {
			relevantProgrammes = append(relevantProgrammes, programme)
		}
	}
	return relevantProgrammes
}

func getTargetsForProgramme(h1 *hackerone.API, programme programme) []target {
	targets := []target{}

	pageOptions := &api.PageOptions{
		PageNumber: 1,
		PageSize:   100,
	}

	nextPageNumber := 1
	var structuredScopes []api.StructuredScope

	for pageNo := 1; nextPageNumber > 0; pageNo++ {
		pageOptions.PageNumber = pageNo
		structuredScopes, nextPageNumber, _ = h1.Hackers.GetStructuredScopes(context.TODO(), programme.handle, pageOptions)

		for _, structuredScope := range structuredScopes {
			target := target{
				programme:             programme,
				assetIdentifier:       structuredScope.Attributes.AssetIdentifier,
				assetType:             structuredScope.Attributes.AssetType,
				eligibleForSubmission: structuredScope.Attributes.EligibleForSubmission,
				eligibleForBounty:     structuredScope.Attributes.EligibleForBounty,
			}
			targets = append(targets, target)
			fmt.Printf("Discovered target %s %s\n", programme.handle, target.assetIdentifier)
		}
	}

	return targets
}

func filterRelevantTargets(targets []target) []target {
	relevantTargets := []target{}
	for _, target := range targets {
		if target.assetType == "URL" && target.eligibleForSubmission && target.eligibleForBounty {
			relevantTargets = append(relevantTargets, target)
		}
	}
	return relevantTargets
}

func writeTargetsToCsv(o csvfiles.OutputDestination, targets []target) {
	for _, target := range targets {
		o.Write(target.StringSlice())
	}
}

func (t *target) StringSlice() []string {
	return []string{
		t.programme.handle,
		t.assetIdentifier,
		t.assetType,
		fmt.Sprintf("%t", t.eligibleForSubmission),
		fmt.Sprintf("%t", t.eligibleForBounty),
	}
}