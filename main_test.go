package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTeamInfoHandler(t *testing.T) {
	tests := []struct {
		name     string
		team     string
		members  []string
		wantBody string
	}{
		{
			name:    "team with members",
			team:    "red",
			members: []string{"Alice Doe", "Bob Smith"},
			wantBody: `{
				"team": "red",
				"members": ["Alice Doe", "Bob Smith"],
				"k8s_labels": {"team": "red"},
				"git_repo": "https://github.com/mincong-classroom/k8s-red",
				"docker_repos": [
					{
						"id": "spring-petclinic-red",
						"name": "Spring PetClinic Monolith (red)",
						"repo_url": "mincongclassroom/spring-petclinic-red",
						"web_url": "https://hub.docker.com/r/mincongclassroom/spring-petclinic-red"
					}
				]
			}`,
		},
		{
			name:    "team without members",
			team:    "black",
			members: parseMembers(""),
			wantBody: `{
				"team": "black",
				"members": [],
				"k8s_labels": {"team": "black"},
				"git_repo": "https://github.com/mincong-classroom/k8s-black",
				"docker_repos": [
					{
						"id": "spring-petclinic-black",
						"name": "Spring PetClinic Monolith (black)",
						"repo_url": "mincongclassroom/spring-petclinic-black",
						"web_url": "https://hub.docker.com/r/mincongclassroom/spring-petclinic-black"
					}
				]
			}`,
		},
		{
			name:    "teacher",
			team:    "teacher",
			members: []string{"Mincong Huang"},
			wantBody: `{
				"team": "teacher",
				"members": ["Mincong Huang"],
				"k8s_labels": {"team": "teacher"},
				"git_repo": "https://github.com/mincong-classroom/k8s-teacher",
				"docker_repos": [
					{
						"id": "spring-petclinic-teacher",
						"name": "Spring PetClinic Monolith (teacher)",
						"repo_url": "mincongclassroom/spring-petclinic-teacher",
						"web_url": "https://hub.docker.com/r/mincongclassroom/spring-petclinic-teacher"
					}
				]
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			teamInfoHandler(tt.team, tt.members).ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
			}

			if got := w.Header().Get("Content-Type"); got != "application/json" {
				t.Errorf("Content-Type = %s, want application/json", got)
			}

			// The whole body, keys included: the consumers of the response read them by name
			if got, want := compactJSON(t, w.Body.String()), compactJSON(t, tt.wantBody); got != want {
				t.Errorf("body = %s, want %s", got, want)
			}
		})
	}
}

// compactJSON returns the JSON document s without its insignificant whitespace, so that two
// documents compare as strings.
func compactJSON(t *testing.T, s string) string {
	t.Helper()
	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(s)); err != nil {
		t.Fatalf("json.Compact(%q) failed: %v", s, err)
	}
	return buf.String()
}

func TestParseMembers(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  []string
	}{
		{
			name:  "empty string yields empty slice",
			value: "",
			want:  []string{},
		},
		{
			name:  "single member",
			value: "Alice Doe",
			want:  []string{"Alice Doe"},
		},
		{
			name:  "multiple members",
			value: "Alice Doe,Bob Smith",
			want:  []string{"Alice Doe", "Bob Smith"},
		},
		{
			name:  "trims surrounding whitespace",
			value: " Alice Doe ,  Bob Smith ",
			want:  []string{"Alice Doe", "Bob Smith"},
		},
		{
			name:  "drops empty entries and trailing comma",
			value: "Alice Doe,,Bob Smith,",
			want:  []string{"Alice Doe", "Bob Smith"},
		},
		{
			name:  "whitespace-only yields empty slice",
			value: "  ,  ",
			want:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseMembers(tt.value)

			// A non-nil slice guarantees the JSON response is "[]" rather than "null".
			if got == nil {
				t.Fatalf("parseMembers(%q) returned nil, want non-nil slice", tt.value)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("parseMembers(%q) = %v, want %v", tt.value, got, tt.want)
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("parseMembers(%q)[%d] = %q, want %q", tt.value, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestMembersJSONMarshaling(t *testing.T) {
	// Empty members must marshal as an empty array, not null, so consumers can
	// safely treat the field as an array.
	empty := TeamInfo{Team: "red", Members: parseMembers("")}
	data, err := json.Marshal(empty)
	if err != nil {
		t.Fatalf("failed to marshal TeamInfo: %v", err)
	}
	if !strings.Contains(string(data), `"members":[]`) {
		t.Errorf("expected empty members to marshal as [], got %s", data)
	}

	// Populated members round-trip correctly.
	populated := TeamInfo{Team: "red", Members: []string{"Alice Doe", "Bob Smith"}}
	data, err = json.Marshal(populated)
	if err != nil {
		t.Fatalf("failed to marshal TeamInfo: %v", err)
	}
	var unmarshaled TeamInfo
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal TeamInfo: %v", err)
	}
	if len(unmarshaled.Members) != 2 || unmarshaled.Members[0] != "Alice Doe" || unmarshaled.Members[1] != "Bob Smith" {
		t.Errorf("unmarshaled members = %v, want [Alice Doe Bob Smith]", unmarshaled.Members)
	}
}

func TestTeamInfoStructure(t *testing.T) {
	info := TeamInfo{
		Team: "red",
		K8sLabels: map[string]string{
			"team": "red",
		},
		GitRepo: "https://example.com/repo",
		DockerRepos: []DockerRepo{
			{
				Id:      "test-repo",
				Name:    "Test Repository",
				RepoUrl: "example.com/test",
				WebUrl:  "https://example.com/docker",
			},
		},
	}

	// Marshal to JSON and verify it works
	data, err := json.Marshal(info)
	if err != nil {
		t.Errorf("failed to marshal TeamInfo: %v", err)
	}

	// Unmarshal and verify structure
	var unmarshaled TeamInfo
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("failed to unmarshal TeamInfo: %v", err)
	}

	if unmarshaled.Team != info.Team {
		t.Errorf("unmarshaled team = %s, want %s", unmarshaled.Team, info.Team)
	}

	if len(unmarshaled.DockerRepos) != 1 {
		t.Errorf("unmarshaled docker_repos count = %d, want 1", len(unmarshaled.DockerRepos))
	}
}

func TestValidateTeam(t *testing.T) {
	tests := []struct {
		name      string
		teamValue string
		wantError bool
	}{
		// Valid: a single lowercase word made of letters only.
		{
			name:      "valid single word red",
			teamValue: "red",
			wantError: false,
		},
		{
			name:      "valid single word black",
			teamValue: "black",
			wantError: false,
		},
		{
			name:      "valid single letter",
			teamValue: "a",
			wantError: false,
		},
		{
			name:      "valid reserved teacher",
			teamValue: "teacher",
			wantError: false,
		},
		// Invalid: anything that is not a single lowercase word.
		{
			name:      "invalid legacy region-digit format",
			teamValue: "east-1",
			wantError: true,
		},
		{
			name:      "invalid capitalized word",
			teamValue: "Red",
			wantError: true,
		},
		{
			name:      "invalid all uppercase",
			teamValue: "RED",
			wantError: true,
		},
		{
			name:      "invalid hyphen",
			teamValue: "red-team",
			wantError: true,
		},
		{
			name:      "invalid underscore",
			teamValue: "red_team",
			wantError: true,
		},
		{
			name:      "invalid trailing digit",
			teamValue: "red1",
			wantError: true,
		},
		{
			name:      "invalid internal space",
			teamValue: "red team",
			wantError: true,
		},
		{
			name:      "invalid empty string",
			teamValue: "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTeam(tt.teamValue)
			if (err != nil) != tt.wantError {
				t.Errorf("validateTeam(%q) error = %v, wantError %v", tt.teamValue, err, tt.wantError)
			}
		})
	}
}

func TestDockerReposJSONMarshaling(t *testing.T) {
	info := TeamInfo{
		Team: "red",
		K8sLabels: map[string]string{
			"team": "red",
		},
		GitRepo: "https://github.com/mincong-classroom/k8s-red",
		DockerRepos: []DockerRepo{
			{
				Id:      "spring-petclinic-red",
				Name:    "Spring PetClinic Monolith (red)",
				RepoUrl: "mincongclassroom/spring-petclinic-red",
				WebUrl:  "https://hub.docker.com/r/mincongclassroom/spring-petclinic-red",
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(info)
	if err != nil {
		t.Errorf("failed to marshal TeamInfo: %v", err)
	}

	// Unmarshal and verify
	var unmarshaled TeamInfo
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("failed to unmarshal TeamInfo: %v", err)
	}

	// Verify docker_repos key is present in JSON
	if len(unmarshaled.DockerRepos) != 1 {
		t.Errorf("unmarshaled docker_repos count = %d, want 1", len(unmarshaled.DockerRepos))
	}

	if unmarshaled.DockerRepos[0].Id != "spring-petclinic-red" {
		t.Errorf("unmarshaled docker_repos[0].id = %q, want %q", unmarshaled.DockerRepos[0].Id, "spring-petclinic-red")
	}

	if unmarshaled.DockerRepos[0].Name != "Spring PetClinic Monolith (red)" {
		t.Errorf("unmarshaled docker_repos[0].name = %q, want %q", unmarshaled.DockerRepos[0].Name, "Spring PetClinic Monolith (red)")
	}
}
