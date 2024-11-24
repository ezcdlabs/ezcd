import DSL from "dsl/DSL";
import { test } from '@playwright/test'

test('should show commit as started the deploy', async ({ page }) => {
    const dsl = new DSL(page);

    await dsl.cli.createProject('project1');

    await dsl.cli.commitStageStarted({ project: 'project1', commitMessage: 'First commit' });
    await dsl.cli.commitStagePassed({ project: 'project1', commitMessage: 'First commit' });
    await dsl.cli.acceptanceStageStarted({ project: 'project1', commitMessage: 'First commit' });
    await dsl.cli.acceptanceStagePassed({ project: 'project1', commitMessage: 'First commit' });

    await dsl.cli.deployStarted({ project: 'project1', commitMessage: 'First commit' });

    await dsl.ui.checkCommit({ project: 'project1', commitMessage: 'First commit', commitStage: 'passed', acceptanceStage: 'passed', deploy: 'started' });
});