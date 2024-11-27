import { test, expect } from "vitest"
// import { render } from "@solidjs/testing-library"
// import userEvent from "@testing-library/user-event"
import groupCommits from "./groupCommits"



const exampleCommit:() => Commit = () => ({
  hash: "123",
  message: "Example commit",
  date: new Date().toString(),
  commitStageStatus: "none",
  acceptanceStageStatus: "none",
  deployStatus: "none",
})

test("groups single commit into commit stage", async () => {
  const commits: Commit[] = [
    { }
  ]
  
  const actual = groupCommits(commits)
})

