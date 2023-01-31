import { Center, Container, Group, Stack, Text, useMantineTheme } from "@mantine/core";
import { CheckIcon, XIcon } from "@primer/octicons-react";
import { List } from "@src/components/List";
import { axiosInstance } from "@src/helpers/axiosInstance";
import { characterState } from "@src/stores/character/state";
import { FC, useEffect } from "react";
import { useRecoilState } from "recoil";

export const BankInfo: FC<{ cid: number }> = ({ cid }) => {
  const [characterBankAccs, setCharacterBankAccs] = useRecoilState(characterState.bank);
  const theme = useMantineTheme();
  const fetchData = async () => {
    try {
      const res = await axiosInstance.get<CharacterState.Bank[]>(`/character/bank/${cid}`);
      if (res.status !== 200) return;
      return res.data;
    } catch (e) {
      console.error(e);
    }
  };

  useEffect(() => {
    let ignore = false;
    if (!characterBankAccs) {
      const doFetch = async () => {
        const data = await fetchData();
        if (data && !ignore) {
          setCharacterBankAccs(data);
        }
      };
      doFetch();
    }
    return () => {
      ignore = true;
    };
  }, []);
  if (!characterBankAccs) {
    return (
      <Container>
        <Center>
          <Text>This character has no bank accounts susge</Text>
        </Center>
      </Container>
    );
  }
  // TODO: Link bank account to bank account page
  return (
    <Center>
      <List>
        {characterBankAccs.map(acc => (
          <List.Entry key={acc.account_id}>
            <Stack w={"100%"} spacing="xs">
              <Text weight={"bolder"}>{acc.name}</Text>
              <Text>€{acc.balance}</Text>
              <Group grow>
                <Group spacing={"xs"}>
                  <Text>Deposit</Text>
                  {acc.permissions.deposit ? (
                    <CheckIcon fill={theme.colors.green[6]} />
                  ) : (
                    <XIcon fill={theme.colors.red[6]} />
                  )}
                </Group>
                <Group spacing={"xs"}>
                  <Text>Withdraw</Text>
                  {acc.permissions.withdraw ? (
                    <CheckIcon fill={theme.colors.green[6]} />
                  ) : (
                    <XIcon fill={theme.colors.red[6]} />
                  )}
                </Group>
                <Group spacing={"xs"}>
                  <Text>Transfer</Text>
                  {acc.permissions.transfer ? (
                    <CheckIcon fill={theme.colors.green[6]} />
                  ) : (
                    <XIcon fill={theme.colors.red[6]} />
                  )}
                </Group>
                <Group spacing={"xs"}>
                  <Text>View transitions</Text>
                  {acc.permissions.transactions ? (
                    <CheckIcon fill={theme.colors.green[6]} />
                  ) : (
                    <XIcon fill={theme.colors.red[6]} />
                  )}
                </Group>
              </Group>
            </Stack>
          </List.Entry>
        ))}
      </List>
    </Center>
  );
};
