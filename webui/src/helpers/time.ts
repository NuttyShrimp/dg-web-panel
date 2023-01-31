import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";

import "dayjs/locale/nl-be";

export const formatRelativeTime = (time: number) => {
  dayjs.extend(relativeTime).locale("nl-be");
  return dayjs.unix(time).fromNow();
};

export const displayDate = (timeStr: string) => {
  return dayjs(timeStr).format("DD MMM YYYY");
};

export const displayUnixDate = (time: number) => {
  dayjs.locale("nl-be");
  return dayjs.unix(time).format("DD/MM/YYYY - HH:mm");
};

// Time here refs to time.Time in golang, when serialized in JSON it is formatted as ISO8601 like string
export const displayTimeDate = (time: string) => {
  dayjs.locale("nl-be");
  return dayjs(time).format("DD/MM/YYYY - HH:mm");
};
