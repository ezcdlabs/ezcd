import DSL from "dsl/DSL";
import { test } from "@playwright/test";

test("should display correctly", async ({ page }) => {
  const dsl = new DSL(page);

  await dsl.cli.createProject("project1");

  await dsl.cli.commitStageStarted({
    project: "project1",
    commitMessage: "Commit 1",
  });
  await dsl.cli.commitStageStarted({
    project: "project1",
    commitMessage: "Commit 2",
  });
  await dsl.cli.commitStageStarted({
    project: "project1",
    commitMessage: "Commit 3",
  });
  await dsl.cli.commitStageStarted({
    project: "project1",
    commitMessage: "Commit 4",
  });
  await dsl.cli.commitStageStarted({
    project: "project1",
    commitMessage: "Commit 5",
  });

  await dsl.cli.commitStageFailed({
    project: "project1",
    commitMessage: "Commit 1",
  });
  await dsl.cli.commitStagePassed({
    project: "project1",
    commitMessage: "Commit 2",
  });
  await dsl.cli.commitStagePassed({
    project: "project1",
    commitMessage: "Commit 3",
  });
  await dsl.cli.commitStagePassed({
    project: "project1",
    commitMessage: "Commit 4",
  });

  await dsl.cli.acceptanceStageStarted({
    project: "project1",
    commitMessage: "Commit 2",
  });
});
