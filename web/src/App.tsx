import type { Component } from 'solid-js';
import { createResource } from 'solid-js';
import { Router, Route } from "@solidjs/router";

const fetchHello = async () => {
  const response = await fetch('/api/hello');
  return response.text();
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
  const [hello] = createResource(fetchHello);

  return <div>Home page: {hello()}</div>
}
function Project() {
  return <div>Project page</div>
}

export default App;
