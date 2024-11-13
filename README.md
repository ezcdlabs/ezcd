<div align="center">
  <a href="https://github.com/ezcdlabs/ezcd/">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="./logo.svg">
      <img alt="EZCD logo" src="./logo.svg" height="128">
    </picture>
  </a>
  <h1>EZCD</h1>


</div>

EZCD (pronounced 'Easy C.D.') makes is easy to create optimal CI/CD pipelines using modern best practices from Continuous Delivery and DevOps. EZCD integrates with your existing CI/CD system (GitHub Actions, GitLab Pipelines, Azure DevOps, CircleCI, Google Cloud Build...) and it can easily be deployed to anywhere that supports Docker containers.

# Overview

Optimized CI/CD pipelines lead to shorter delivery times and faster feedback.

__EZCD visualises the status of your pipeline__ - The first step of a CI/CD pipeline is to run unit tests and build a new release candidate. The EZCD dashboard visualizes the state of your pipeline, helping you prioritize keeping it green.

__EZCD optimizes the acceptance stage__ - The acceptance stage, where a release candidate is deployed to a production-like environment and tested, is often the slowest. EZCD helps you skip ahead to the most recent release candidate, avoiding an ever-growing backlog.

__EZCD tracks lead times__ - To improve lead times, we need to measure them. EZCD does this for you.

<!-- ## Getting Started -->
<!-- 


- **Dashboard**: Run the dashboard in Docker, e.g., on Render.com, Heroku, Kubernetes (k8s).

- **CLI**: Use the CLI from your CI/CD pipeline such as GitHub Actions, etc.

## Usage

### Dashboard

To run the dashboard, you need a PostgreSQL database. Pass the connection string as the environment variable `EZCD_DATABASE_URL`.

### CLI

Use the CLI in your CI/CD pipeline to interact with the dashboard and manage your deployments.

You can use the GitHub Action, or you can download EZCD from our releases page
```
  commit-stage:
    runs-on: ubuntu-latest
    steps:
      - name: Setup EZCD CLI using github action
        uses: ezcdlabs/ezcd@main
        with:
          version: 0.1.0

      - name: Commit stage started
        run: ezcd-cli commit-stage-started

      - name: Test and build
        run: ./your-test-and-build.sh
      
      - name: Commit stage passed
        run: ezcd-cli commit-stage-passed
      
      - name: Commit stage failed
        if: !success()
        run: ezcd-cli commit-stage-failed
``` -->