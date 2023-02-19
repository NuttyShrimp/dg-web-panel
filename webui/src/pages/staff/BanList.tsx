import { Container, Flex, Stack, Text, Title } from "@mantine/core";
import { UnbanActionMenu } from "@src/components/Bans/UnbanActionMenu";
import { List } from "@src/components/List";
import { LoadingSpinner } from "@src/components/LoadingSpinner";
import { Link } from "@src/components/Router/Link";
import { basicGet } from "@src/lib/actions/basicReq";
import { useQuery } from "@tanstack/react-query";
import dayjs from "dayjs";

export const BanListPage = () => {
  const {
    data: banList,
    error,
    isLoading,
    isError,
  } = useQuery<CfxState.Penalty[], Error>({
    queryKey: ["ban-list"],
    queryFn: () => basicGet<CfxState.Penalty[]>("/staff/ban/list"),
  });

  return (
    <Container>
      <Title order={2}>Ban list</Title>
      {isLoading ? (
        <LoadingSpinner />
      ) : isError ? (
        <Text>Failed to load banlist: {error.message}</Text>
      ) : (
        <List>
          {banList.map(b => (
            <List.Entry key={b.id}>
              <Flex justify={"space-between"} w={"100%"}>
                <Stack spacing={4}>
                  <Link to={`/staff/players/${b.steamId}`} style={{ width: "100%" }}>
                    <Text weight={"bolder"}>{b.steamId}</Text>
                  </Link>
                  <Flex>
                    <Text weight={"bolder"}>Reason:</Text>
                    <Text>&nbsp;{b.reason}</Text>
                  </Flex>
                  <Flex>
                    <Text weight={"bolder"}>Until:</Text>
                    <Text>
                      &nbsp;
                      {b.length === -1 ? "Permanently" : dayjs(b.date).add(b.length, "d").format("DD/MM/YYYY HH:mm")}
                    </Text>
                  </Flex>
                </Stack>
                <Flex align={"center"}>
                  <UnbanActionMenu penalty={b} />
                </Flex>
              </Flex>
            </List.Entry>
          ))}
        </List>
      )}
    </Container>
  );
};
