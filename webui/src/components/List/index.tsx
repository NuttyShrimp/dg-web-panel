import { Divider } from "@mantine/core";
import { Children, ReactNode } from "react";
import "./style.scss";

declare interface ListProps {
  children: ReactNode;
  highlightHover?: boolean;
  hideOverflow?: boolean;
}

export const List = ({ children, highlightHover, hideOverflow }: ListProps) => {
  const dividedChilds = Children.map(children, child => (
    <>
      {child}
      <Divider />
    </>
  ));
  return (
    <div
      className={`list-container ${highlightHover ? "list-container-highlight" : ""} ${
        hideOverflow ? "list-container-no-overflow" : ""
      }`}
    >
      {dividedChilds}
    </div>
  );
};

List.Entry = ({ children }: { children: ReactNode }) => <div className="list-entry">{children}</div>;
