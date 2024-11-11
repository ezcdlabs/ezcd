import { test } from '@playwright/test';
import DSL from './dsl/DSL';

test('should create commits', async ({ page }) => {

  const dsl = new DSL(page);

  await dsl.cli.getVersion();

  await dsl.cli.createProject('project1');
  await dsl.cli.commitPhaseStarted('project1', 'commit1');
  await dsl.cli.commitPhaseStarted('project1', 'commit2');
  await dsl.cli.commitPhaseStarted('project1', 'commit3');
  
  await dsl.ui.verifyProjectCommits('project1', ['commit1', 'commit2', 'commit3']);
});