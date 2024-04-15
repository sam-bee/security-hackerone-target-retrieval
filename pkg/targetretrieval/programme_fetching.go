package targetretrieval

import (
	"context"
	"fmt"
	"io"

	"github.com/sam-bee/security-hackerone-api-client"
	"github.com/sam-bee/security-hackerone-api-client/pkg/api"
)

func fetchProgrammes(h1 *hackerone.API, stdOut *io.Writer, out chan<- programme, filter func(programme) bool) {

	pageOptions := &api.PageOptions{
		PageNumber: 1,
		PageSize:   100,
	}

	nextPageNumber := 1
	var programmesFull []api.Program

	for pageNo := 1; nextPageNumber > 0; pageNo++ {
		pageOptions.PageNumber = pageNo
		var err error
		programmesFull, nextPageNumber, err = h1.Hackers.GetPrograms(context.TODO(), pageOptions)
		if err != nil {
			fmt.Fprintf(*stdOut, "Error fetching programmes: %s\n", err)
			continue
		}

		for _, programmeFull := range programmesFull {
			programme := programme{
				handle:          programmeFull.Attributes.Handle,
				submissionState: programmeFull.Attributes.SubmissionState,
			}
			if filter(programme) {
				out <- programme
			}
			fmt.Fprintf(*stdOut, "Discovered programme %s\n", programme.handle)
		}
	}

	close(out)
}
