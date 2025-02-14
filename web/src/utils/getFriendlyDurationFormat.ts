import * as dateFns from "date-fns";

export default function getFriendlyDurationFormat(start: Date, end: Date) {
  const nearestStartSecond = dateFns.startOfSecond(start);

  const duration = dateFns.intervalToDuration({
    start: nearestStartSecond,
    end,
  });
  const days = dateFns.differenceInDays(end, nearestStartSecond);

  let output = "";
  if (days > 0) {
    output += `${days}d `;
  }
  if (duration?.hours || output) {
    output += `${(duration.hours ?? 0)
      .toString()
      .padStart(output ? 2 : 1, "0")}h `;
  }

  if (duration?.minutes || output) {
    output += `${(duration.minutes ?? 0)
      .toString()
      .padStart(output ? 2 : 1, "0")}m `;
  }

  output += `${(duration.seconds ?? 0)
    .toString()
    .padStart(output ? 2 : 1, "0")}s`;
  return output;
}
