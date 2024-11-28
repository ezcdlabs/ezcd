import { Page } from "@playwright/test";
import { pipelineSection } from "./DSL";

export default class UiDriver {
  private page: Page;

  constructor(page: Page) {
    this.page = page;
  }

  async getProject(projectId: string) {
    // Just so that the home page shows up in the trace.
    await this.page.goto("/");

    await this.page.goto(`/project/${projectId}`);

    await this.page.waitForSelector("[data-label='projectName']");

    return await this.page.innerText("[data-label='projectName']");
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

  async getProjectStageInfo(params: { projectId: string; stage: string }) {
    await this.page.goto(`/project/${params.projectId}`);

    // wait for the commits to load
    await this.page.waitForSelector(`[data-commits=loaded]`);

    // find the element with the data tag [data-pipelineStage="stage"]
    const stageElement = await this.page.$(`[data-section="${params.stage}"]`);

    return {
      status: await stageElement.getAttribute(`data-status`),
    };
  }

  async getProjectCommit(params: { projectId: string; commitHash: string }) {
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
      acceptanceStageStatus: await commit.$(
        "[data-label=acceptanceStageStatus]"
      ),
      deployStatus: await commit.$("[data-label=deployStatus]"),

      leadTime: await commit.$("[data-label=leadTime]"),

      // figure out which section the commit is within:
      section: await this.page
        .locator("section")
        .filter({
          has: this.page.locator(`[data-commit="${params.commitHash}"]`),
        })
        .first()
        .elementHandle(),
    };

    return {
      commitMessage: await elements.commitMessage?.innerText(),
      commitAuthor: await elements.commitAuthor?.innerText(),
      commitStageStatus: await elements.commitStageStatus?.getAttribute(
        "data-value"
      ),
      acceptanceStageStatus: await elements.acceptanceStageStatus?.getAttribute(
        "data-value"
      ),
      deployStatus: await elements.deployStatus?.getAttribute("data-value"),

      isLeadTimeStopped:
        (await elements.leadTime?.getAttribute("data-stopped")) === "true",

      section: (await elements.section?.getAttribute(
        "data-section"
      )) as pipelineSection,
    };
  }
}
