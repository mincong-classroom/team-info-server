# team-info-server

Team Info Server is a simple Go web server that returns team metadata in JSON format. It provides the team name, team members, Kubernetes labels, Git repository, and Docker repository information. It serves as a teaching tool for students to learn about environment variables, namespaces, and service discovery in Kubernetes.

- **Source code:** [`mincong-classroom/team-info-server`](https://github.com/mincong-classroom/team-info-server) on GitHub
- **Container image:** [`mincongclassroom/team-info-server`](https://hub.docker.com/r/mincongclassroom/team-info-server) on Docker Hub

## Core Configuration

This server is configured **entirely through environment variables** — there are no config files or command-line flags. This mirrors how applications are configured in Kubernetes, where you set these values in your Deployment manifest.

If this is your first time here, start with this table. There is exactly **one required variable** (`TEAM_ID`) and **one optional variable** (`TEAM_MEMBERS`):

| Variable       | Required | Description                                                                                                         | Example                |
| -------------- | -------- | ----------------------------------------------------------------------------------------------------------------- | ---------------------- |
| `TEAM_ID`      | Yes      | Your team's identifier. Must match `{region}-{digit}` (see [Team ID](#team-id)). The server refuses to start without a valid value. | `east-1`               |
| `TEAM_MEMBERS` | No       | Comma-separated list of member names (see [Team Members](#team-members)). Defaults to an empty list.               | `Alice Doe, Bob Smith` |

A complete example that sets both variables:

```bash
TEAM_ID=east-1 TEAM_MEMBERS="Alice Doe, Bob Smith" go run main.go
```

The two sections below explain each variable in detail.

### Team ID

`TEAM_ID` is **required**. The server validates it against the format `{region}-{digit}`, where:

- **region**: One of `east`, `west`, `south`, or `north`
- **digit**: A single digit from `0` to `9`

Valid examples: `east-1`, `west-2`, `south-0`, `north-9`

Invalid examples: `East-1`, `east-1-1`, `east-12`, `central-1`

If `TEAM_ID` is missing or does not match this format, the server prints an error and exits. This is intentional: troubleshooting the missing variable and fixing it in the Kubernetes manifest is part of the exercise.

The value is echoed back in the response as `team`, used to build the `k8s_labels`, and substituted into the Git and Docker repository URLs.

### Team Members

`TEAM_MEMBERS` is **optional** and complements `TEAM_ID` by naming who is on the team. It is a comma-separated list of names; surrounding whitespace is trimmed and empty entries are ignored:

```bash
TEAM_ID=east-1 TEAM_MEMBERS="Alice Doe, Bob Smith" go run main.go
```

The names are returned in the `members` array of the response. Unlike `TEAM_ID`, `TEAM_MEMBERS` is optional — when it is not set, `members` is an empty array (`[]`), never `null`, so consumers can always treat it as a list.

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
  "members": ["Alice Doe", "Bob Smith"],
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

## Container Image

Released versions are published to Docker Hub as [`mincongclassroom/team-info-server`](https://hub.docker.com/r/mincongclassroom/team-info-server). The same [Core Configuration](#core-configuration) applies — pass the environment variables with `-e`:

```bash
docker run -e TEAM_ID=east-1 -e TEAM_MEMBERS="Alice Doe, Bob Smith" -p 8090:8090 mincongclassroom/team-info-server
```

The image is built from this repository. Its source and documentation always point back to [github.com/mincong-classroom/team-info-server](https://github.com/mincong-classroom/team-info-server).

## Design Notes

Here are some choices made:

* We use `Go` because it is lightweight. It does not require any dependencies to bootstrap a new web server.
* We use a lightweight Docker image "scratch" to reduce the size of the image.
* We use port 8090 to avoid conflicts with the other exercises based on the Spring PetClinic project.
* We return JSON to provide structured team metadata that students can parse and use in their applications.
