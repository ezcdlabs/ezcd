import DSL from "dsl/DSL";
import { test } from "@playwright/test";

async function group(dsl: DSL, project: string, size: number, deploy: boolean) {
  const randomGroup = Math.random().toString(16).slice(2, 7);

  for (let i = 1; i <= size; i++) {
    await dsl.cli.commitStageStarted({
      project: project,
      commitMessage: `Commit ${randomGroup} ${i}`,
    });
    await dsl.cli.commitStagePassed({
      project: project,
      commitMessage: `Commit ${randomGroup} ${i}`,
    });
  }

  await dsl.cli.acceptanceStageStarted({
    project: project,
    commitMessage: `Commit ${randomGroup} ${size}`,
  });
  await dsl.cli.acceptanceStagePassed({
    project: project,
    commitMessage: `Commit ${randomGroup} ${size}`,
  });

  if (deploy) {
    await dsl.cli.deployStarted({
      project: project,
      commitMessage: `Commit ${randomGroup} ${size}`,
    });
    await dsl.cli.deployPassed({
      project: project,
      commitMessage: `Commit ${randomGroup} ${size}`,
    });
  }
}

test("should display correctly", async ({ page }) => {
  const dsl = new DSL(page);

  await dsl.cli.createProject("project1");

  const groups = [13, 7, 3, 3, 1, 8, 12, 7, 3, 5, 8, 9];
  //const groups = [3, 4, 5];

  for (const size of groups) {
    await group(dsl, "project1", size, true);
  }

  await group(dsl, "project1", 3, false);
});
