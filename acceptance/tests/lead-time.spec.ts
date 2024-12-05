import DSL from "dsl/DSL";
import { test } from "@playwright/test";

test("should not stop lead time for commits that are not deployed", async ({
  page,
}) => {
  const dsl = new DSL(page);

  await dsl.cli.createProject("project1");

  await dsl.cli.commitStageStarted({
    project: "project1",
    commitMessage: "First commit",
  });

  await dsl.ui.checkCommit({
    project: "project1",
    commitMessage: "First commit",
    isLeadTimeStopped: false,
  });
});

test("should stop lead time for commits that are deployed", async ({
  page,
}) => {
  const dsl = new DSL(page);

  await dsl.cli.createProject("project1");

  await dsl.cli.commitStageStarted({
    project: "project1",
    commitMessage: "First commit",
  });
  await dsl.cli.commitStagePassed({
    project: "project1",
    commitMessage: "First commit",
  });
  await dsl.cli.acceptanceStageStarted({
    project: "project1",
    commitMessage: "First commit",
  });
  await dsl.cli.acceptanceStagePassed({
    project: "project1",
    commitMessage: "First commit",
  });
  await dsl.cli.deployStarted({
    project: "project1",
    commitMessage: "First commit",
  });
  await dsl.cli.deployPassed({
    project: "project1",
    commitMessage: "First commit",
  });

  await dsl.ui.checkCommit({
    project: "project1",
    commitMessage: "First commit",
    isLeadTimeStopped: true,
  });
});

test("should stop lead time for commits that are included in another later commit that is deployed", async ({
  page,
}) => {
  const dsl = new DSL(page);

  await dsl.cli.createProject("project1");

  await dsl.cli.commitStageStarted({
    project: "project1",
    commitMessage: "First commit",
  });
  await dsl.cli.commitStagePassed({
    project: "project1",
    commitMessage: "First commit",
  });

  await dsl.cli.commitStageStarted({
    project: "project1",
    commitMessage: "Second commit",
  });
  await dsl.cli.commitStagePassed({
    project: "project1",
    commitMessage: "Second commit",
  });
  await dsl.cli.acceptanceStageStarted({
    project: "project1",
    commitMessage: "Second commit",
  });
  await dsl.cli.acceptanceStagePassed({
    project: "project1",
    commitMessage: "Second commit",
  });
  await dsl.cli.deployStarted({
    project: "project1",
    commitMessage: "Second commit",
  });
  await dsl.cli.deployPassed({
    project: "project1",
    commitMessage: "Second commit",
  });

  await dsl.ui.checkCommit({
    project: "project1",
    commitMessage: "First commit",
    isLeadTimeStopped: true,
  });
  await dsl.ui.checkCommit({
    project: "project1",
    commitMessage: "Second commit",
    isLeadTimeStopped: true,
  });
});
