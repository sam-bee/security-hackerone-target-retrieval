package targetretrieval

import (
	"context"
	"fmt"
	"io"
	"github.com/liamg/hackerone"
	"github.com/liamg/hackerone/pkg/api"
)

func getProgrammes(h1 *hackerone.API, stdOut io.Writer) []programme {

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
			fmt.Fprintf(stdOut, "Discovered programme %s\n", programme.handle)
		}
	}

	return programmes
}
