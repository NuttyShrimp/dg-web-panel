import { Flex, Stack, Text } from "@mantine/core";
import { cfxState } from "@src/stores/cfx/state";
import { useEffect } from "react";
import { useRecoilRefresher_UNSTABLE, useRecoilValue } from "recoil";
import { List } from "../List";
import { Link } from "../Router/Link";

export const BusinessEmployees = ({ id }: { id: number }) => {
  const employees = useRecoilValue(cfxState.businessEmployees(id));
  const refreshList = useRecoilRefresher_UNSTABLE(cfxState.businessEmployees(id));
  useEffect(() => {
    refreshList();
  }, []);
  return (
    <List>
      {employees.map(e => (
        <List.Entry key={e.id}>
          <Stack spacing={4}>
            <Text weight={"bolder"}>
              {e.character.info.firstname} {e.character.info.lastname} |{" "}
              <Link to={`/staff/characters/${e.character.citizenid}`}>{e.character.citizenid}</Link> |{" "}
              <Link to={`/staff/users/${e.character.steamId}`}>{e.character.steamId}</Link>
            </Text>
            <Text>
              {e.role.name} | {e.role.permissions.join(",")}
            </Text>
          </Stack>
        </List.Entry>
      ))}
    </List>
  );
};
