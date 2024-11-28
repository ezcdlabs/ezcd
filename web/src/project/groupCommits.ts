import * as dateFns from "date-fns";

export interface CommitForGrouping {
  hash: string;
  commitStageStatus: "none" | "started" | "passed" | "failed";
  acceptanceStageStatus: "none" | "started" | "passed" | "failed";
  deployStatus: "none" | "started" | "passed" | "failed";
  deployCompletedAt?: string;
}

export interface Section<T extends CommitForGrouping> {
  name: string;
  status: "ok" | "failing";
  groups: {
    name: string;
    commits: T[];
  }[];
}

export default function groupCommits<T extends CommitForGrouping>(
  commits: T[],
): Section<T>[] {
  const commitStage: Section<T> = {
    name: "commit-stage",
    status: "ok",
    groups: [],
  };

  const acceptanceStage: Section<T> = {
    name: "acceptance-stage",
    status: "ok",
    groups: [],
  };

  const deploy: Section<T> = {
    name: "deploy",
    status: "ok",
    groups: [],
  };

  let index = 0;

  const runningCommitStage: T[] = [];

  // find any that are still running the commit stage
  for (; index < commits.length; index++) {
    if (commits[index].commitStageStatus !== "started") {
      break;
    }

    runningCommitStage.push(commits[index]);
  }

  if (runningCommitStage.length > 0) {
    commitStage.groups.push({
      name: "Running commit stage:",
      commits: runningCommitStage,
    });
  }

  const failingCommitStage: T[] = [];
  // find any that are still queuing for the acceptance stage
  for (; index < commits.length; index++) {
    if (commits[index].commitStageStatus !== "failed") {
      break;
    }

    failingCommitStage.push(commits[index]);
  }

  if (failingCommitStage.length > 0) {
    commitStage.status = "failing";
    commitStage.groups.push({
      name: "Failed commit stage:",
      commits: failingCommitStage,
    });
  }

  const queuedForAcceptanceStage: T[] = [];
  // find any that are still queuing for the acceptance stage
  for (; index < commits.length; index++) {
    if (commits[index].acceptanceStageStatus !== "none") {
      break;
    }

    queuedForAcceptanceStage.push(commits[index]);
  }

  if (queuedForAcceptanceStage.length > 0) {
    acceptanceStage.groups.push({
      name: "Queued for acceptance stage:",
      commits: queuedForAcceptanceStage,
    });
  }

  // find any that are still running the acceptance stage
  const runningAcceptanceStage: T[] = [];

  for (; index < commits.length; index++) {
    if (
      commits[index].acceptanceStageStatus === "passed" ||
      commits[index].acceptanceStageStatus === "failed"
    ) {
      break;
    }

    runningAcceptanceStage.push(commits[index]);
  }

  if (runningAcceptanceStage.length > 0) {
    acceptanceStage.groups.push({
      name: "Running acceptance stage:",
      commits: runningAcceptanceStage,
    });
  }

  // find any that failed the acceptance stage
  const failedAcceptanceStage: T[] = [];

  for (; index < commits.length; index++) {
    if (commits[index].acceptanceStageStatus === "passed") {
      break;
    }

    failedAcceptanceStage.push(commits[index]);
  }

  if (failedAcceptanceStage.length > 0) {
    acceptanceStage.status = "failing";
    acceptanceStage.groups.push({
      name: "Failed acceptance stage:",
      commits: failedAcceptanceStage,
    });
  }

  // find any that are queued for deploy
  const deployQueue: T[] = [];

  for (; index < commits.length; index++) {
    if (commits[index].deployStatus !== "none") {
      break;
    }

    deployQueue.push(commits[index]);
  }

  if (deployQueue.length > 0) {
    deploy.groups.push({
      name: "Queued for deploy:",
      commits: deployQueue,
    });
  }

  const deploying: T[] = [];

  for (; index < commits.length; index++) {
    if (
      commits[index].deployStatus === "passed" ||
      commits[index].deployStatus === "failed"
    ) {
      break;
    }

    deploying.push(commits[index]);
  }

  if (deploying.length > 0) {
    deploy.groups.push({
      name: "Deploying:",
      commits: deploying,
    });
  }

  const failedToDeploy: T[] = [];

  for (; index < commits.length; index++) {
    if (commits[index].deployStatus === "passed") {
      break;
    }

    failedToDeploy.push(commits[index]);
  }

  if (failedToDeploy.length > 0) {
    deploy.status = "failing";
    deploy.groups.push({
      name: "Failed to deploy:",
      commits: failedToDeploy,
    });
  }

  const result = [commitStage, acceptanceStage, deploy];

  for (; index < commits.length; index++) {
    const week: Section<T> = {
      name: `Deployed in week of ${dateFns.format(
        dateFns.startOfWeek(commits[index].deployCompletedAt!, {
          weekStartsOn: 1,
        }),
        "EEE, dd MMM yyyy",
      )}:`,
      status: "ok",
      groups: [],
    };

    for (; index < commits.length; index++) {
      if (!commits[index].deployCompletedAt) {
        // throw an error?
        return result;
      }

      week.groups.push({
        name: `Deployed on ${dateFns.format(commits[index].deployCompletedAt!, "EEE, dd MMM yyyy")} at ${dateFns.format(commits[index].deployCompletedAt!, "hh:mm a")}:`,
        commits: [commits[index]],
      });
    }

    result.push(week);
  }

  return result;
}
