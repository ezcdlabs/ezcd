<div align="center">
  <a href="https://github.com/ezcdlabs/ezcd/">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="./logo.svg">
      <img alt="EZCD logo" src="./logo.svg" height="128">
    </picture>
  </a>
  <h1>EZCD</h1>
  
  ![CI Stage](https://github.com/ezcdlabs/ezcd/actions/workflows/ci.yml/badge.svg)
  ![Release](https://github.com/ezcdlabs/ezcd/actions/workflows/release.yml/badge.svg)
</div>

EZCD (pronounced 'Easy C.D.') is a CI/CD enhancement tool that integrates with your existing pipeline to provide faster feedback and accelerate your development

![alt text](images/dashboard-example.png)

# Why do you need EZCD?

EZCD helps you create effective CI/CD pipelines, enabling you to achieve shorter delivery times and produce higher quality software.

## 1. Get fast feedback using the EZCD Dashboard

Imagine working without an IDE's error checking or syntax highlighting. We need fast feedback loops to help us write the best code; the same applies to creating reliable and high-performing pipelines.

The EZCD dasboard provides a real-time, transparent view of your pipeline. You can instantly see which commits are in each stage of your workflow.

EZCD tracks the health of your pipeline, and alerts you to failures as soon as they occur. It keeps a notice in place so you can act quickly to unblock your team.

## 2. Improve your pipeline over time with DORA metrics

EZCD helps you continuously improve your processes by tracking the following metrics:

- **Lead time**: EZCD tracks the lead times of your commits. You can see how long it took before each commit was deployed to production.
- **Deployment frequency**: EZCD tracks your deployments and shows you how frequently you are delivering changes to your users.

By comparing this over time you can see how changes to your workflow impact the speed and efficiency of your pipeline.

## 3. Smart coordination of your acceptance test runs

Before you can deploy your changes, they need to run through the two main stages of a CI/CD pipeline.

1. Commit stage
   - Runs on every commit and can run in parallel on multiple commits at once
   - Builds or creates a release candidate from the source
   - Runs unit tests and other checks that can be completed quickly
2. Acceptance stage
   - Deploys a release candidate to a production-like environment and runs acceptance tests
   - Runs other checks that are too slow to complete in the commit stage
   - Can only run on one release candidate at a time

If you schedule an acceptance run after every successful release candidate you will create an ever-growing backlog of work.

EZCD tracks the successful release candidates and ensures that your acceptance stage jumps ahead to the newest release candidate. 

<details>
  <summary>Why isn't this possible natively in GitHub Actions?</summary>
  <blockquote>
  GitHub Actions actually comes quite close - closer than GitLab Pipelines.

  You can restrict the parallelism of your acceptance workflow and configure it to queue up a new one after each sucessful commit-stage.
  But the queuing order is based on when the acceptance workflow was triggered, which can be different to the order in which the commits were created - imagine if your most recent commit completed its commit stage really quickly, then older commits would trigger the acceptance workflow and cancel the workflow for your latest commit: Not good!
  </blockquote>
</details>


## 4. Prioritisation of trunk based continuous delivery workflows

High performing teams need to use modern workflows like Trunk Based Development and Continuous Delivery.

This can make it hard for your team to adopt tools that assume your workflow includes Pull Requests and Feature Branches.

EZCD is the other way around and works best for modern trunk based workflows.

# Getting Started

```mermaid
graph LR;
  A["Your CI/CD pipeline"] --> B["EZCD CLI"];
  B --> C["Postgres"];
  D["EZCD Dashboard
  Docker Container"] --> C;
```

## 1. Set up postgres

EZCD needs an external Postgres database to track information about your pipeline.

Set up the database using the schema file [./schema/projects_and_commits.sql](./schema/projects_and_commits.sql)

## 2. Add the EZCD_DATABASE_URL environment variable to your pipeline

The format should match:

```
postgresql://[user[:password]@][host][:port][/dbname][?param1=value1&...]
```

For example:

```
postgresql://user:password@your-ezcd-database-server:5432/ezcd
```

## 3. Install the CLI commands into your pipeline

<details>
  <summary>GitHub Actions</summary>

```yaml
env:
  PROJECT_ID: your_project_id_here

commit-stage:
  runs-on: ubuntu-latest
  steps:
    - name: Setup EZCD CLI using github action
      uses: ezcdlabs/ezcd@main
      with:
        version: 0.1.0

    - name: Commit stage started
      run: ezcd-cli commit-stage-started

    - name: Commit stage started
      run: ezcd-cli commit-stage-started --project $PROJECT_ID --hash "${{ github.sha }}"
        --author-name "${{ github.actor }}"
        --author-email "${{ github.actor }}"
        --message "${{ github.event.head_commit.message }}"
        --date "${{ github.event.head_commit.timestamp }}"

    - name: Test and build
      run: ./your-test-and-build.sh

    - name: Commit stage passed
      run: ezcd-cli commit-stage-passed --project $PROJECT_ID --hash "${{ github.sha }}"

    - name: Commit stage failed
      if: !success()
      run: ezcd-cli commit-stage-failed --project $PROJECT_ID --hash "${{ github.sha }}"
```

</details>

<details>
  <summary>Other CI/CD Pipelines</summary>

```sh
# Install EZCD CLI
version=0.1.0
url=https://github.com/ezcdlabs/ezcd/releases/download/$version/ezcd_linux_amd64.tar.gz
curl -L -o ezcdcli.tar.gz $url
tar -xzf ezcdcli.tar.gz -C /usr/local/bin
chmod +x /usr/local/bin/ezcd-cli
```

</details>

## 4. Run the EZCD dashboard using Docker

```sh
docker run -d \
  --name ezcd-dashboard \
  -e EZCD_DATABASE_URL=postgres://username:password@hostname:port/dbname \
  -p 3923:3923 \
  ghcr.io/ezcdlabs/ezcd-server:latest
```

<!-- 1. Create a persistent Postgres database that can be accessed from your CI/CD pipeline
2. Install the EZCD CLI into your pipeline

EZCD consists of two components:

- **EZCD CLI** - designed to be used in your CI/CD pipeline
- **EZCD Dashboard** - a docker image that you can

/**
 * ezcd is a tool that consists of two main components: the CLI and the dashboard.
 *
 * The CLI is designed to be used in your CI/CD pipeline. You can install the CLI binary by downloading it from the releases page, or you can use the GitHub Action 'ezcdlabs/ezcd'.
 *
 * The dashboard is published as a Docker image and provides a graphical interface for managing your deployments.
 */ -->

<!--
EZCD is more than just a monitoring dashboard you check when something goes wrong; it's your homepage. It's the central hub you return to between each commit, keeping you focused on delivering value to your users.  -->

<!--
EZCD is more than just a monitoring dashboard you check when something goes wrong; it's your homepage. It's the central hub you return to between each commit, keeping you focused on delivering value to your users. With EZCD, you can see how you're doing—whether you're making as many small commits as you thought and if you're actually achieving the lead times you expect. It provides a reassuring space where you can quickly verify that everything is going to plan, ensuring continuous progress and peace of mind. -->

<!-- **EZCD visualises the status of your pipeline**






Fast feedback about our pipeline itself




We can achieve shorter delivery times and create higher quality software when we have optimised and healthy CI/CD pipelines.

We can create higher quality software with shorter delivery times

High performing teams produce higher quality software with shorter


High performing teams use high-quality, healthy, optimised CI/CD pipelines to achieve shorter delivery times and they produce higher quality,




High performing dev teams need a high performing CI/CD pipeline. When we keep it healthy and monitor lead times we can achieve much shorter delivery times.


ide makes it effortless to ensure our code is syntactically correct through fast feedback

- Monitor your lead times
- Coordinate your acceptance test runs
- React quickly to




- high performing teams need high performing pipeline
- faster feedback matters
- keep it green and keep an eye on the lead times is feedback that guides you to make your pipeline better
- when we can see what the issue is, we can fix it with ease (like red squiggly or ide error)
- EZCD is here to make these things much much easier!!


to provide observability, monitoring and optimisation, helping you run pipelines faster,


 that gives you fast feedback on your CI/CD pipelines.

EZCD integrates with your existing CI/CD system



makes is easy to create optimal CI/CD pipelines using modern best practices from Continuous Delivery and DevOps.

EZCD integrates with your existing CI/CD system (GitHub Actions, GitLab Pipelines, Azure DevOps, CircleCI, Google Cloud Build...) to coordinate your pipeline runs and visualises the state of your CI/CD pipeline in the EZCD dashboard.



# Overview

Fast feedback matters. Imagine working without an IDE's error checking or syntax highlighting. IDEs accelerate your development by providing real-time feedback on your code in an easy-to-consume yet comprehensive manner, helping you stay on track and quickly address any issues as they arise.

EZCD gives you fast feedback on your CI/CD pipeline.

EZCD provides a real-time, transparent view of your pipeline, structured to offer a clear overview while surfacing the most important details. This enables you to:
 - act quickly when there are failing builds, tests, or deploys
 - monitor your lead times
 -


EZCD is designed to support teams that need to move quickly.
EZCD works with trunk based development. In trunk based development there are no pull requests, which is good because pull requests introduce unnecessary time waste into the workflow, but it is often difficult to achieve good DX


, enabling you to act quickly to failing builds or tests,

EZCD provides a real-time, transparent view of your pipeline, structured to offer both a clear overview and surface the most important details all in one place.


So you can:
- Act quickly when there is a build or test failure
- Monitor the lead times of your commits (commit date to when the commit is first included in a deploy to production)
-


EZCD provides rich information about your pipeline .

EZCD shows you feedback about your pipeline


EZCD helps you adopt

Support trunk based development


# Why do you need EZCD?

The priority for Product teams is to get working software into the hands of their users.

IDE
DX

Delays

Waste:

- failing builds and tests blocking pipelines
- wasteful or poor scheduling of acceptance testing and deploys
-

Nothing is more important in software product development than getting working software into the hands of your users.

Optimizing your CI/CD pipeline will help your team achieve shorter delivery times and get faster feedback!

The EZCD dasboard clearly shows all of the information you need to monitor your CI/CD pipeline,

makes your team efficient

Keep your pipeline clear

The EZCD dashboard makes the state of your CI/CD clear, enabling your team to:

- address failures quickly
- keep track of lead times
-

Central place to visualise your pipeline:

- lead times
- blocked or red pipeline
- commits waiting to be deployed

Optimized CI/CD pipelines lead to shorter delivery times and faster feedback.

## Trunk based development

Trunk based development is true Continuous Integration, where your team makes

**EZCD visualises the status of your pipeline** - The first step of a CI/CD pipeline is to run unit tests and build a new release candidate. The EZCD dashboard visualizes the state of your pipeline, helping you prioritize keeping it green.

**EZCD optimizes the acceptance stage** - The acceptance stage, where a release candidate is deployed to a production-like environment and tested, is often the slowest. EZCD helps you skip ahead to the most recent release candidate, avoiding an ever-growing backlog.

**EZCD tracks lead times** - To improve lead times, we need to measure them. EZCD does this for you. -->

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
