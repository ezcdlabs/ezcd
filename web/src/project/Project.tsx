import { useParams } from "@solidjs/router";
import { For, JSX, Show, Suspense } from "solid-js";
import CommitListItem from "./components/CommitListItem";
import classNames from "../utils/classNames";
import logo from "../logo.svg";
import getMedianLeadTime from "../utils/getMedianLeadTime";
import getShortDurationFormatFromSeconds from "../utils/getShortDurationFormatFromSeconds";
import DataLoader, { useData } from "./components/DataLoader";

interface Project {
  // Define the structure of a project here
  id: string;
  name: string;
}

export default function Project() {
  return (
    <DataLoader>
      <Suspense>
        <ProjectHeader />
        <ProjectBody />
      </Suspense>
    </DataLoader>
  );
}

function ProjectHeader() {
  const projectId = useParams().projectId;

  return (
    <div class="fixed inset-x-0 top-0 z-20 flex h-11 shrink-0 items-center bg-neutral-900/75 backdrop-blur-sm">
      <div class="container flex items-center gap-4">
        <a href="/" class="group flex items-center">
          <svg
            class="h-5 w-5 fill-neutral-300 group-hover:fill-neutral-100"
            aria-hidden="true"
            viewBox="0 0 24 24"
            data-testid="ArrowBackIcon"
            aria-label="fontSize small"
          >
            <path d="M20 11H7.83l5.59-5.59L12 4l-8 8 8 8 1.41-1.41L7.83 13H20z"></path>
          </svg>
          <div class="h-8 w-8 rounded-full border-2 border-neutral-300 bg-neutral-800 hover:border-neutral-100">
            <img src={logo} alt="EZCD Logo" class="h-full w-full" />
          </div>
        </a>

        <h1 data-label="projectName" class="font-semibold">
          {projectId}
        </h1>
      </div>
    </div>
  );
}

function ProjectBody() {
  const params = useParams();
  const data = useData();

  return (
    <Suspense>
      <Show
        when={data.project()?.name}
        fallback={<div>Project not found.</div>}
      >
        <main class="flex h-screen flex-col pt-11">
          <div class="fixed inset-x-0 top-0 z-20 flex h-11 shrink-0 items-center bg-neutral-900/75 backdrop-blur-sm">
            <div class="container flex items-center gap-4">
              <a href="/" class="group flex items-center">
                <svg
                  class="h-5 w-5 fill-neutral-300 group-hover:fill-neutral-100"
                  aria-hidden="true"
                  viewBox="0 0 24 24"
                  data-testid="ArrowBackIcon"
                  aria-label="fontSize small"
                >
                  <path d="M20 11H7.83l5.59-5.59L12 4l-8 8 8 8 1.41-1.41L7.83 13H20z"></path>
                </svg>
                <div class="h-8 w-8 rounded-full border-2 border-neutral-300 bg-neutral-800 hover:border-neutral-100">
                  <img src={logo} alt="EZCD Logo" class="h-full w-full" />
                </div>
              </a>

              <h1 data-label="projectName" class="font-semibold">
                {data.project()?.name}
              </h1>
            </div>
          </div>

          <div class="flex min-h-0 grow flex-col">
            <Commits projectId={params.projectId} />
          </div>
        </main>
      </Show>
    </Suspense>
  );
}

function Commits(props: { projectId: string }) {
  const { isSuccess, groupedCommits, failures, project } = useData();

  return (
    <Suspense>
      <Show
        when={isSuccess() && groupedCommits()?.length !== 0}
        fallback={
          <div class="flex grow items-center justify-center">
            <div class="container">
              <p>
                There doesn't appear to be any commits yet for this project.
              </p>
            </div>
          </div>
        }
      >
        <div data-commits="loaded">
          <Show when={failures().length > 0}>
            <div class="container">
              <div class="my-4 flex gap-2 rounded-md bg-red-950 p-2">
                <div class="flex h-8 w-11 items-center justify-center">🚨</div>
                <div class="flex flex-col gap-2 text-sm text-red-300">
                  <h2 class="text-lg font-semibold text-red-100">
                    Pipeline Failure: Immediate Attention Required!
                  </h2>
                  <p>
                    The pipeline is currently blocked due to a failure in the{" "}
                    {failures().join(", ")}. No new feature commits should be
                    pushed until the issue is resolved.
                  </p>
                </div>
              </div>
            </div>
          </Show>

          {groupedCommits().map((section) => (
            <Section
              name={section.name}
              status={section.status}
              medianLeadTime={
                !section.status
                  ? (getMedianLeadTime(
                      section.groups.flatMap((x) => x.commits),
                    ) ?? 0)
                  : undefined
              }
              deploys={
                !section.status
                  ? section.groups
                      .flatMap((x) => x.commits)
                      .filter((x) => x.deployStatus === "passed").length
                  : undefined
              }
            >
              <For each={section.groups}>
                {(group) => (
                  <Group name={group.name}>
                    <For each={group.commits}>
                      {(commit) => (
                        <CommitListItem
                          commit={commit}
                          isCauseOfFailure={Boolean(
                            section.brokenBy &&
                              commit.hash === section.brokenBy,
                          )}
                        />
                      )}
                    </For>
                  </Group>
                )}
              </For>
              <Show when={section.groups.length === 0}>
                <div class="container py-10 text-center text-neutral-500">
                  No {section.name} commits right now
                </div>
              </Show>
            </Section>
          ))}
        </div>
      </Show>
    </Suspense>
  );
}

function Section(props: {
  children: JSX.Element;
  status?: string;
  medianLeadTime?: number;
  deploys?: number;
  name: string;
}) {
  return (
    <section
      class={classNames(
        "py-8",
        props.status === "failing" ? "bg-red-950ds" : "",
      )}
      data-section={props.name}
      data-status={props.status}
    >
      <div class="sticky top-11 z-10 flex justify-center py-1">
        <div
          class={classNames(
            "flex items-center justify-center gap-4 rounded-md px-3 py-1",

            props.status === "failing"
              ? "bg-red-900 text-red-200"
              : "bg-neutral-900 text-neutral-200",
          )}
        >
          <div>
            {props.name?.[0].toUpperCase()}
            {props.name.slice(1).replaceAll("-", " ")}
            {props.status !== undefined && (
              <>
                &nbsp;&nbsp;
                <span class="text-xs uppercase opacity-50">{props.status}</span>
              </>
            )}
            {props.deploys !== undefined && (
              <>
                &nbsp;&nbsp;
                <span class="text-xs text-cornflower-blue-500 opacity-50">
                  {props.deploys} deploys
                </span>
              </>
            )}
            {props.medianLeadTime !== undefined && (
              <>
                &nbsp;&nbsp;
                <span class="text-xs text-cyan-500 opacity-50">
                  {getShortDurationFormatFromSeconds(props.medianLeadTime)}{" "}
                  median lead time
                </span>
              </>
            )}
          </div>
        </div>
      </div>

      {props.children}
    </section>
  );
}

function Group(props: { children: JSX.Element; name: string }) {
  return (
    <div class="font-mono">
      <h3 class="container text-sm">
        <br />
        {props.name}
      </h3>
      <ul>{props.children}</ul>
    </div>
  );
}
