import { promisify } from 'util';
import { exec as execCallback } from 'child_process';
import { test } from '@playwright/test'
const exec = promisify(execCallback);

async function ezcdCli(args: string) {
  return test.step(`ezcd-cli ${args}`, async () => {
    try {
      const { stdout, stderr } = await exec(`../dist/ezcd-cli ${args}`, {
        env: {
          EZCD_DATABASE_URL: process.env.EZCD_DATABASE_URL,
        }
      });
      console.log(`stdout: ${stdout}`);
      console.error(`stderr: ${stderr}`);

      // record the result in playwright report:
      test.info().attach(`ezcd-cli ${args}`, {
        body: `stdout: ${stdout}\n\nstderr: ${stderr}`,
      });

      return stdout;
    } catch (error) {
      throw new Error(`Failed to run ezcd-cli. EZCD_DATABASE_URL was '${process.env.EZCD_DATABASE_URL}' error:${error} stdout:${error.stdout} stderr:${error.stderr}`);
    }
  });
}

export default class CLIDriver {

  getVersion = async () => {
    return await ezcdCli('--version');
  }

  createProject = async (project: string) => {
    return await ezcdCli(`create-project ${project}`);
  }

  commitPhaseStarted = async (projectId: string, commitMessage: string) => {
    return await ezcdCli(`commit ${projectId} "${commitMessage}"`);
  }

}