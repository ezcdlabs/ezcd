import { createContext, createMemo, useContext, JSX } from "solid-js";
import { Commit } from "../types";
import { useParams } from "@solidjs/router";
import { createQuery } from "@tanstack/solid-query";
import groupCommits from "../../utils/groupCommits";

interface Project {
  // Define the structure of a project here
  id: string;
  name: string;
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

interface ProjectData {
  project: () => Project | undefined;
  groupedCommits: () => ReturnType<typeof groupCommits<Commit>>;
  failures: () => string[];
  isSuccess: () => boolean;
}

const ProjectDataContext = createContext<ProjectData | undefined>();

export default function DataLoader(props: { children: JSX.Element }) {
  const params = useParams();
  const projectQuery = createQuery(() => ({
    queryKey: ["project", params.projectId],
    queryFn: () => fetchProject(params.projectId),
  }));

  const commitsQuery = createQuery(() => ({
    queryKey: ["commits", params.projectId],
    queryFn: () => fetchCommits(params.projectId),
    refetchInterval: 5000,
    enabled: !!projectQuery.data,
  }));

  const groupedCommits = createMemo(() =>
    groupCommits(commitsQuery.data ?? []),
  );

  const failures = createMemo(() => {
    return groupedCommits()
      .filter((section) => section.status === "failing")
      .map((x) => x.name);
  });

  return (
    <ProjectDataContext.Provider
      value={{
        project: () => projectQuery.data,
        isSuccess: () => commitsQuery.isSuccess,
        groupedCommits,
        failures,
      }}
    >
      {props.children}
    </ProjectDataContext.Provider>
  );
}

export function useData() {
  const context = useContext(ProjectDataContext);
  if (!context) {
    throw new Error("useData must be used within a DataLoader");
  }
  return context;
}
