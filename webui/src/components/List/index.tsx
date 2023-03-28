import { ReactNode, MouseEventHandler } from "react";
import "./style.scss";

declare interface ListProps {
  children: ReactNode;
  highlightHover?: boolean;
  hideOverflow?: boolean;
}

export const List = ({ children, highlightHover, hideOverflow }: ListProps) => {
  return (
    <div
      className={`list-container ${highlightHover ? "list-container-highlight" : ""} ${
        hideOverflow ? "list-container-no-overflow" : ""
      }`}
    >
      {children}
    </div>
  );
};

declare interface ListEntryProps {
  children: ReactNode;
  onClick?: MouseEventHandler<HTMLDivElement>;
}

List.Entry = ({ children, onClick }: ListEntryProps) => (
  <div className="list-entry" onClick={onClick}>
    {children}
  </div>
);
