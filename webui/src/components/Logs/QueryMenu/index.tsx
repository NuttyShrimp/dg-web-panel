import { Button, Menu } from "@mantine/core";
import { closeAllModals, openModal } from "@mantine/modals";
import { ChevronLeftIcon } from "@primer/octicons-react";
import { InputsModal } from "@src/components/Modals/InputsModal";
import { fieldsToQuery, menus, queries } from "./actions";

export const QueryMenu = ({ setQuery }: { setQuery: (q: string) => void }) => {
  const createQuery = (actionName: keyof typeof queries) => {
    const data = queries[actionName];
    let query = data.prefix ? data.prefix + " AND " : "";
    data.types.forEach(logtype => {
      query += ` OR logtype:"${logtype}"`;
    });
    query = query.replace(/^ OR /, "");
    if (data.inputs) {
      query = `(${query}) AND `;
      openModal({
        title: "Query parameters",
        children: (
          <InputsModal
            fields={data.inputs}
            onSubmit={f => {
              Object.keys(f).forEach((fn: string) => {
                if (f[fn] === "") return;
                let subQuery = "(";
                fieldsToQuery[fn].forEach(fq => {
                  subQuery += `${fq}: "${f[fn]}" OR `;
                });
                subQuery = subQuery.replace(/ OR $/, "");
                subQuery += ") OR ";
                query += subQuery;
              });
              query = query.replace(/ (OR|AND) $/, "");
              setQuery(query);
              closeAllModals();
            }}
          />
        ),
      });
    } else {
      setQuery(query);
    }
  };

  return (
    <Menu>
      <Menu.Target>
        <Button>Create query</Button>
      </Menu.Target>
      <Menu.Dropdown>
        {menus.map(m => (
          <Menu position="left" trigger="hover" key={`staff-query-builder-${m.label}`}>
            <Menu.Target>
              <Menu.Item icon={<ChevronLeftIcon />}>{m.label}</Menu.Item>
            </Menu.Target>
            <Menu.Dropdown>
              {Object.entries(m.actions).map(([action, label]) => (
                <Menu.Item key={`staff-query-builder-${m.label}-${action}`} onClick={() => createQuery(action)}>
                  {label}
                </Menu.Item>
              ))}
            </Menu.Dropdown>
          </Menu>
        ))}
      </Menu.Dropdown>
    </Menu>
  );
};
