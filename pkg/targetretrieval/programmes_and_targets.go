package targetretrieval

type Filter struct {
	ProgrammeIsRelevant func(programme programme) bool
	TargetIsRelevant    func(target target) bool
}

func NewFilter(
	programmeIsRelevant func(programme programme) bool,
	targetIsRelevant func(target target) bool,
) Filter {
	return Filter{
		ProgrammeIsRelevant: programmeIsRelevant,
		TargetIsRelevant:    targetIsRelevant,
	}
}

func NullFilter() Filter {
	return Filter{
		ProgrammeIsRelevant: func(programme programme) bool {
			return true
		},
		TargetIsRelevant: func(target target) bool {
			return true
		},
	}
}

func UrlsWithBounties() Filter {
	programmeOpen := func(programme programme) bool {
		return programme.submissionState == "open"
	}
	targetIsUrlWithBounty := func(target target) bool {
		return target.assetType == "URL" && target.eligibleForSubmission && target.eligibleForBounty
	}
	return NewFilter(programmeOpen, targetIsUrlWithBounty)
}
