package main //+integration

import (
	"testing"

	"github.com/weaveworks/flux/update"
)

func TestReleaseCommand_CLIConversion(t *testing.T) {
	for _, v := range []struct {
		args           []string
		expectedParams map[string]string
	}{
		{[]string{"--update-all-images", "--all"}, map[string]string{
			"service": string(update.ResourceSpecAll),
			"image":   string(update.ImageSpecLatest),
			"kind":    string(update.ReleaseKindExecute),
		}},
		{[]string{"--update-all-images", "--all", "--dry-run"}, map[string]string{
			"service": string(update.ResourceSpecAll),
			"image":   string(update.ImageSpecLatest),
			"kind":    string(update.ReleaseKindPlan),
		}},
		{[]string{"--update-image=alpine:latest", "--all"}, map[string]string{
			"service": string(update.ResourceSpecAll),
			"image":   "alpine:latest",
			"kind":    string(update.ReleaseKindExecute),
		}},
		{[]string{"--update-all-images", "--controller=deployment/flux"}, map[string]string{
			"service": "default:deployment/flux",
			"image":   string(update.ImageSpecLatest),
			"kind":    string(update.ReleaseKindExecute),
		}},
		{[]string{"--update-all-images", "--all", "--exclude=deployment/test,deployment/yeah"}, map[string]string{
			"service": string(update.ResourceSpecAll),
			"image":   string(update.ImageSpecLatest),
			"kind":    string(update.ReleaseKindExecute),
			"exclude": "default:deployment/test,default:deployment/yeah",
		}},
	} {
		svc := testArgs(t, v.args, false, "")

		// Check that PostRelease was called with correct args
		method := "UpdateImages"
		if calledURL(method, svc.requestHistory) == nil {
			t.Fatalf("Expecting fluxctl to request %q, but did not.", method)
		}
		vars := calledRequest(method, svc.requestHistory).Vars
		for kk, vv := range v.expectedParams {
			assertString(t, vv, vars[kk])
		}

		// Check that GetRelease was polled for status
		method = "JobStatus"
		if calledURL(method, svc.requestHistory) == nil {
			t.Fatalf("Expecting fluxctl to request %q, but did not.", method)
		}
	}
}

func TestReleaseCommand_InputFailures(t *testing.T) {
	for _, v := range []struct {
		args []string
		msg  string
	}{
		{[]string{}, "Should error when no args"},
		{[]string{"--all"}, "Should error when not specifying image spec"},
		{[]string{"--all", "--update-image=alpine"}, "Should error with invalid image spec"},
		{[]string{"--update-all-images"}, "Should error when not specifying controller spec"},
		{[]string{"--controller=invalid&controller", "--update-all-images"}, "Should error with invalid controller"},
		{[]string{"subcommand"}, "Should error when given subcommand"},
	} {
		testArgs(t, v.args, true, v.msg)
	}

}
