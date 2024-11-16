import { useParams } from "@solidjs/router";
import { createResource, For, Show, Suspense } from "solid-js";

interface Project {
  // Define the structure of a project here
  id: string;
  name: string;
}

interface Commit {
  // Define the structure of a commit here
  hash: string;
  message: string;
  authorName: string;
  authorEmail: string;
  date: string;

  commitStageStatus: string;
  commitStageStartedAt: string;
  commitStageCompletedAt: string;
}

const fetchProject = async (id: string): Promise<Project> => {
  const response = await fetch(`/api/projects/${id}`);
  return response.json();
};

const fetchCommits = async (projectId: string): Promise<Commit[]> => {
  const response = await fetch(`/api/projects/${projectId}/commits`);
  return response.json();
};

export default function Project() {

  const params = useParams();
  const [project] = createResource(() => fetchProject(params.projectId));

  return <Suspense fallback={<div>Loading...</div>}>
    <Show when={project()?.name} fallback={<div>Project not found.</div>}>
      <div>Project page
        <h1 class="project-name">{project()?.name}</h1>
        <p>Project ID: {params.projectId}</p>
        <Commits projectId={params.projectId} />
      </div>
    </Show>
  </Suspense>
}

function Commits(props: { projectId: string }) {
  const [commits] = createResource(() => fetchCommits(props.projectId));
  return <Suspense fallback={<div>Loading...</div>}>
    <Show when={commits()} fallback={<div>Commits not found.</div>}>
      <ul data-commits="loaded">
        <For each={commits()}>
          {commit => (
            <li data-commit={commit.hash}>
              Hash: <span data-label="commitHash">{commit.hash}</span>
              &nbsp;
              Author: <span data-label="commitAuthorName">{commit.authorName}</span>(<span data-label="commitAuthorEmail">{commit.authorEmail}</span>)
              &nbsp;
              Message: <span data-label="commitMessage">{commit.message}</span>
              &nbsp;
              Message: <span data-label="commitDate">{commit.date}</span>
              &nbsp;
              <br />
              Commit Stage: <span data-label="commitStageStatus" data-value={commit.commitStageStatus}>{commit.commitStageStatus}</span>
              &nbsp;
              Start: <span data-label="commitStageStartedAt" data-value={commit.commitStageStartedAt}>{commit.commitStageStartedAt}</span>
              &nbsp;
              Completed: <span data-label="commitStageCompletedAt" data-value={commit.commitStageCompletedAt}>{commit.commitStageCompletedAt}</span>
            </li>
          )}
        </For>
      </ul>
    </Show>
  </Suspense>

}