package github

import (
	"context"
	"testing"

	"github.com/techreloaded-ar/ARchetipo/cli/internal/config"
)

// TestInitializeConnector_HappyPath verifies the gh call sequence used by
// initialize_connector resolves repo + project metadata into a SetupInfo.
func TestInitializeConnector_HappyPath(t *testing.T) {
	m := newMock(t).
		on("repo view --json", `{
			"id":"R_abc","owner":{"login":"acme"},"name":"web","nameWithOwner":"acme/web"
		}`).
		on("project list --owner acme", `{
			"projects":[{"number":4,"id":"PVT_kw","title":"web Backlog","url":"https://gh/p/4"}]
		}`).
		on("project field-list 4", `{
			"fields":[
				{"id":"FID_status","name":"Status","type":"SINGLE_SELECT","options":[
					{"id":"OPT_todo","name":"TODO"},{"id":"OPT_planned","name":"PLANNED"}
				]},
				{"id":"FID_pri","name":"Priority","type":"SINGLE_SELECT","options":[
					{"id":"OPT_high","name":"HIGH"}
				]},
				{"id":"FID_sp","name":"Story Points","type":"NUMBER"}
			]
		}`)

	cfg := config.Default()
	cfg.Connector = config.ConnectorGitHub
	c := NewWithRunner(cfg, m)

	info, err := c.InitializeConnector(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if info.Repo == nil || info.Repo.Slug != "acme/web" {
		t.Errorf("repo not resolved: %+v", info.Repo)
	}
	if info.Project == nil || info.Project.Number != 4 {
		t.Errorf("project not resolved: %+v", info.Project)
	}
	if info.Project.Fields.StatusOptions["TODO"] != "OPT_todo" {
		t.Errorf("status option lost: %+v", info.Project.Fields.StatusOptions)
	}
	if !m.calledWithPrefix("repo view") {
		t.Errorf("expected gh repo view to be called")
	}
}

// TestProjectPreference_TitleContainsBacklog covers the second tier of the
// preference pipeline: when no exact title match exists, prefer one that
// contains "Backlog".
func TestProjectPreference_TitleContainsBacklog(t *testing.T) {
	m := newMock(t).
		on("repo view --json", `{"id":"R","owner":{"login":"o"},"name":"n","nameWithOwner":"o/n"}`).
		on("project list --owner o", `{"projects":[
			{"number":7,"id":"PVT7","title":"Tracking","url":""},
			{"number":3,"id":"PVT3","title":"Sprint Backlog","url":""}
		]}`).
		on("project field-list 3", `{"fields":[]}`)

	cfg := config.Default()
	cfg.Connector = config.ConnectorGitHub
	c := NewWithRunner(cfg, m)

	info, err := c.InitializeConnector(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if info.Project.Number != 3 {
		t.Errorf("expected project 3 (contains Backlog), got %d", info.Project.Number)
	}
}
