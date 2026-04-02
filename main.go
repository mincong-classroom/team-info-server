package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
)

const port = "8090" // avoids conflicts with the Spring PetClinic app (8080)

var teamPattern = regexp.MustCompile(`^(east|west|south|north)-[0-9]$`)

func validateTeam(team string) error {
	if !teamPattern.MatchString(team) {
		return fmt.Errorf("invalid team format: %q. Expected format: {region}-{digit} where region is one of: east, west, south, north", team)
	}
	return nil
}

type TeamInfo struct {
	Team        string            `json:"team"`
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

func main() {
	fmt.Println("Starting server...")

	fmt.Println("Validating environment variables...")
	var team string
	if value, ok := os.LookupEnv("TEAM"); ok {
		team = value
	}
	if team == "" {
		// This is intentional to let students practice troubleshooting errors in Kubernetes and
		// setting environment variables in the manifest.
		fmt.Printf("Environment variable %q is not set\n", "TEAM")
		fmt.Println("Exiting...")
		os.Exit(1)
	}

	if err := validateTeam(team); err != nil {
		fmt.Println(err)
		fmt.Println("Exiting...")
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received %s request to %s\n", r.Method, r.URL.Path)

		info := TeamInfo{
			Team: team,
			K8sLabels: map[string]string{
				"team": team,
			},
			GitRepo: fmt.Sprintf("https://github.com/mincong-classroom/k8s-%s", team),
			DockerRepos: []DockerRepo{
				{
					Id:      fmt.Sprintf("spring-petclinic-%s", team),            // ex: "spring-petclinic-east-1"
					Name:    fmt.Sprintf("Spring PetClinic Monolith (%s)", team), // ex: "Spring PetClinic Monolith (east-1)"
					RepoUrl: fmt.Sprintf("mincongclassroom/spring-petclinic-%s", team),
					WebUrl:  fmt.Sprintf("https://hub.docker.com/r/mincongclassroom/spring-petclinic-%s", team),
				},
				{
					Id:      fmt.Sprintf("spring-petclinic-api-gateway-%s", team),                   // ex: "spring-petclinic-api-gateway-east-1"
					Name:    fmt.Sprintf("Spring PetClinic Microservices - API Gateway (%s)", team), // ex: "Spring PetClinic Microservices - API Gateway (east-1)"
					RepoUrl: fmt.Sprintf("mincongclassroom/spring-petclinic-api-gateway-%s", team),
					WebUrl:  fmt.Sprintf("https://hub.docker.com/r/mincongclassroom/spring-petclinic-api-gateway-%s", team),
				},
				{
					Id:      fmt.Sprintf("spring-petclinic-customers-service-%s", team),
					Name:    fmt.Sprintf("Spring PetClinic Microservices - Customers Service (%s)", team),
					RepoUrl: fmt.Sprintf("mincongclassroom/spring-petclinic-customers-service-%s", team),
					WebUrl:  fmt.Sprintf("https://hub.docker.com/r/mincongclassroom/spring-petclinic-customers-service-%s", team),
				},
				{
					Id:      fmt.Sprintf("spring-petclinic-vets-service-%s", team),
					Name:    fmt.Sprintf("Spring PetClinic Microservices - Veterinarians Service (%s)", team),
					RepoUrl: fmt.Sprintf("mincongclassroom/spring-petclinic-vets-service-%s", team),
					WebUrl:  fmt.Sprintf("https://hub.docker.com/r/mincongclassroom/spring-petclinic-vets-service-%s", team),
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info)
	})

	fmt.Println("Server running on port " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
