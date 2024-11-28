import DSL from "dsl/DSL";
import { test } from "@playwright/test";

test("should move fixed commits to acceptance stage", async ({ page }) => {
  const dsl = new DSL(page);

  await dsl.cli.createProject("project1");

  await dsl.cli.commitStageStarted({
    project: "project1",
    commitMessage: "First commit",
  });
  await dsl.cli.commitStageFailed({
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

  await dsl.ui.checkCommit({
    project: "project1",
    commitMessage: "First commit",
    section: "acceptance-stage",
  });

  await dsl.ui.checkStage({
    project: "project1",
    stage: "commit-stage",
    status: "ok",
  });
});
