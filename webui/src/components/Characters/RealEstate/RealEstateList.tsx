import { Center, Flex, Paper, Stack, Text, useMantineTheme } from "@mantine/core";
import { List } from "@src/components/List";
import { LoadingSpinner } from "@src/components/LoadingSpinner";
import { Link } from "@src/components/Router/Link";
import { basicGet } from "@src/lib/actions/basicReq";
import { useQuery } from "@tanstack/react-query";

export const RealEstateList = ({ cid }: { cid: number }) => {
  const { data, isError, isLoading, error } = useQuery<CfxState.RealEstate.Location[], Error>({
    queryKey: ["real-estate", cid],
    queryFn: () => basicGet(`/character/realestate/owned/${cid}`),
    refetchOnWindowFocus: false,
  });

  const theme = useMantineTheme();

  if (isLoading) {
    return <LoadingSpinner />;
  }

  if (isError) {
    return <Text>Failed to load real estate locations: {error.message}</Text>;
  }

  if (data.length === 0) {
    return <Text>Persoon bezit geen real estate</Text>;
  }

  return (
    <Center>
      <List>
        {data.map(loc => (
          <List.Entry key={loc.id}>
            <Stack w={"100%"} spacing={"xs"}>
              <Text weight={"bolder"}>{loc.name}</Text>
              <Flex wrap="wrap">
                {loc.access.map(acc => (
                  <Paper
                    p="xs"
                    m="xs"
                    shadow="lg"
                    radius={"sm"}
                    key={acc.citizenId}
                    bg={theme.colors.dark[6]}
                    withBorder={acc.owner}
                  >
                    <Link to={`/staff/characters/${acc.citizenId}`}>
                      <Text>
                        {acc.character.info.firstname} {acc.character.info.lastname}
                      </Text>
                    </Link>
                  </Paper>
                ))}
              </Flex>
            </Stack>
          </List.Entry>
        ))}
      </List>
    </Center>
  );
};
