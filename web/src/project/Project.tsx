import { useParams } from "@solidjs/router";
import { createMemo, createResource, For, JSX, Show, Suspense } from "solid-js";
import { Commit, pipelineSection } from "./types";
import CommitListItem from "./CommitListItem";
import groupCommits from "./groupCommits";

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

  const groupedCommits = createMemo(() => groupCommits(commits() ?? []));

  return (
    <Suspense>
      <Show when={commits()} fallback={<div>Commits not found.</div>}>
        <div data-commits="loaded">
          {groupedCommits().map((section) => (
            <Section name={section.name}>
              <For each={section.groups}>
                {(group) => (
                  <Group name={group.name}>
                    <For each={group.commits}>
                      {(commit) => <CommitListItem commit={commit} />}
                    </For>
                  </Group>
                )}
              </For>
            </Section>
          ))}
        </div>
      </Show>
    </Suspense>
  );
}

function Section(props: { children: JSX.Element; name: string }) {
  return (
    <section
      class="border-t border-white-secondary py-4"
      data-section={props.name}
    >
      <h2 class="container mb-6 text-lg font-semibold">{props.name}</h2>
      {props.children}
    </section>
  );
}

function Group(props: { children: JSX.Element; name: string }) {
  return (
    <div>
      <h3 class="container text-sm">{props.name}</h3>
      <ul>{props.children}</ul>
    </div>
  );
}
