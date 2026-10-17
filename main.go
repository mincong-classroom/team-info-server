package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const port = "8090" // avoids conflicts with the Spring PetClinic app (8080)

// teamPattern matches a team identifier: a single lowercase word made of letters
// only, with no digits, hyphens, underscores, or other characters (e.g. "red",
// "black", "teacher").
var teamPattern = regexp.MustCompile(`^[a-z]+$`)

func validateTeam(team string) error {
	if teamPattern.MatchString(team) {
		return nil
	}
	return fmt.Errorf("invalid team format: %q. Expected a single lowercase word with letters only, e.g. \"red\" or \"black\"", team)
}

// parseMembers turns the comma-separated TEAM_MEMBERS env var into a slice of names.
// Surrounding whitespace is trimmed and empty entries are dropped, so values
// like "Alice Doe, Bob Smith," yield ["Alice Doe", "Bob Smith"]. It always
// returns a non-nil slice so the JSON response renders an empty array (not
// null) when no members are configured.
func parseMembers(value string) []string {
	members := []string{}
	for _, part := range strings.Split(value, ",") {
		if name := strings.TrimSpace(part); name != "" {
			members = append(members, name)
		}
	}
	return members
}

type TeamInfo struct {
	Team        string            `json:"team"`
	Members     []string          `json:"members"`
	K8sLabels   map[string]string `json:"k8s_labels"`
	GitRepo     string            `json:"git_repo"`
	DockerRepos []DockerRepo      `json:"docker_repos"`
}

type DockerRepo struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	RepoUrl string `json:"repo_url"`
	WebUrl  string `json:"web_url"`
}

// newTeamInfo returns the metadata of a team: its Kubernetes labels, its Git repository, and the
// Docker repositories that the team pushes its images to.
func newTeamInfo(team string, members []string) TeamInfo {
	return TeamInfo{
		Team:    team,
		Members: members,
		K8sLabels: map[string]string{
			"team": team,
		},
		GitRepo: fmt.Sprintf("https://github.com/mincong-classroom/k8s-%s", team),
		DockerRepos: []DockerRepo{
			{
				Id:      fmt.Sprintf("spring-petclinic-%s", team),            // ex: "spring-petclinic-red"
				Name:    fmt.Sprintf("Spring PetClinic Monolith (%s)", team), // ex: "Spring PetClinic Monolith (red)"
				RepoUrl: fmt.Sprintf("mincongclassroom/spring-petclinic-%s", team),
				WebUrl:  fmt.Sprintf("https://hub.docker.com/r/mincongclassroom/spring-petclinic-%s", team),
			},
		},
	}
}

// teamInfoHandler answers every request with the metadata of the team, in JSON.
func teamInfoHandler(team string, members []string) http.HandlerFunc {
	info := newTeamInfo(team, members)
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received %s request to %s\n", r.Method, r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(info); err != nil {
			// The status is already sent: the client gets a truncated body, and the log says why.
			fmt.Printf("Failed to write the response: %v\n", err)
		}
	}
}

func main() {
	fmt.Println("Starting server...")

	fmt.Println("Validating environment variables...")
	var team string
	if value, ok := os.LookupEnv("TEAM_ID"); ok {
		team = value
	}
	if team == "" {
		// This is intentional to let students practice troubleshooting errors in Kubernetes and
		// setting environment variables in the manifest.
		fmt.Printf("Environment variable %q is not set\n", "TEAM_ID")
		fmt.Println("Exiting...")
		os.Exit(1)
	}

	if err := validateTeam(team); err != nil {
		fmt.Println(err)
		fmt.Println("Exiting...")
		os.Exit(1)
	}

	// TEAM_MEMBERS is optional (comma-separated). Unlike TEAM_ID, a missing value
	// is not an error: the response simply reports an empty members list.
	members := parseMembers(os.Getenv("TEAM_MEMBERS"))
	fmt.Printf("Team members: %v\n", members)

	http.HandleFunc("/", teamInfoHandler(team, members))

	fmt.Println("Server running on port " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
