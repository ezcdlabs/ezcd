import classNames from "../../utils/classNames";
import { createPolled } from "@solid-primitives/timer";
import getFriendlyDurationFormat from "../../utils/getFriendlyDurationFormat";

const now = createPolled(() => new Date(), 500);

export default function TimerDisplay({
  start,
  end,
}: {
  start: Date;
  end?: Date | null;
}) {
  return (
    <span
      class={
        classNames()
        // "font-mono text-right text-xs text-[#797979]"
        // end ? "" : "text-azure-radiance-400"
      }
    >
      {getFriendlyDurationFormat(start, end ?? now())}
    </span>
  );
}
