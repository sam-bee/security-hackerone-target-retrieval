package hackeronetargetretrieval

import (
	"io"
	"os"

	"github.com/sam-bee/security-hackerone-target-retrieval/pkg/targetretrieval"

	hackerone "github.com/sam-bee/security-hackerone-api-client"
)

type targetRetriever struct {
	api    *hackerone.API
	filter targetretrieval.Filter
	output targetretrieval.OutputDestinationInterface
	stdOut io.Writer
	stdErr io.Writer
}

func NewTargetRetriever(user string, token string, outputPath string, filter targetretrieval.Filter) *targetRetriever {
	api := hackerone.New(user, token)
	output := targetretrieval.OutputFile{Path: outputPath}
	stdOut := io.Writer(os.Stdout)
	stdErr := io.Writer(os.Stderr)
	return NewConfiguredTargetRetriever(*api, filter, &output, &stdOut, &stdErr)
}

func NewConfiguredTargetRetriever(api hackerone.API, filter targetretrieval.Filter, output targetretrieval.OutputDestinationInterface, stdOut *io.Writer, stdErr *io.Writer) *targetRetriever {
	return &targetRetriever{
		api:    &api,
		filter: filter,
		output: output,
		stdOut: *stdOut,
		stdErr: *stdErr,
	}
}

func (t *targetRetriever) RetrieveTargets() {
	targetretrieval.RetrieveTargets(t.api, t.filter, t.output, &t.stdOut, &t.stdErr)
}
