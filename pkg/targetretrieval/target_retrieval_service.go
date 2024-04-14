package targetretrieval

import (
	"github.com/liamg/hackerone"
	"io"
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

func SearchForTargets(o OutputDestinationInterface, api *hackerone.API, filter Filter, stdErr io.Writer, stdOut io.Writer) {
	programmesCh := make(chan programme)
	targetsCh := make(chan target)
	signalCh := make(chan bool)

	go fetchProgrammes(api, stdOut, programmesCh, filter.ProgrammeIsRelevant)
	runTargetFetchingWorkerPool(100, api, stdOut, programmesCh, targetsCh, filter.TargetIsRelevant)
	go writeTargetsToCsv(o, targetsCh, signalCh, stdErr)

	<-signalCh
}

func runTargetFetchingWorkerPool(noOfWorkers int, api *hackerone.API, stdOut io.Writer, in <-chan programme, out chan<- target, filter func(target) bool) {
	pool := newTargetFetchingWorkerPool(in, out, api, filter, stdOut)
	pool.run(noOfWorkers)
}
