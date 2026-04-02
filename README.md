# team-info-server

Team Info Server is a simple web server that returns team metadata in JSON format. It provides team name, Kubernetes labels, Git repository, and Docker repository information. It serves as a teaching tool for students to learn about environment variables, namespaces, and service discovery in Kubernetes.

The server expects the `TEAM` environment variable to be set (e.g., `TEAM=east-1`). Based on this value, it constructs and returns team-related information.

Example response:
```json
{
  "team": "east-1",
  "k8s_labels": {
    "team": "east-1"
  },
  "git_repo": "https://github.com/mincong-classroom/k8s-east-1",
  "docker_repo": "https://hub.docker.com/r/mincongclassroom/spring-petclinic-east-1"
}
```

Here are some choices made:

* We use `Go` because it is lightweight. It does not require any dependencies to bootstrap a new web server.
* We use a lightweight Docker image "scratch" to reduce the size of the image.
* We use port 8090 to avoid conflicts with the other exercises based on the Spring PetClinic project.
* We return JSON to provide structured team metadata that students can parse and use in their applications.