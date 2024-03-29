workflow "Build and push docker image" {
  on = "push"
  resolves = ["Deploy on lighsail"]
}

action "Build image" {
  uses = "actions/docker/cli@c08a5fc9e0286844156fefff2c141072048141f6"
  args = "build -t whalesan/the-needle-dropped:latest ."
}

action "Login to hub" {
  uses = "actions/docker/login@c08a5fc9e0286844156fefff2c141072048141f6"
  needs = ["Build image"]
  secrets = [
    "DOCKER_PASSWORD",
    "DOCKER_USERNAME",
  ]
}

action "Push to hub " {
  uses = "actions/docker/cli@c08a5fc9e0286844156fefff2c141072048141f6"
  needs = ["Login to hub"]
  args = "push whalesan/the-needle-dropped:latest"
}

action "Deploy on lighsail" {
  uses = "maddox/actions/ssh@master"
  needs = ["Push to hub "]
  secrets = [
    "HOST",
    "USER",
    "PRIVATE_KEY",
    "PUBLIC_KEY",
  ]
  args = "docker-compose pull app && docker-compose up -d --no-deps app && docker image prune -f"
}
