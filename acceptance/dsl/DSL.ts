import { expect, Page } from '@playwright/test';
import UiDriver from './UiDriver';
import * as uuid from 'uuid'
import CLIDriver from './CliDriver';
import * as dateFns from 'date-fns';

export default class DSL {
    private projects = new Map<string, string>();
    private commits = new Map<string, string>();

    private getOrThrow(map: Map<string, string>, key: string) {
        const value = map.get(key);
        if (!value) {
            throw new Error(`Key ${key} not found`);
        }
        return value;
    }

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
            const projectId = randomSuffix(project)
            await this.cliDriver.createProject(projectId);
            this.projects.set(project, projectId);
        },

        commitStageStarted: async (params: { project: string, commitMessage: string, commitAuthorName?: string, commitDate?: string }) => {
            const projectId = this.getOrThrow(this.projects, params.project)
            const commitHash = uuid.v4();
            const commitAuthorName = params.commitAuthorName ?? "John Doe";
            const commitAuthorEmail = params.commitAuthorName ?? "jd@ezcd.com";
            const commitDate = params.commitDate ? dateFns.parse(params.commitDate, "yyyy-MM-dd HH:mm:ss", new Date()) : new Date();

            this.commits.set(params.commitMessage, commitHash);

            await this.cliDriver.commitStageStarted({
                projectId,
                commitHash,
                commitMessage: params.commitMessage,
                commitAuthorName,
                commitAuthorEmail,
                commitDate
            });
        },

        commitStagePassed: async (params: { project: string, commitMessage: string }) => {
            const projectId = this.getOrThrow(this.projects, params.project)
            const commitHash = this.getOrThrow(this.commits, params.commitMessage);

            await this.cliDriver.commitStagePassed({
                projectId,
                commitHash,
            });
        },

        commitStageFailed: async (params: { project: string, commitMessage: string }) => {
            const projectId = this.getOrThrow(this.projects, params.project)
            const commitHash = this.getOrThrow(this.commits, params.commitMessage);

            await this.cliDriver.commitStageFailed({
                projectId,
                commitHash,
            });
        },

        acceptanceStageStarted: async (params: { project: string, commitMessage: string }) => {
            const projectId = this.getOrThrow(this.projects, params.project)
            const commitHash = this.getOrThrow(this.commits, params.commitMessage);

            await this.cliDriver.acceptanceStageStarted({
                projectId,
                commitHash,
            });
        },

        acceptanceStagePassed: async (params: { project: string, commitMessage: string }) => {
            const projectId = this.getOrThrow(this.projects, params.project)
            const commitHash = this.getOrThrow(this.commits, params.commitMessage);

            await this.cliDriver.acceptanceStagePassed({
                projectId,
                commitHash,
            });
        },

        acceptanceStageFailed: async (params: { project: string, commitMessage: string }) => {
            const projectId = this.getOrThrow(this.projects, params.project)
            const commitHash = this.getOrThrow(this.commits, params.commitMessage);

            await this.cliDriver.acceptanceStageFailed({
                projectId,
                commitHash,
            });
        },

        deployStarted: async (params: { project: string, commitMessage: string }) => {
            const projectId = this.getOrThrow(this.projects, params.project)
            const commitHash = this.getOrThrow(this.commits, params.commitMessage);

            await this.cliDriver.deployStarted({
                projectId,
                commitHash,
            });
        },

        deployPassed: async (params: { project: string, commitMessage: string }) => {
            const projectId = this.getOrThrow(this.projects, params.project)
            const commitHash = this.getOrThrow(this.commits, params.commitMessage);

            await this.cliDriver.deployPassed({
                projectId,
                commitHash,
            });
        },

        deployFailed: async (params: { project: string, commitMessage: string }) => {
            const projectId = this.getOrThrow(this.projects, params.project)
            const commitHash = this.getOrThrow(this.commits, params.commitMessage);

            await this.cliDriver.deployFailed({
                projectId,
                commitHash,
            });
        }

    }

    ui = {
        verifyProject: async (project: string) => {
            const projectId = this.projects.get(project);
            if (!projectId) {
                throw new Error(`Project ${project} not found`);
            }
            const actual = await this.uiDriver.getProject(projectId);

            expect(actual).toEqual(projectId);
        },

        checkCommit: async (params: { project: string, commitMessage: string, commitStage?: string, acceptanceStage?: string, deploy?: string }) => {
            const projectId = this.getOrThrow(this.projects, params.project)
            const commitHash = this.getOrThrow(this.commits, params.commitMessage);

            const commit = await this.uiDriver.getProjectCommit({ projectId, commitHash });

            if (params.commitStage) {
                expect(commit.commitStageStatus).toEqual(params.commitStage);
            }
            if (params.acceptanceStage) {
                expect(commit.acceptanceStageStatus).toEqual(params.acceptanceStage);
            }
            if (params.deploy) {
                expect(commit.deployStatus).toEqual(params.deploy);
            }
        }
    }
}

function randomSuffix(input: string) {
    return input + "_" + Math.random().toString(36).substring(7);
}