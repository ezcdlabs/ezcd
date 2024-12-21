import DSL from "dsl/DSL";
import { test } from "@playwright/test";

test("should show the commit that is ready for acceptance in the cli", async ({
  page,
}) => {
  const dsl = new DSL(page);

  await dsl.cli.createProject("project1");

  await dsl.cli.commitStageStarted({
    project: "project1",
    commitMessage: "First commit",
  });

  await dsl.cli.commitStageStarted({
    project: "project1",
    commitMessage: "Second commit",
  });

  await dsl.cli.commitStagePassed({
    project: "project1",
    commitMessage: "First commit",
  });

  await dsl.cli.checkQueuedForAcceptance({
    project: "project1",
    commitMessage: "First commit",
  });
});
