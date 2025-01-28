// Copyright 2022 OpenSSF Scorecard Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package checks

import (
	"os"

	"github.com/ossf/scorecard/v5/checker"
	"github.com/ossf/scorecard/v5/checks/evaluation"
	"github.com/ossf/scorecard/v5/checks/raw"
	sce "github.com/ossf/scorecard/v5/errors"
	"github.com/ossf/scorecard/v5/probes"
	"github.com/ossf/scorecard/v5/probes/zrunner"
)

const (
	// CheckWebHooks is the registered name for WebHooks.
	CheckWebHooks = "Webhooks"
)

//nolint:gochecknoinits
func init() {
	if err := registerCheck(CheckWebHooks, WebHooks, nil); err != nil {
		// this should never happen
		panic(err)
	}
}

// WebHooks run Webhooks check.
func WebHooks(c *checker.CheckRequest) checker.CheckResult {
	// TODO: remove this check when v6 is released
	_, enabled := os.LookupEnv("SCORECARD_EXPERIMENTAL")
	if !enabled {
		c.Dlogger.Warn(&checker.LogMessage{
			Text: "SCORECARD_EXPERIMENTAL is not set, not running the Webhook check",
		})

		e := sce.WithMessage(sce.ErrUnsupportedCheck, "SCORECARD_EXPERIMENTAL is not set, not running the Webhook check")
		return checker.CreateRuntimeErrorResult(CheckWebHooks, e)
	}

	rawData, err := raw.WebHook(c)
	if err != nil {
		e := sce.WithMessage(sce.ErrScorecardInternal, err.Error())
		return checker.CreateRuntimeErrorResult(CheckWebHooks, e)
	}

	// Set the raw results.
	pRawResults := getRawResults(c)
	pRawResults.WebhookResults = rawData

	// Evaluate the probes.
	findings, err := zrunner.Run(pRawResults, probes.Webhook)
	if err != nil {
		e := sce.WithMessage(sce.ErrScorecardInternal, err.Error())
		return checker.CreateRuntimeErrorResult(CheckWebHooks, e)
	}

	// Return the score evaluation.
	ret := evaluation.Webhooks(CheckWebHooks, findings, c.Dlogger)
	ret.Findings = findings
	return ret
}
