export type status = "passed" | "failed" | "started" | "none";

export type pipelineSection = "commit-stage" | "acceptance-stage" | "deploy";

export interface Commit {
  // Define the structure of a commit here
  hash: string;
  message: string;
  authorName: string;
  authorEmail: string;
  date: string;

  commitStageStatus: status;
  commitStageStartedAt: string;
  commitStageCompletedAt: string;

  acceptanceStageStatus: status;
  acceptanceStageStartedAt: string;
  acceptanceStageCompletedAt: string;

  deployStatus: status;
  deployStartedAt: string;
  deployCompletedAt: string;

  leadTimeCompletedAt: string;
}
