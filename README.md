# team-info-server

Team Info Server is a simple web server that returns team metadata in JSON format. It provides team name, Kubernetes labels, Git repository, and Docker repository information. It serves as a teaching tool for students to learn about environment variables, namespaces, and service discovery in Kubernetes.

## Team Validation

The server expects the `TEAM` environment variable to be set with the format `{region}-{digit}`, where:
- **region**: One of `east`, `west`, `south`, or `north`
- **digit**: A single digit from 0-9

Valid examples: `east-1`, `west-2`, `south-0`, `north-9`

Invalid examples: `East-1`, `east-1-1`, `east-12`, `central-1`

## Docker Repositories

The server returns an array of Docker repositories supporting both monolithic and microservices architectures:
- **Spring PetClinic Monolith** — Traditional monolithic architecture
- **API Gateway** — Microservices entry point
- **Customers Service** — Microservice for customer data
- **Vets Service** — Microservice for veterinarian data

Each repository is identified with:
- `id`: Unique identifier
- `name`: Human-readable name with team context
- `repo_url`: Repository reference (e.g., `mincongclassroom/spring-petclinic-east-1`)
- `web_url`: Full Docker Hub URL

Example response:
```json
{
  "team": "east-1",
  "k8s_labels": {
    "team": "east-1"
  },
  "git_repo": "https://github.com/mincong-classroom/k8s-east-1",
  "docker_repos": [
    {
      "id": "spring-petclinic-east-1",
      "name": "Spring PetClinic Monolith (east-1)",
      "repo_url": "mincongclassroom/spring-petclinic-east-1",
      "web_url": "https://hub.docker.com/r/mincongclassroom/spring-petclinic-east-1"
    },
    {
      "id": "spring-petclinic-api-gateway-east-1",
      "name": "Spring PetClinic Microservices - API Gateway (east-1)",
      "repo_url": "mincongclassroom/spring-petclinic-api-gateway-east-1",
      "web_url": "https://hub.docker.com/r/mincongclassroom/spring-petclinic-api-gateway-east-1"
    },
    {
      "id": "spring-petclinic-customers-service-east-1",
      "name": "Spring PetClinic Microservices - Customers Service (east-1)",
      "repo_url": "mincongclassroom/spring-petclinic-customers-service-east-1",
      "web_url": "https://hub.docker.com/r/mincongclassroom/spring-petclinic-customers-service-east-1"
    },
    {
      "id": "spring-petclinic-vets-service-east-1",
      "name": "Spring PetClinic Microservices - Veterinarians Service (east-1)",
      "repo_url": "mincongclassroom/spring-petclinic-vets-service-east-1",
      "web_url": "https://hub.docker.com/r/mincongclassroom/spring-petclinic-vets-service-east-1"
    }
  ]
}
```

Here are some choices made:

* We use `Go` because it is lightweight. It does not require any dependencies to bootstrap a new web server.
* We use a lightweight Docker image "scratch" to reduce the size of the image.
* We use port 8090 to avoid conflicts with the other exercises based on the Spring PetClinic project.
* We return JSON to provide structured team metadata that students can parse and use in their applications.