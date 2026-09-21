package function

import "github.com/GoogleCloudPlatform/functions-framework-go/functions"

func init() {
	functions.CloudEvent("LoanHandler", LoanHandler)
}