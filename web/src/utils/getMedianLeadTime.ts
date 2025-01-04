import * as dateFns from "date-fns";

export interface CommitForLeadTime {
  date: string;
  leadTimeCompletedAt?: string;
}

export default function getMedianLeadTime<T extends CommitForLeadTime>(
  commits: T[],
) {
  if (commits.length === 0) {
    return null;
  }

  // lead time if from date to leadTimeCompletedAt
  let leadTimes: number[] = [];

  for (let i = 0; i < commits.length; i++) {
    const commit = commits[i];

    if (commit.leadTimeCompletedAt) {
      const leadTime = dateFns.differenceInSeconds(
        new Date(commit.leadTimeCompletedAt),
        new Date(commit.date),
      );

      // if the lead time is negative, it's an invalid lead time and we discard it
      if (leadTime > 0) {
        leadTimes.push(leadTime);
      }
    }
  }

  return leadTimes.length ? median(leadTimes) : null;
}

function median(values: number[]): number {
  if (values.length === 0) {
    throw new Error("Cannot calculate median. Input array is empty.");
  }

  // Sorting values, preventing original array
  // from being mutated.
  values = [...values].sort((a, b) => a - b);

  const half = Math.floor(values.length / 2);

  return values.length % 2
    ? values[half]
    : (values[half - 1] + values[half]) / 2;
}
