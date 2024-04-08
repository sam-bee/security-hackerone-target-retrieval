package targetretrieval

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/liamg/hackerone"
	"github.com/liamg/hackerone/pkg/api"
)

type targetFetchingWorkerPool struct {
	in     <-chan programme
	out    chan<- target
	api    *hackerone.API
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

func newTargetFetchingWorkerPool(in <-chan programme, out chan<- target) *targetFetchingWorkerPool {
	return &targetFetchingWorkerPool{
		in:  in,
		out: out,
	}
}

func (p *targetFetchingWorkerPool) targetFetchingWorker() {
	for prog := range p.in {
		targets := getTargetsForProgramme(p.api, prog, p.stdOut)
		for _, target := range targets {
			p.out <- target
		}
	}
}

func (p *targetFetchingWorkerPool) shutDown() {
	close(p.out)
}

func getTargetsForProgramme(h1 *hackerone.API, programme programme, stdOut io.Writer) []target {
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
			fmt.Fprintf(stdOut, "Discovered target %s %s\n", programme.handle, target.assetIdentifier)
		}
	}

	return targets
}
