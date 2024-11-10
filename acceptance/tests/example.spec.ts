import { test, expect } from '@playwright/test';
import { promisify } from 'util';
import { exec as execCallback } from 'child_process';

const exec = promisify(execCallback);
async function ezcdCli(args: string) {
  try {
    const { stdout, stderr } = await exec(`../dist/ezcd-cli ${args}`);
    console.log(`stdout: ${stdout}`);
    console.error(`stderr: ${stderr}`);
    return stdout;
  } catch (error) {
    console.error(`exec error: ${error}`);
    throw error;
  }
}

// test('has title', async ({ page }) => {
//   await page.goto('https://playwright.dev/');

//   // Expect a title "to contain" a substring.
//   await expect(page).toHaveTitle(/Playwright/);
// });

// test('get started link', async ({ page }) => {
//   await page.goto('https://playwright.dev/');

//   // Click the get started link.
//   await page.getByRole('link', { name: 'Get started' }).click();

//   // Expects page to have a heading with the name of Installation.
//   await expect(page.getByRole('heading', { name: 'Installation' })).toBeVisible();
// });

test('should show the API result', async ({ page }) => {
  const result = await ezcdCli('--version')

  await expect(result).toBeTruthy();

  await page.goto("http://localhost:3000/");

  await expect(page.getByText('Hello, API!!!')).toBeVisible();
});