import { useParams } from "@solidjs/router";
import { createResource, Show, Suspense } from "solid-js";

const fetchProject = async (id: string) => {
  const response = await fetch(`/api/projects/${id}`);
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
        <ul>
          <li class="commit">commit1</li>
          <li class="commit">commit2</li>
          <li class="commit">commit3</li>
        </ul>
      </div>
    </Show>
  </Suspense>
}