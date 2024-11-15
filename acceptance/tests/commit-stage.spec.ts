import { test } from '@playwright/test'
import DSL from './dsl/DSL'

test('should add commits and start the commit stage', async ({ page }) => {
    const dsl = new DSL(page);
    
    await dsl.cli.createProject('project1');
    await dsl.cli.commitStageStarted({ project: 'project1', commitMessage: 'First commit' });
    await dsl.cli.commitStageStarted({ project: 'project1', commitMessage: 'Second commit' });
    
    await dsl.ui.checkCommit({ project: 'project1', commitMessage: 'First commit', commitStage: 'STARTED' });
    await dsl.ui.checkCommit({ project: 'project1', commitMessage: 'Second commit', commitStage: 'STARTED' });
});