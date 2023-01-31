import { Text } from "@mantine/core";
import { FC } from "react";

export const FontAwesomeIcon: FC<Icon.Props> = ({ icon, lib, ...props }) => {
  return <Text component="i" className={`${lib ?? "fas"} fa-${icon}`} {...props} />;
};
