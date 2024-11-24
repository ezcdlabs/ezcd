import DSL from "dsl/DSL";
import { test } from '@playwright/test'

test('should show commit as started the acceptance stage', async ({ page }) => {
    const dsl = new DSL(page);

    await dsl.cli.createProject('project1');

    await dsl.cli.commitStageStarted({ project: 'project1', commitMessage: 'First commit' });
    await dsl.cli.commitStagePassed({ project: 'project1', commitMessage: 'First commit' });

    await dsl.cli.acceptanceStageStarted({ project: 'project1', commitMessage: 'First commit' });

    await dsl.ui.checkCommit({ project: 'project1', commitMessage: 'First commit', commitStage: 'passed', acceptanceStage: 'started' });
});

test('should show commit as passed the acceptance stage', async ({ page }) => {
    const dsl = new DSL(page);

    await dsl.cli.createProject('project1');

    await dsl.cli.commitStageStarted({ project: 'project1', commitMessage: 'First commit' });
    await dsl.cli.commitStagePassed({ project: 'project1', commitMessage: 'First commit' });

    await dsl.cli.acceptanceStageStarted({ project: 'project1', commitMessage: 'First commit' });
    await dsl.cli.acceptanceStagePassed({ project: 'project1', commitMessage: 'First commit' });

    await dsl.ui.checkCommit({ project: 'project1', commitMessage: 'First commit', commitStage: 'passed', acceptanceStage: 'passed' });
});

test('should show commit as failed the acceptance stage', async ({ page }) => {
    const dsl = new DSL(page);

    await dsl.cli.createProject('project1');

    await dsl.cli.commitStageStarted({ project: 'project1', commitMessage: 'First commit' });
    await dsl.cli.commitStagePassed({ project: 'project1', commitMessage: 'First commit' });

    await dsl.cli.acceptanceStageStarted({ project: 'project1', commitMessage: 'First commit' });
    await dsl.cli.acceptanceStageFailed({ project: 'project1', commitMessage: 'First commit' });

    await dsl.ui.checkCommit({ project: 'project1', commitMessage: 'First commit', commitStage: 'passed', acceptanceStage: 'failed' });
});