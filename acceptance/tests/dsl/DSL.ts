import { expect, Page } from '@playwright/test';
import UiDriver from './UiDriver';
import * as uuid from 'uuid'
import CLIDriver from './CliDriver';

export default class DSL {
    private projects = new Map<string, string>();
    private uiDriver: UiDriver;
    private cliDriver: CLIDriver;

    constructor(page: Page) {
        this.uiDriver = new UiDriver(page);
        this.cliDriver = new CLIDriver();
    }

    cli = {
        getVersion: async () => {
            return this.cliDriver.getVersion();
        },

        createProject: async (project: string) => {
            const projectId = await this.cliDriver.createProject(project);
            this.projects.set(project, projectId);
        },

        commitPhaseStarted: async (project: string, commitMessage: string) => {
            var projectId = this.projects.get(project);
            if (!projectId) {
                throw new Error(`Project ${project} not found`);
            }

            // TODO, actuall call the CLI and create a commit
            // await this.cliDriver.commitPhaseStarted(projectId, commitMessage);
        }
    }

    ui = {
        verifyProject: async (project: string) => {
            var projectId = this.projects.get(project);
            if (!projectId) {
                throw new Error(`Project ${project} not found`);
            }
            const actual = await this.uiDriver.getProject(projectId);

            expect(actual).toEqual(project);
        },

        verifyProjectCommits: async (project: string, expected: string[]) => {
            var projectId = this.projects.get(project);
            if (!projectId) {
                throw new Error(`Project ${project} not found`);
            }
            const actual = await this.uiDriver.getProjectCommits(projectId);

            expect(actual).toEqual(expected);
        }
    }
}

function randomSuffix(input: string) {
    return input + "_" + Math.random().toString(36).substring(7);
}