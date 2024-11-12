import { Page } from "@playwright/test";

export default class UiDriver {
    private page: Page;

    constructor(page: Page) {
        this.page = page;
    }

    async getProject(projectId: string) {
        // Just so that the home page shows up in the trace.
        await this.page.goto('/');

        await this.page.goto(`/project/${projectId}`);

        await this.page.waitForSelector(".project-name");

        return await this.page.innerText(".project-name");
    }

    async getProjectCommits(projectId: string) {
        await this.page.goto(`/project/${projectId}`);

        // wait for the commits to load
        await this.page.waitForSelector(".commit");

        const commitElements = await this.page.$$(".commit");
        const commitTexts: string[] = [];

        for (const commit of commitElements) {
            const text = await commit.innerText();
            commitTexts.push(text);
        }

        return commitTexts;
    }
}