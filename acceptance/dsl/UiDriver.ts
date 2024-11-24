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

    async getProjectCommit(params: { projectId: string, commitHash: string }) {
        await this.page.goto(`/project/${params.projectId}`);

        // wait for the commits to load
        await this.page.waitForSelector(`[data-commits=loaded]`);

        // find the element with the data tag [commit-hash="commitHash"]
        const commit = await this.page.$(`[data-commit="${params.commitHash}"]`);

        if (!commit) {
            throw new Error(`Commit with hash ${params.commitHash} not found`);
        }

        const elements = {
            commitMessage: await commit.$("[data-label=commitMessage]"),
            commitAuthor: await commit.$("[data-label=commitAuthor]"),
            commitStageStatus: await commit.$("[data-label=commitStageStatus]"),
            acceptanceStageStatus: await commit.$("[data-label=acceptanceStageStatus]"),
            deployStatus: await commit.$("[data-label=deployStatus]"),
        }

        return {
            commitMessage: await elements.commitMessage?.innerText(),
            commitAuthor: await elements.commitAuthor?.innerText(),
            commitStageStatus: await elements.commitStageStatus?.getAttribute("data-value"),
            acceptanceStageStatus: await elements.acceptanceStageStatus?.getAttribute("data-value"),
            deployStatus: await elements.deployStatus?.getAttribute("data-value"),
        }
    }
}