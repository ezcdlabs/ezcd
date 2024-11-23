import { test } from '@playwright/test';
import DSL from '../dsl/DSL';

test('should create new project', async ({ page }) => {

  const dsl = new DSL(page);

  await dsl.cli.createProject('project1');
  await dsl.ui.verifyProject('project1');
});