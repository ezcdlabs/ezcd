import type { Component } from "solid-js";
import { createResource, For } from "solid-js";
import { Router, Route, A } from "@solidjs/router";
import Project from "./project/Project";
import logo from "./logo.svg";

const fetchProjects = async () => {
  const response = await fetch("/api/projects");
  return response.json();
};

const App: Component = () => {
  return (
    <Router>
      <Route path="/project/:projectId" component={Project} />
      <Route path="/" component={Home} />
    </Router>
  );
};

function Home() {
  const [projects] = createResource(fetchProjects);

  return (
    <main>
      <div class="container flex flex-col items-center py-10">
        <img src={logo} alt="EZCD Logo" class="h-20 w-20" />
        <h1 class="text-xl font-bold">EZCD</h1>
      </div>
      <ul class="container flex flex-col">
        <For each={projects()}>
          {(project: any) => (
            <A href={`/project/${project.id}`}>
              <li class="rounded-lg p-4 hover:bg-neutral-900">
                {project.name}
              </li>
            </A>
          )}
        </For>
      </ul>
    </main>
  );
}

export default App;
