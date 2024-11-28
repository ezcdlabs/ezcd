import { test, expect, describe } from "vitest";
// import { render } from "@solidjs/testing-library"
// import userEvent from "@testing-library/user-event"
import groupCommits, { CommitForGrouping } from "./groupCommits";

function makeCommit(commit: Partial<CommitForGrouping>): CommitForGrouping {
  return {
    hash: "default",
    commitStageStatus: "none",
    acceptanceStageStatus: "none",
    deployStatus: "none",
    ...commit,
  };
}

const defaultCommit: CommitForGrouping = {
  hash: "default",
  commitStageStatus: "none",
  acceptanceStageStatus: "none",
  deployStatus: "none",
};

test("should return correct 3 sections by default", async () => {
  const actual = groupCommits([]);

  expect(actual).toEqual([
    { name: "commit-stage", status: "ok", groups: [] },
    { name: "acceptance-stage", status: "ok", groups: [] },
    { name: "deploy", status: "ok", groups: [] },
  ]);
});

test("should put commit in commit-stage", async () => {
  const commit1 = makeCommit({
    ...defaultCommit,
    hash: "1",
    commitStageStatus: "started",
  });

  const actual = groupCommits([commit1]);

  expect(actual).toEqual([
    {
      name: "commit-stage",
      status: "ok",
      groups: [
        {
          name: "Running commit stage:",
          commits: [commit1],
        },
      ],
    },
    { name: "acceptance-stage", status: "ok", groups: [] },
    { name: "deploy", status: "ok", groups: [] },
  ]);
});

test("should put commit in acceptance-stage", async () => {
  const commit1 = makeCommit({
    ...defaultCommit,
    hash: "1",
    commitStageStatus: "passed",
    acceptanceStageStatus: "started",
  });

  const actual = groupCommits([commit1]);

  expect(actual).toEqual([
    {
      name: "commit-stage",
      status: "ok",
      groups: [],
    },
    {
      name: "acceptance-stage",
      status: "ok",
      groups: [
        {
          name: "Running acceptance stage:",
          commits: [commit1],
        },
      ],
    },
    { name: "deploy", status: "ok", groups: [] },
  ]);
});

test("should put commit in deploy queue", async () => {
  const commit1 = makeCommit({
    ...defaultCommit,
    hash: "1",
    commitStageStatus: "passed",
    acceptanceStageStatus: "passed",
  });

  const actual = groupCommits([commit1]);

  expect(actual).toEqual([
    {
      name: "commit-stage",
      status: "ok",
      groups: [],
    },
    {
      name: "acceptance-stage",
      status: "ok",
      groups: [],
    },
    {
      name: "deploy",
      status: "ok",
      groups: [
        {
          name: "Queued for deploy:",
          commits: [commit1],
        },
      ],
    },
  ]);
});

test("should put commit in deployed by week", async () => {
  const commit1 = makeCommit({
    ...defaultCommit,
    hash: "1",
    commitStageStatus: "passed",
    acceptanceStageStatus: "passed",
    deployStatus: "passed",
    deployCompletedAt: "2024-11-28T10:00:00Z",
  });

  const actual = groupCommits([commit1]);

  expect(actual).toEqual([
    {
      name: "commit-stage",
      status: "ok",
      groups: [],
    },
    {
      name: "acceptance-stage",
      status: "ok",
      groups: [],
    },
    {
      name: "deploy",
      status: "ok",
      groups: [],
    },
    {
      name: "Deployed in week of Mon, 25 Nov 2024:",
      status: "ok",
      groups: [
        {
          name: "Deployed on Thu, 28 Nov 2024 at 10:00 AM:",
          commits: [commit1],
        },
      ],
    },
  ]);
});

test("should show commit as deploying", async () => {
  const commit1 = makeCommit({
    ...defaultCommit,
    hash: "1",
    commitStageStatus: "passed",
    acceptanceStageStatus: "passed",
    deployStatus: "started",
  });

  const actual = groupCommits([commit1]);

  expect(actual).toEqual([
    {
      name: "commit-stage",
      status: "ok",
      groups: [],
    },
    {
      name: "acceptance-stage",
      status: "ok",
      groups: [],
    },
    {
      name: "deploy",
      status: "ok",
      groups: [
        {
          name: "Deploying:",
          commits: [commit1],
        },
      ],
    },
  ]);
});

test("should show commit as failed deploy", async () => {
  const commit1 = makeCommit({
    ...defaultCommit,
    hash: "1",
    commitStageStatus: "passed",
    acceptanceStageStatus: "passed",
    deployStatus: "failed",
  });

  const actual = groupCommits([commit1]);

  expect(actual).toEqual([
    {
      name: "commit-stage",
      status: "ok",
      groups: [],
    },
    {
      name: "acceptance-stage",
      status: "ok",
      groups: [],
    },
    {
      name: "deploy",
      status: "failing",
      groups: [
        {
          name: "Failed to deploy:",
          commits: [commit1],
        },
      ],
    },
  ]);
});

test("should queue second commit for acceptance-stage", async () => {
  const commit1 = makeCommit({
    ...defaultCommit,
    hash: "1",
    commitStageStatus: "passed",
    acceptanceStageStatus: "started",
  });
  const commit2 = makeCommit({
    ...defaultCommit,
    hash: "2",
    commitStageStatus: "passed",
    acceptanceStageStatus: "none",
  });

  const actual = groupCommits([commit2, commit1]);

  expect(actual).toEqual([
    {
      name: "commit-stage",
      status: "ok",
      groups: [],
    },
    {
      name: "acceptance-stage",
      status: "ok",
      groups: [
        {
          name: "Queued for acceptance stage:",
          commits: [commit2],
        },
        {
          name: "Running acceptance stage:",
          commits: [commit1],
        },
      ],
    },
    { name: "deploy", status: "ok", groups: [] },
  ]);
});

test("should show failed commit stage in failed group and set the commit stage to failing", async () => {
  const commit1 = makeCommit({
    ...defaultCommit,
    hash: "1",
    commitStageStatus: "failed",
  });

  const actual = groupCommits([commit1]);

  expect(actual).toEqual([
    {
      name: "commit-stage",
      status: "failing",
      groups: [
        {
          name: "Failed commit stage:",
          commits: [commit1],
        },
      ],
    },
    {
      name: "acceptance-stage",
      status: "ok",
      groups: [],
    },
    { name: "deploy", status: "ok", groups: [] },
  ]);
});

test("should show fixed commit stage if newer commit fixes previous error", async () => {
  const commit1 = makeCommit({
    ...defaultCommit,
    hash: "1",
    commitStageStatus: "failed",
  });
  const commit2 = makeCommit({
    ...defaultCommit,
    hash: "2",
    commitStageStatus: "passed",
  });

  const actual = groupCommits([commit2, commit1]);

  expect(actual).toEqual([
    {
      name: "commit-stage",
      status: "ok",
      groups: [],
    },
    {
      name: "acceptance-stage",
      status: "ok",
      groups: [
        {
          name: "Queued for acceptance stage:",
          commits: [commit2, commit1],
        },
      ],
    },
    { name: "deploy", status: "ok", groups: [] },
  ]);
});

test("should show failed acceptance stage in failed group and set the acceptance stage to failing", async () => {
  const commit1 = makeCommit({
    ...defaultCommit,
    hash: "1",
    commitStageStatus: "passed",
    acceptanceStageStatus: "failed",
  });

  const actual = groupCommits([commit1]);

  expect(actual).toEqual([
    {
      name: "commit-stage",
      status: "ok",
      groups: [],
    },
    {
      name: "acceptance-stage",
      status: "failing",
      groups: [
        {
          name: "Failed acceptance stage:",
          commits: [commit1],
        },
      ],
    },
    { name: "deploy", status: "ok", groups: [] },
  ]);
});
