package targetretrieval

import (
	"fmt"
	"io"

	"github.com/liamg/hackerone"
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


func (t *target) StringSlice() []string {
	return []string{
		t.programme.handle,
		t.assetIdentifier,
		t.assetType,
		fmt.Sprintf("%t", t.eligibleForSubmission),
		fmt.Sprintf("%t", t.eligibleForBounty),
	}
}

func SearchForTargets(o OutputDestinationInterface, username string, token string, stdErr io.Writer, stdOut io.Writer) {

	h1 := hackerone.New(username, token)

	programmes := getProgrammes(h1, stdOut)
	programmes = filterRelevantProgrammes(programmes)

	err := o.Open()
	if (err != nil) {
		fmt.Fprintf(stdErr, "Error opening output file: %s\n", err)
		return
	}
	defer o.Close()

	wg := sync.WaitGroup{}

	for i := 0; i < len(programmes); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			programme := programmes[i]
			fmt.Printf("Getting structured scopes for programme %s\n", programme.handle)
			targets := getTargetsForProgramme(h1, programme, stdOut)
			targets = filterRelevantTargets(targets)
			writeTargetsToCsv(o, targets, stdErr)
		}()
	}
	wg.Wait()
}