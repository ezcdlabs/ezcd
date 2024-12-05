import classNames from "../utils/classNames";
import TimerDisplay from "./TimerDisplay";
import { Commit, status } from "./types";
import * as dateFns from "date-fns";

export default function CommitListItem(props: {
  commit: Commit;
  isCauseOfFailure: boolean;
}) {
  return (
    <li
      class={classNames(
        "text-sm",
        props.isCauseOfFailure
          ? "bg-red-950 text-red-500"
          : "text-white-secondary",
      )}
      data-commit={props.commit.hash}
    >
      <div class="container flex">
        <div class="grow">
          &nbsp;&nbsp;
          <span>
            <span
              data-label="commitStageStatus"
              data-value={props.commit.commitStageStatus}
            >
              <StatusIndicator
                status={props.commit.commitStageStatus}
                color="white"
              />
              <span
                class="hidden"
                data-label="commitStageStartedAt"
                data-value={props.commit.commitStageStartedAt}
              ></span>

              <span
                class="hidden"
                data-label="commitStageCompletedAt"
                data-value={props.commit.commitStageCompletedAt}
              ></span>
            </span>
            <span
              data-label="acceptanceStageStatus"
              data-value={props.commit.acceptanceStageStatus}
            >
              <StatusIndicator
                status={props.commit.acceptanceStageStatus}
                color="yellow"
              />
              <span
                class="hidden"
                data-label="acceptanceStageStartedAt"
                data-value={props.commit.acceptanceStageStartedAt}
              ></span>

              <span
                class="hidden"
                data-label="acceptanceStageCompletedAt"
                data-value={props.commit.acceptanceStageCompletedAt}
              ></span>
            </span>
            <span
              data-label="deployStatus"
              data-value={props.commit.deployStatus}
            >
              <StatusIndicator
                status={props.commit.deployStatus}
                color="blue"
              />
              <span
                class="hidden"
                data-label="deployStartedAt"
                data-value={props.commit.deployStartedAt}
              ></span>

              <span
                class="hidden"
                data-label="deployCompletedAt"
                data-value={props.commit.deployCompletedAt}
              ></span>
            </span>
          </span>
          <span data-label="commitAuthorName">
            &nbsp;{props.commit.authorName}
          </span>
          <span
            data-label="commitAuthorEmail"
            data-value={props.commit.authorEmail}
          ></span>
          : <span data-label="commitMessage">{props.commit.message}</span>
          &nbsp;(
          <span data-label="commitHash" data-value={props.commit.hash}>
            {props.commit.hash?.slice(0, 7)}
          </span>
          )<span data-label="commitDate" data-value={props.commit.date}></span>
        </div>
        {props.isCauseOfFailure && (
          <span class="mr-2 rounded-full bg-red-900 px-3 text-red-400">
            BREAKING
          </span>
        )}
        <span
          data-label="leadTime"
          data-stopped={Boolean(props.commit.leadTimeCompletedAt)}
          class={classNames(
            props.commit.leadTimeCompletedAt && "text-cyan-500",
          )}
        >
          <TimerDisplay
            start={dateFns.parseISO(props.commit.date)}
            end={
              props.commit.leadTimeCompletedAt
                ? dateFns.parseISO(props.commit.leadTimeCompletedAt)
                : undefined
            }
          />
        </span>
      </div>
    </li>
  );
}

function StatusIndicator(props: {
  status: status;
  color: "white" | "yellow" | "blue";
}) {
  if (props.status === "failed") {
    return <FailedIndicator />;
  }

  if (props.status === "passed") {
    return <PassedIndicator color={props.color} />;
  }

  if (props.status === "started") {
    return <StartedIndicator color={props.color} />;
  }

  return <span class="inline-block h-3 w-3"></span>;
}

function FailedIndicator() {
  return <span class="inline-block h-3 w-3 text-red-400">x</span>;
}

function PassedIndicator(props: { color: "white" | "yellow" | "blue" }) {
  const className =
    "w-3 h-3 inline-flex items-center " +
    ({
      blue: "fill-cornflower-blue-400",
      yellow: "fill-saffron-400",
      white: "fill-neutral-100",
    }[props.color] ?? "fill-neutral-100");

  return (
    <svg class={className} aria-hidden="true" viewBox="0 0 24 24">
      <path d="M9 16.2 4.8 12l-1.4 1.4L9 19 21 7l-1.4-1.4z"></path>
    </svg>
  );
}

function StartedIndicator(props: { color: "white" | "yellow" | "blue" }) {
  const className =
    "w-3 h-3 inline-flex items-center " +
    ({
      blue: "text-cornflower-blue-400",
      yellow: "text-saffron-400",
      white: "text-neutral-100",
    }[props.color] ?? "text-neutral-100");

  return (
    <span class={`braile-spinner ${className}`} />
    // <svg class={className} viewBox="0 0 24 24">
    //   <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2M7 13.5c-.83 0-1.5-.67-1.5-1.5s.67-1.5 1.5-1.5 1.5.67 1.5 1.5-.67 1.5-1.5 1.5m5 0c-.83 0-1.5-.67-1.5-1.5s.67-1.5 1.5-1.5 1.5.67 1.5 1.5-.67 1.5-1.5 1.5m5 0c-.83 0-1.5-.67-1.5-1.5s.67-1.5 1.5-1.5 1.5.67 1.5 1.5-.67 1.5-1.5 1.5"></path>
    // </svg>
  );
}
