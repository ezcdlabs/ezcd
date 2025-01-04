import * as dateFns from "date-fns";

export default function getShortDurationFormatFromSeconds(seconds: number) {
  const duration = dateFns.intervalToDuration({
    start: 0,
    end: seconds * 1000,
  });

  if (duration.days) {
    const days = duration.days + (duration.hours ?? 0) / 24;
    return `${days.toFixed(1)}d`;
  }

  if (duration.hours) {
    const hours = duration.hours + (duration.minutes ?? 0) / 60;
    return `${hours.toFixed(1)}h`;
  }

  if (duration.minutes) {
    const minutes = duration.minutes + (duration.seconds ?? 0) / 60;
    return `${minutes.toFixed(1)}m`;
  }

  return `${duration.seconds ?? 0}s`;
}
