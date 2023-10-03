import { TextInput } from "@mantine/core";
import { ChangeEventHandler, EventHandler, FC, useState } from "react";
import { animated, easings, useTransition } from "react-spring";
import "./styles.scss";

declare interface SearchInputProps {
  value: string;
  onChange: (value: string) => void;
}

export const SearchInput: FC<SearchInputProps> = props => {
  const [focused, setFocused] = useState(false);
  const [showIcon, setShowIcon] = useState(true);
  const [searchValue, setSearchValue] = useState<string>(props.value);

  const focus = () => {
    setFocused(true);
    setShowIcon(false);
  };

  const unFocus: EventHandler<any> = (e: Event) => {
    if (searchValue !== "") return;
    if (document.activeElement === e.currentTarget) return;
    setFocused(false);
  };

  const updateValue: ChangeEventHandler<HTMLInputElement> = e => {
    setSearchValue(e.currentTarget.value);
    props.onChange(e.currentTarget.value);
  };

  const transitions = useTransition(focused, {
    from: {
      width: 50,
    },
    enter: {
      width: 200,
    },
    leave: {
      width: 50,
    },
    config: {
      duration: 250,
      easing: easings.linear,
    },
    onRest: () => {
      if (!focused) setShowIcon(true);
    },
  });

  return (
    <div className="searchInput-wrapper">
      {transitions(
        (styles, focused) =>
          focused && (
            <animated.div style={styles}>
              <TextInput
                leftSection={<i className="fas fa-magnifying-glass" />}
                placeholder="Search"
                radius={"xl"}
                onChange={updateValue}
                onMouseOut={unFocus}
                onBlur={unFocus}
              />
            </animated.div>
          )
      )}
      {showIcon && (
        <div onMouseOver={focus} onFocus={focus}>
          <i className="fas fa-magnifying-glass" />
        </div>
      )}
    </div>
  );
};
