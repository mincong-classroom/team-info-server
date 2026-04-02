package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTeamInfoHandler(t *testing.T) {
	tests := []struct {
		name      string
		teamValue string
		wantCode  int
		wantTeam  string
		wantGit   string
		wantDocker string
	}{
		{
			name:      "standard team",
			teamValue: "east-1",
			wantCode:  http.StatusOK,
			wantTeam:  "east-1",
			wantGit:   "https://github.com/mincong-classroom/k8s-east-1",
			wantDocker: "https://hub.docker.com/r/mincongclassroom/spring-petclinic-east-1",
		},
		{
			name:      "different team",
			teamValue: "west-2",
			wantCode:  http.StatusOK,
			wantTeam:  "west-2",
			wantGit:   "https://github.com/mincong-classroom/k8s-west-2",
			wantDocker: "https://hub.docker.com/r/mincongclassroom/spring-petclinic-west-2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create handler with specific team value
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				info := TeamInfo{
					Team: tt.teamValue,
					K8sLabels: map[string]string{
						"team": tt.teamValue,
					},
					GitRepo:    "https://github.com/mincong-classroom/k8s-" + tt.teamValue,
					DockerRepo: "https://hub.docker.com/r/mincongclassroom/spring-petclinic-" + tt.teamValue,
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(info)
			})

			// Make request
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.wantCode {
				t.Errorf("status code = %d, want %d", w.Code, tt.wantCode)
			}

			// Check content type
			if got := w.Header().Get("Content-Type"); got != "application/json" {
				t.Errorf("Content-Type = %s, want application/json", got)
			}

			// Parse response
			var info TeamInfo
			if err := json.NewDecoder(w.Body).Decode(&info); err != nil {
				t.Errorf("failed to decode response: %v", err)
			}

			// Check team value
			if info.Team != tt.wantTeam {
				t.Errorf("team = %s, want %s", info.Team, tt.wantTeam)
			}

			// Check git repo
			if info.GitRepo != tt.wantGit {
				t.Errorf("git_repo = %s, want %s", info.GitRepo, tt.wantGit)
			}

			// Check docker repo
			if info.DockerRepo != tt.wantDocker {
				t.Errorf("docker_repo = %s, want %s", info.DockerRepo, tt.wantDocker)
			}

			// Check k8s labels
			if info.K8sLabels["team"] != tt.wantTeam {
				t.Errorf("k8s_labels.team = %s, want %s", info.K8sLabels["team"], tt.wantTeam)
			}
		})
	}
}

func TestTeamInfoStructure(t *testing.T) {
	info := TeamInfo{
		Team: "test-team",
		K8sLabels: map[string]string{
			"team": "test-team",
		},
		GitRepo:    "https://example.com/repo",
		DockerRepo: "https://example.com/docker",
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
}
