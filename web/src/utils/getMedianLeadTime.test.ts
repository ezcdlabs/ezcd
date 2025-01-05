import { describe, expect, test } from "vitest";
import getMedianLeadTime, {
  CommitForLeadTime,
} from "../utils/getMedianLeadTime";

function makeCommit(commit: Partial<CommitForLeadTime>): CommitForLeadTime {
  return {
    date: "2021-01-01T09:00:00Z",
    ...commit,
  };
}

test("should return empty metrics for empty commits", () => {
  const actual = getMedianLeadTime([]);
  expect(actual).toEqual(null);
});

test("should return lead time of the only commit", () => {
  const commit1 = makeCommit({
    date: "2021-01-01T09:00:00Z",
    leadTimeCompletedAt: "2021-01-01T12:00:00Z",
  });

  const actual = getMedianLeadTime([commit1]);
  expect(actual).toEqual(3 * 60 * 60);
});

test("should return lead time of the middle commit", () => {
  const commit1 = makeCommit({
    date: "2021-01-01T09:00:00Z",
    leadTimeCompletedAt: "2021-01-01T12:00:00Z",
  });
  const commit2 = makeCommit({
    date: "2021-01-01T09:00:00Z",
    leadTimeCompletedAt: "2021-01-01T13:00:00Z",
  });
  const commit3 = makeCommit({
    date: "2021-01-01T09:00:00Z",
    leadTimeCompletedAt: "2021-01-01T14:00:00Z",
  });

  const actual = getMedianLeadTime([commit1, commit3, commit2]);
  expect(actual).toEqual(4 * 60 * 60);
});

test("should return lead time between the middle commits", () => {
  const commit1 = makeCommit({
    date: "2021-01-01T09:00:00Z",
    leadTimeCompletedAt: "2021-01-01T12:00:00Z",
  });
  const commit2 = makeCommit({
    date: "2021-01-01T09:00:00Z",
    leadTimeCompletedAt: "2021-01-01T13:00:00Z",
  });
  const commit3 = makeCommit({
    date: "2021-01-01T09:00:00Z",
    leadTimeCompletedAt: "2021-01-01T14:00:00Z",
  });
  const commit4 = makeCommit({
    date: "2021-01-01T09:00:00Z",
    leadTimeCompletedAt: "2021-01-01T15:00:00Z",
  });

  const actual = getMedianLeadTime([commit1, commit3, commit2, commit4]);
  expect(actual).toEqual(4.5 * 60 * 60);
});
