import { Button, Center, Container, Group, Stack, Text, useMantineTheme } from "@mantine/core";
import { openModal } from "@mantine/modals";
import { CheckIcon, XIcon } from "@primer/octicons-react";
import { List } from "@src/components/List";
import { axiosInstance } from "@src/helpers/axiosInstance";
import { characterState } from "@src/stores/character/state";
import { FC, useEffect } from "react";
import { useRecoilState } from "recoil";
import { ChangeBalanceModal } from "./modals/changeBalanceModal";

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

  const openBalanceProfile = (accId: string, balance: number) => {
    openModal({
      title: "Update Bank Balance",
      children: <ChangeBalanceModal balance={balance} accountId={accId} />,
    });
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
          <List.Entry key={acc.accountId}>
            <Stack w={"100%"} gap="xs">
              <Group justify="space-between">
                <Text fw={"bolder"}>{acc.name}</Text>
                <Button onClick={() => openBalanceProfile(acc.accountId, acc.balance)}>Change Balance</Button>
              </Group>
              <Text>€{acc.balance}</Text>
              <Group grow>
                <Group gap={"xs"}>
                  <Text>Deposit</Text>
                  {acc.permissions.deposit ? (
                    <CheckIcon fill={theme.colors.green[6]} />
                  ) : (
                    <XIcon fill={theme.colors.red[6]} />
                  )}
                </Group>
                <Group gap={"xs"}>
                  <Text>Withdraw</Text>
                  {acc.permissions.withdraw ? (
                    <CheckIcon fill={theme.colors.green[6]} />
                  ) : (
                    <XIcon fill={theme.colors.red[6]} />
                  )}
                </Group>
                <Group gap={"xs"}>
                  <Text>Transfer</Text>
                  {acc.permissions.transfer ? (
                    <CheckIcon fill={theme.colors.green[6]} />
                  ) : (
                    <XIcon fill={theme.colors.red[6]} />
                  )}
                </Group>
                <Group gap={"xs"}>
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
