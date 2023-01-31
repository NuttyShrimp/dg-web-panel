import { FC, PropsWithChildren } from "react";
import { Title } from "@mantine/core";

import "./centerbox.scss";

export const Centerbox: FC<PropsWithChildren<{ title: string }>> = ({ children, title }) => {
  return (
    <div className={"center centerBox"}>
      <div className={"centerBox-inner"}>
        <Title>{title}</Title>
        {children}
      </div>
    </div>
  );
};
