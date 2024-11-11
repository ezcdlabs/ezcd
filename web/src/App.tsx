import type { Component } from 'solid-js';
import { createResource, For } from 'solid-js';
import { Router, Route } from "@solidjs/router";
import Project from './pages/Project';

const fetchProjects = async () => {
  const response = await fetch('/api/projects');
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

  return <main>
    <h1>
      Your Projects
    </h1>
    <ul>
      <For each={projects()}>
        {(project: any) => (
          <li>
            <a href={`/project/${project.id}`}>Project: {project.name}</a>
          </li>
        )}
      </For>
    </ul>
  </main>
}


export default App;
