workflow "Build and push docker image" {
  on = "push"
  resolves = ["Push image"]
}

action "Build image" {
  uses = "actions/docker/cli@c08a5fc9e0286844156fefff2c141072048141f6"
  args = "build -t whalesan/the-needle-dropped:latest ."
}

action "Push image" {
  uses = "actions/docker/login@c08a5fc9e0286844156fefff2c141072048141f6"
  needs = ["Build image"]
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
}
