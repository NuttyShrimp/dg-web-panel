import { Combobox, ComboboxItem, TextInput, useCombobox } from "@mantine/core";
import { useMemo, useState } from "react";

interface AutoCompleteProps<T = Record<string, any> & ComboboxItem> {
  itemComponent: React.FC<any>;
  data: T[];
  placeholder?: string;
  defaultValue?: string;
  filter: (item: T, search: string) => boolean;
  onOptionSubmit: (item: string | null) => void;
}

export const AutoComplete = ({
  itemComponent,
  data,
  placeholder,
  defaultValue,
  filter,
  onOptionSubmit,
}: AutoCompleteProps) => {
  const combobox = useCombobox({
    onDropdownClose: () => combobox.resetSelectedOption(),
  });

  const [value, setValue] = useState<string>(defaultValue ?? "");

  const filteredData = useMemo(() => data.filter(v => filter(v, value)), [value, data, filter]);

  const options = useMemo(
    () =>
      filteredData.map(item => (
        <Combobox.Option value={item.value} key={item.value}>
          {itemComponent && itemComponent(item)}
        </Combobox.Option>
      )),
    [itemComponent, filteredData]
  );

  return (
    <Combobox
      store={combobox}
      onOptionSubmit={val => {
        setValue(val);
        onOptionSubmit?.(val);
        combobox.closeDropdown();
      }}
    >
      <Combobox.Target>
        <TextInput
          label={placeholder}
          placeholder={placeholder}
          value={value}
          onChange={event => {
            setValue(event.currentTarget.value);
            combobox.openDropdown();
            combobox.updateSelectedOptionIndex();
          }}
          onClick={() => combobox.openDropdown()}
          onFocus={() => combobox.openDropdown()}
          onBlur={() => combobox.closeDropdown()}
        />
      </Combobox.Target>

      <Combobox.Dropdown>
        <Combobox.Options>{options}</Combobox.Options>
      </Combobox.Dropdown>
    </Combobox>
  );
};
