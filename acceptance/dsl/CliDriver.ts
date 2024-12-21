import { promisify } from "util";
import { exec as execCallback } from "child_process";
import { test } from "@playwright/test";
const exec = promisify(execCallback);

async function ezcdCli(args: string) {
  return test.step(`ezcd-cli ${args}`, async () => {
    try {
      const { stdout, stderr } = await exec(`../dist/ezcd-cli ${args}`, {
        env: {
          EZCD_DATABASE_URL: process.env.EZCD_DATABASE_URL,
        },
      });
      console.log(`stdout: ${stdout}`);
      console.error(`stderr: ${stderr}`);

      // record the result in playwright report:
      test.info().attach(`ezcd-cli ${args}`, {
        body: `stdout: ${stdout}\n\nstderr: ${stderr}`,
      });

      return stdout.trimEnd();
    } catch (error) {
      throw new Error(
        `Failed to run ezcd-cli. EZCD_DATABASE_URL was '${process.env.EZCD_DATABASE_URL}' error:${error} stdout:${error.stdout} stderr:${error.stderr}`
      );
    }
  });
}

export default class CLIDriver {
  getVersion = async () => {
    return await ezcdCli("--version");
  };

  createProject = async (project: string) => {
    return await ezcdCli(`create-project ${project}`);
  };

  commitStageStarted = async (action: {
    projectId: string;
    commitHash: string;
    commitMessage: string;
    commitAuthorName: string;
    commitAuthorEmail: string;
    commitDate: Date;
  }) => {
    return await ezcdCli(
      `commit-stage-started --project ${action.projectId} --hash "${
        action.commitHash
      }" --message "${action.commitMessage}" --author-name "${
        action.commitAuthorName
      }" --author-email "${
        action.commitAuthorEmail
      }" --date "${action.commitDate.toISOString()}"`
    );
  };

  commitStagePassed = async (action: {
    projectId: string;
    commitHash: string;
  }) => {
    return await ezcdCli(
      `commit-stage-passed --project ${action.projectId} --hash "${action.commitHash}"`
    );
  };

  commitStageFailed = async (action: {
    projectId: string;
    commitHash: string;
  }) => {
    return await ezcdCli(
      `commit-stage-failed --project ${action.projectId} --hash "${action.commitHash}"`
    );
  };

  acceptanceStageStarted = async (action: {
    projectId: string;
    commitHash: string;
  }) => {
    return await ezcdCli(
      `acceptance-stage-started --project ${action.projectId} --hash "${action.commitHash}"`
    );
  };
  acceptanceStagePassed = async (action: {
    projectId: string;
    commitHash: string;
  }) => {
    return await ezcdCli(
      `acceptance-stage-passed --project ${action.projectId} --hash "${action.commitHash}"`
    );
  };
  acceptanceStageFailed = async (action: {
    projectId: string;
    commitHash: string;
  }) => {
    return await ezcdCli(
      `acceptance-stage-failed --project ${action.projectId} --hash "${action.commitHash}"`
    );
  };
  deployStarted = async (action: { projectId: string; commitHash: string }) => {
    return await ezcdCli(
      `deploy-started --project ${action.projectId} --hash "${action.commitHash}"`
    );
  };
  deployPassed = async (action: { projectId: string; commitHash: string }) => {
    return await ezcdCli(
      `deploy-passed --project ${action.projectId} --hash "${action.commitHash}"`
    );
  };
  deployFailed = async (action: { projectId: string; commitHash: string }) => {
    return await ezcdCli(
      `deploy-failed --project ${action.projectId} --hash "${action.commitHash}"`
    );
  };

  getQueuedForAcceptance = async (action: { projectId: string }) => {
    return await ezcdCli(
      `get-queued-for-acceptance --project ${action.projectId}`
    );
  };
}
