export interface CommitForGrouping {
  commitStageStatus: "none" | "started" | "passed" | "failed"
  acceptanceStageStatus: "none" | "started" | "passed" | "failed"
  deployStatus: "none" | "started" | "passed" | "failed"
}

export default function groupCommits(commits: CommitForGrouping[]) {
  return null;
}
