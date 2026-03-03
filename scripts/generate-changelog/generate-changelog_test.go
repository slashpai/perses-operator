// Copyright The Perses Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/perses/perses/scripts/pkg/changelog"
	"github.com/stretchr/testify/assert"
)

func TestGenerateChangelog(t *testing.T) {
	now := time.Now()
	title := fmt.Sprintf("## 0.20.0 / %s", now.Format("2006-01-02"))
	testSuite := []struct {
		title    string
		clog     *changelog.Changelog
		expected string
	}{
		{
			title:    "empty changelog",
			clog:     &changelog.Changelog{},
			expected: fmt.Sprintf("%s\n%s\n", title, ""),
		},
		{
			title: "changelog with every entry",
			clog: &changelog.Changelog{
				Features: []string{"Discard Changes Confirmation Dialog (#834)"},
				Enhancements: []string{"Variable UX fixes (#842)",
					"legend options editor UX improvements (#845)",
					"Make it possible to adjust the height of the time range controls (#829)",
				},
				BugFixes:        []string{"Fix time units display, allow decimalPlaces to be used (#837)"},
				BreakingChanges: []string{"legend.position now required in time series panel (#848)"},
				Docs:            []string{"Complete documentation about the API. (#1471) (##1479) (##1483) (#1490) (#1491) (#1500)"},
				Unknown:         []string{"Use exact versions for internal npm dependencies (#846)", "Support snapshot UI releases (#844)"},
			},
			expected: fmt.Sprintf("%s\n%s", title, `
- [FEATURE] Discard Changes Confirmation Dialog (#834)
- [ENHANCEMENT] Variable UX fixes (#842)
- [ENHANCEMENT] legend options editor UX improvements (#845)
- [ENHANCEMENT] Make it possible to adjust the height of the time range controls (#829)
- [BUGFIX] Fix time units display, allow decimalPlaces to be used (#837)
- [BREAKINGCHANGE] legend.position now required in time series panel (#848)
- [DOC] Complete documentation about the API. (#1471) (##1479) (##1483) (#1490) (#1491) (#1500)

[//]: <UNKNOWN ENTRIES. Release shepherd, please review the following list and categorize them or remove them>

- [UNKNOWN] Use exact versions for internal npm dependencies (#846)
- [UNKNOWN] Support snapshot UI releases (#844)
`),
		},
	}
	for _, test := range testSuite {
		t.Run(test.title, func(t *testing.T) {
			assert.Equal(t, test.expected, generateChangelog(test.clog, "0.20.0"))
		})
	}
}

func TestPreprocessEntry(t *testing.T) {
	testSuite := []struct {
		title    string
		entry    string
		expected string
		kept     bool
	}{
		{
			title:    "already categorized entry is unchanged",
			entry:    "abc1234 [FEATURE] Add new widget (#100)",
			expected: "abc1234 [FEATURE] Add new widget (#100)",
			kept:     true,
		},
		{
			title: "build(deps) entry is filtered",
			entry: "abc1234 build(deps): Bump github.com/foo from 1.0 to 2.0 (#101)",
			kept:  false,
		},
		{
			title: "Build(deps) entry is filtered case-insensitively",
			entry: "abc1234 Build(deps): Bump sigs.k8s.io/controller-runtime (#102)",
			kept:  false,
		},
		{
			title:    "feat: prefix maps to FEATURE",
			entry:    "abc1234 feat: add dashboard tags (#200)",
			expected: "abc1234 [FEATURE] add dashboard tags (#200)",
			kept:     true,
		},
		{
			title:    "feat(scope): prefix maps to FEATURE",
			entry:    "abc1234 feat(api): add new CRD field (#201)",
			expected: "abc1234 [FEATURE] add new CRD field (#201)",
			kept:     true,
		},
		{
			title:    "fix: prefix maps to BUGFIX",
			entry:    "abc1234 fix: correct logging issues (#202)",
			expected: "abc1234 [BUGFIX] correct logging issues (#202)",
			kept:     true,
		},
		{
			title:    "docs: prefix maps to DOC",
			entry:    "abc1234 docs: update README (#203)",
			expected: "abc1234 [DOC] update README (#203)",
			kept:     true,
		},
		{
			title:    "refactor: prefix maps to ENHANCEMENT",
			entry:    "abc1234 refactor: simplify reconciler logic (#204)",
			expected: "abc1234 [ENHANCEMENT] simplify reconciler logic (#204)",
			kept:     true,
		},
		{
			title:    "perf: prefix maps to ENHANCEMENT",
			entry:    "abc1234 perf: reduce memory allocations (#208)",
			expected: "abc1234 [ENHANCEMENT] reduce memory allocations (#208)",
			kept:     true,
		},
		{
			title: "chore: prefix is filtered",
			entry: "abc1234 chore: update tooling (#205)",
			kept:  false,
		},
		{
			title: "ci: prefix is filtered",
			entry: "abc1234 ci: fix linting job (#206)",
			kept:  false,
		},
		{
			title: "test: prefix is filtered",
			entry: "abc1234 test: add unit tests (#207)",
			kept:  false,
		},
		{
			title:    "unrecognized prefix is passed through",
			entry:    "abc1234 something unusual happened",
			expected: "abc1234 something unusual happened",
			kept:     true,
		},
	}
	for _, test := range testSuite {
		t.Run(test.title, func(t *testing.T) {
			result, ok := preprocessEntry(test.entry)
			assert.Equal(t, test.kept, ok)
			if ok {
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestPreprocessEntries(t *testing.T) {
	entries := []string{
		"aaa1111 [FEATURE] Add widget (#1)",
		"bbb2222 build(deps): Bump foo from 1.0 to 2.0 (#2)",
		"ccc3333 feat: add bar (#3)",
		"ddd4444 chore: update deps (#4)",
		"eee5555 fix: resolve crash (#5)",
		"fff6666 docs: improve guide (#6)",
	}
	expected := []string{
		"aaa1111 [FEATURE] Add widget (#1)",
		"ccc3333 [FEATURE] add bar (#3)",
		"eee5555 [BUGFIX] resolve crash (#5)",
		"fff6666 [DOC] improve guide (#6)",
	}
	assert.Equal(t, expected, preprocessEntries(entries))
}
