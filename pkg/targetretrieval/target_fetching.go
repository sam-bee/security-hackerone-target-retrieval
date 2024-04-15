package targetretrieval

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/sam-bee/security-hackerone-api-client"
	"github.com/sam-bee/security-hackerone-api-client/pkg/api"
)

type targetFetchingWorkerPool struct {
	in     <-chan programme
	out    chan<- target
	api    *hackerone.API
	filter func(target) bool
	stdOut io.Writer
}

func (p *targetFetchingWorkerPool) run(noOfWorkers int) {
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(noOfWorkers)
		for i := 0; i < noOfWorkers; i++ {
			go func() {
				defer wg.Done()
				p.targetFetchingWorker()
			}()
		}
		wg.Wait()
		p.shutDown()
	}()
}

func newTargetFetchingWorkerPool(in <-chan programme, out chan<- target, api *hackerone.API, filter func(target) bool, stdOut *io.Writer) *targetFetchingWorkerPool {
	return &targetFetchingWorkerPool{
		in:     in,
		out:    out,
		api:    api,
		filter: filter,
		stdOut: *stdOut,
	}
}

func (p *targetFetchingWorkerPool) targetFetchingWorker() {
	for prog := range p.in {
		targets := getRelevantTargetsForProgramme(p.api, prog, p.filter, p.stdOut)
		for _, target := range targets {
			p.out <- target
		}
	}
}

func (p *targetFetchingWorkerPool) shutDown() {
	close(p.out)
}

func getRelevantTargetsForProgramme(h1 *hackerone.API, programme programme, filter func(target) bool, stdOut io.Writer) []target {
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

			if filter(target) {
				targets = append(targets, target)
				fmt.Fprintf(stdOut, "Discovered target %s %s\n", programme.handle, target.assetIdentifier)
			}
		}
	}

	return targets
}
