import { useParams } from "@solidjs/router";
import { createMemo, createResource, For, JSX, Show, Suspense } from "solid-js";
import { Commit, pipelineSection } from "./types";
import CommitListItem from "./CommitListItem";

interface Project {
  // Define the structure of a project here
  id: string;
  name: string;
}

interface CommitGroup {
  name: string;
  commits: Commit[];
}

const fetchProject = async (id: string): Promise<Project> => {
  const response = await fetch(`/api/projects/${id}`);
  return response.json();
};

const fetchCommits = async (projectId: string): Promise<Commit[]> => {
  const response = await fetch(`/api/projects/${projectId}/commits`);
  const commits = (await response.json()) as Commit[];
  return commits.sort(
    (a, b) => new Date(b.date).getTime() - new Date(a.date).getTime(),
  );
};

export default function Project() {
  const params = useParams();
  const [project] = createResource(() => fetchProject(params.projectId));

  return (
    <Suspense>
      <Show when={project()?.name} fallback={<div>Project not found.</div>}>
        <main class="font-mono">
          <div class="container">
            <h1 data-label="projectName" class="py-10 text-lg font-bold">
              {project()?.name}
            </h1>
          </div>
          <Commits projectId={params.projectId} />
        </main>
      </Show>
    </Suspense>
  );
}

function Commits(props: { projectId: string }) {
  const [commits] = createResource(() => fetchCommits(props.projectId));

  const groupedCommits = createMemo(() => {
    const commitList = commits();
    if (!commitList) return {};

    const sections = {
      commitStage: { status: "ok" },
      acceptance: { status: "ok" },
      deploy: { status: "ok" },
    };

    const groups = {
      runningCommitStage: [] as Commit[],
      queuedAcceptanceStage: [] as Commit[],
      runningAcceptanceStage: [] as Commit[],
      queuedForDeploy: [] as Commit[],
      runningDeploy: [] as Commit[],
    }

    // find all the commit stage ones
    for (const commit of commitList) {
      if (commit.commitStageStatus === "passed") {
        break;
      }

      groups.runningCommitStage.push(commit);
    }

    // find all the queued acceptance stage ones
    for (const commit of commitList) {
      if (commit.acceptanceStageStatus === "started") {
        break;
      }

      groups.queuedAcceptanceStage.push(commit);
    }

    // find the running acceptance stage ones
    for (const commit of commitList) {
      if (commit.acceptanceStageStatus === "passed") {
        break;
      }

      groups.runningAcceptanceStage.push(commit);
    }
    
    // find the running acceptance stage ones
    for (const commit of commitList) {
      if (commit.acceptanceStageStatus === "passed") {
        break;
      }

      groups.runningAcceptanceStage.push(commit);
    }
  });

  return (
    <Suspense>
      <Show when={commits()} fallback={<div>Commits not found.</div>}>
        <div data-commits="loaded">
          <Section name="commit-stage">
            <Group name="Running commit stage:">
              <For each={commits()}>
                {(commit) => <CommitListItem commit={commit} />}
              </For>
            </Group>
          </Section>
          <Section name="acceptance-stage">
            <Group name="Running acceptance stage:">
              <For each={[]}>
                {(commit) => <CommitListItem commit={commit} />}
              </For>
            </Group>
          </Section>
          <Section name="deploy">
            <Group name="Running deploy:">
              <For each={[]}>
                {(commit) => <CommitListItem commit={commit} />}
              </For>
            </Group>
          </Section>
        </div>
      </Show>
    </Suspense>
  );
}

function Section(props: { children: JSX.Element; name: pipelineSection }) {
  return <section data-section={props.name}>{props.children}</section>;
}

function Group(props: { children: JSX.Element; name: string }) {
  return (
    <div>
      <h2 class="container text-sm">{props.name}:</h2>
      <ul>{props.children}</ul>
    </div>
  );
}
