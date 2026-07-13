package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTeamInfoHandler(t *testing.T) {
	tests := []struct {
		name                string
		teamValue           string
		membersValue        []string
		wantCode            int
		wantTeam            string
		wantGit             string
		wantDockerRepoCount int
		wantMembers         []string
	}{
		{
			name:                "standard team",
			teamValue:           "red",
			membersValue:        []string{"Alice Doe", "Bob Smith"},
			wantCode:            http.StatusOK,
			wantTeam:            "red",
			wantGit:             "https://github.com/mincong-classroom/k8s-red",
			wantDockerRepoCount: 4,
			wantMembers:         []string{"Alice Doe", "Bob Smith"},
		},
		{
			name:                "different team",
			teamValue:           "black",
			membersValue:        []string{},
			wantCode:            http.StatusOK,
			wantTeam:            "black",
			wantGit:             "https://github.com/mincong-classroom/k8s-black",
			wantDockerRepoCount: 4,
			wantMembers:         []string{},
		},
		{
			name:                "teacher",
			teamValue:           "teacher",
			membersValue:        []string{"Mincong Huang"},
			wantCode:            http.StatusOK,
			wantTeam:            "teacher",
			wantGit:             "https://github.com/mincong-classroom/k8s-teacher",
			wantDockerRepoCount: 4,
			wantMembers:         []string{"Mincong Huang"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create handler with specific team value
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				info := TeamInfo{
					Team:    tt.teamValue,
					Members: tt.membersValue,
					K8sLabels: map[string]string{
						"team": tt.teamValue,
					},
					GitRepo: "https://github.com/mincong-classroom/k8s-" + tt.teamValue,
					DockerRepos: []DockerRepo{
						{
							Id:      fmt.Sprintf("spring-petclinic-%s", tt.teamValue),
							Name:    fmt.Sprintf("Spring PetClinic Monolith (%s)", tt.teamValue),
							RepoUrl: fmt.Sprintf("mincongclassroom/spring-petclinic-%s", tt.teamValue),
							WebUrl:  fmt.Sprintf("https://hub.docker.com/r/mincongclassroom/spring-petclinic-%s", tt.teamValue),
						},
						{
							Id:      fmt.Sprintf("spring-petclinic-api-gateway-%s", tt.teamValue),
							Name:    fmt.Sprintf("Spring PetClinic Microservices - API Gateway (%s)", tt.teamValue),
							RepoUrl: fmt.Sprintf("mincongclassroom/spring-petclinic-api-gateway-%s", tt.teamValue),
							WebUrl:  fmt.Sprintf("https://hub.docker.com/r/mincongclassroom/spring-petclinic-api-gateway-%s", tt.teamValue),
						},
						{
							Id:      fmt.Sprintf("spring-petclinic-customers-service-%s", tt.teamValue),
							Name:    fmt.Sprintf("Spring PetClinic Microservices - Customers Service (%s)", tt.teamValue),
							RepoUrl: fmt.Sprintf("mincongclassroom/spring-petclinic-customers-service-%s", tt.teamValue),
							WebUrl:  fmt.Sprintf("https://hub.docker.com/r/mincongclassroom/spring-petclinic-customers-service-%s", tt.teamValue),
						},
						{
							Id:      fmt.Sprintf("spring-petclinic-vets-service-%s", tt.teamValue),
							Name:    fmt.Sprintf("Spring PetClinic Microservices - Veterinarians Service (%s)", tt.teamValue),
							RepoUrl: fmt.Sprintf("mincongclassroom/spring-petclinic-vets-service-%s", tt.teamValue),
							WebUrl:  fmt.Sprintf("https://hub.docker.com/r/mincongclassroom/spring-petclinic-vets-service-%s", tt.teamValue),
						},
					},
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

			// Check docker repos count
			if len(info.DockerRepos) != tt.wantDockerRepoCount {
				t.Errorf("docker_repos count = %d, want %d", len(info.DockerRepos), tt.wantDockerRepoCount)
			}

			// Check k8s labels
			if info.K8sLabels["team"] != tt.wantTeam {
				t.Errorf("k8s_labels.team = %s, want %s", info.K8sLabels["team"], tt.wantTeam)
			}

			// Check members
			if len(info.Members) != len(tt.wantMembers) {
				t.Errorf("members count = %d, want %d", len(info.Members), len(tt.wantMembers))
			}
			for i := range tt.wantMembers {
				if i < len(info.Members) && info.Members[i] != tt.wantMembers[i] {
					t.Errorf("members[%d] = %q, want %q", i, info.Members[i], tt.wantMembers[i])
				}
			}
		})
	}
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

func TestDockerReposStructure(t *testing.T) {
	teamValue := "red"
	repos := []DockerRepo{
		{
			Id:      fmt.Sprintf("spring-petclinic-%s", teamValue),
			Name:    fmt.Sprintf("Spring PetClinic Monolith (%s)", teamValue),
			RepoUrl: fmt.Sprintf("mincongclassroom/spring-petclinic-%s", teamValue),
			WebUrl:  fmt.Sprintf("https://hub.docker.com/r/mincongclassroom/spring-petclinic-%s", teamValue),
		},
		{
			Id:      fmt.Sprintf("spring-petclinic-api-gateway-%s", teamValue),
			Name:    fmt.Sprintf("Spring PetClinic Microservices - API Gateway (%s)", teamValue),
			RepoUrl: fmt.Sprintf("mincongclassroom/spring-petclinic-api-gateway-%s", teamValue),
			WebUrl:  fmt.Sprintf("https://hub.docker.com/r/mincongclassroom/spring-petclinic-api-gateway-%s", teamValue),
		},
		{
			Id:      fmt.Sprintf("spring-petclinic-customers-service-%s", teamValue),
			Name:    fmt.Sprintf("Spring PetClinic Microservices - Customers Service (%s)", teamValue),
			RepoUrl: fmt.Sprintf("mincongclassroom/spring-petclinic-customers-service-%s", teamValue),
			WebUrl:  fmt.Sprintf("https://hub.docker.com/r/mincongclassroom/spring-petclinic-customers-service-%s", teamValue),
		},
		{
			Id:      fmt.Sprintf("spring-petclinic-vets-service-%s", teamValue),
			Name:    fmt.Sprintf("Spring PetClinic Microservices - Veterinarians Service (%s)", teamValue),
			RepoUrl: fmt.Sprintf("mincongclassroom/spring-petclinic-vets-service-%s", teamValue),
			WebUrl:  fmt.Sprintf("https://hub.docker.com/r/mincongclassroom/spring-petclinic-vets-service-%s", teamValue),
		},
	}

	expectedIds := []string{
		"spring-petclinic-red",
		"spring-petclinic-api-gateway-red",
		"spring-petclinic-customers-service-red",
		"spring-petclinic-vets-service-red",
	}

	expectedNames := []string{
		"Spring PetClinic Monolith (red)",
		"Spring PetClinic Microservices - API Gateway (red)",
		"Spring PetClinic Microservices - Customers Service (red)",
		"Spring PetClinic Microservices - Veterinarians Service (red)",
	}

	// Check count
	if len(repos) != 4 {
		t.Errorf("expected 4 docker repos, got %d", len(repos))
	}

	// Check IDs and names
	for i, repo := range repos {
		if repo.Id != expectedIds[i] {
			t.Errorf("repo[%d].Id = %q, want %q", i, repo.Id, expectedIds[i])
		}
		if repo.Name != expectedNames[i] {
			t.Errorf("repo[%d].Name = %q, want %q", i, repo.Name, expectedNames[i])
		}
		if repo.RepoUrl == "" {
			t.Errorf("repo[%d].RepoUrl is empty", i)
		}
		if repo.WebUrl == "" {
			t.Errorf("repo[%d].WebUrl is empty", i)
		}
	}
}

func TestDockerReposURLFormatting(t *testing.T) {
	tests := []struct {
		name          string
		teamValue     string
		expectedRepos []string
	}{
		{
			name:      "team red",
			teamValue: "red",
			expectedRepos: []string{
				"mincongclassroom/spring-petclinic-red",
				"mincongclassroom/spring-petclinic-api-gateway-red",
				"mincongclassroom/spring-petclinic-customers-service-red",
				"mincongclassroom/spring-petclinic-vets-service-red",
			},
		},
		{
			name:      "team black",
			teamValue: "black",
			expectedRepos: []string{
				"mincongclassroom/spring-petclinic-black",
				"mincongclassroom/spring-petclinic-api-gateway-black",
				"mincongclassroom/spring-petclinic-customers-service-black",
				"mincongclassroom/spring-petclinic-vets-service-black",
			},
		},
		{
			name:      "team green",
			teamValue: "green",
			expectedRepos: []string{
				"mincongclassroom/spring-petclinic-green",
				"mincongclassroom/spring-petclinic-api-gateway-green",
				"mincongclassroom/spring-petclinic-customers-service-green",
				"mincongclassroom/spring-petclinic-vets-service-green",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repos := []DockerRepo{
				{
					Id:      fmt.Sprintf("spring-petclinic-%s", tt.teamValue),
					Name:    fmt.Sprintf("Spring PetClinic Monolith (%s)", tt.teamValue),
					RepoUrl: fmt.Sprintf("mincongclassroom/spring-petclinic-%s", tt.teamValue),
					WebUrl:  fmt.Sprintf("https://hub.docker.com/r/mincongclassroom/spring-petclinic-%s", tt.teamValue),
				},
				{
					Id:      fmt.Sprintf("spring-petclinic-api-gateway-%s", tt.teamValue),
					Name:    fmt.Sprintf("Spring PetClinic Microservices - API Gateway (%s)", tt.teamValue),
					RepoUrl: fmt.Sprintf("mincongclassroom/spring-petclinic-api-gateway-%s", tt.teamValue),
					WebUrl:  fmt.Sprintf("https://hub.docker.com/r/mincongclassroom/spring-petclinic-api-gateway-%s", tt.teamValue),
				},
				{
					Id:      fmt.Sprintf("spring-petclinic-customers-service-%s", tt.teamValue),
					Name:    fmt.Sprintf("Spring PetClinic Microservices - Customers Service (%s)", tt.teamValue),
					RepoUrl: fmt.Sprintf("mincongclassroom/spring-petclinic-customers-service-%s", tt.teamValue),
					WebUrl:  fmt.Sprintf("https://hub.docker.com/r/mincongclassroom/spring-petclinic-customers-service-%s", tt.teamValue),
				},
				{
					Id:      fmt.Sprintf("spring-petclinic-vets-service-%s", tt.teamValue),
					Name:    fmt.Sprintf("Spring PetClinic Microservices - Veterinarians Service (%s)", tt.teamValue),
					RepoUrl: fmt.Sprintf("mincongclassroom/spring-petclinic-vets-service-%s", tt.teamValue),
					WebUrl:  fmt.Sprintf("https://hub.docker.com/r/mincongclassroom/spring-petclinic-vets-service-%s", tt.teamValue),
				},
			}

			if len(repos) != len(tt.expectedRepos) {
				t.Errorf("repos count = %d, want %d", len(repos), len(tt.expectedRepos))
			}

			for i, repo := range repos {
				if i < len(tt.expectedRepos) && repo.RepoUrl != tt.expectedRepos[i] {
					t.Errorf("repos[%d].RepoUrl = %q, want %q", i, repo.RepoUrl, tt.expectedRepos[i])
				}

				// Verify WebUrl format is correct
				expectedWebUrl := fmt.Sprintf("https://hub.docker.com/r/%s", repo.RepoUrl)
				if repo.WebUrl != expectedWebUrl {
					t.Errorf("repos[%d].WebUrl = %q, want %q", i, repo.WebUrl, expectedWebUrl)
				}
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
