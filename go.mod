module github.com/sam-bee/security-hackerone-target-retrieval

go 1.22.1

require github.com/liamg/hackerone v0.0.8

require (
	github.com/avast/retry-go v3.0.0+incompatible // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
)

replace github.com/liamg/hackerone => github.com/sam-bee/security_hackerone-client v0.0.8
