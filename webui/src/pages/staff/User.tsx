import { Container, Divider, Flex, Group, Text, Title } from "@mantine/core";
import { showNotification } from "@mantine/notifications";
import { UserBanStatus } from "@src/components/Users/BanStatus";
import { displayDate } from "@src/helpers/time";
import { cfxState } from "@src/stores/cfx/state";
import { useMemo } from "react";
import { Navigate, useParams } from "react-router-dom";
import { useRecoilValue } from "recoil";

export const UserPage = () => {
  const { steamid } = useParams();
  const players = useRecoilValue(cfxState.players);
  const player = useMemo(() => {
    return players.find(p => p.steamId === steamid);
  }, [steamid, players]);

  if (!player || !steamid) {
    showNotification({
      title: "Player not found",
      message: "seems like the player you were searching doesn't exists",
      color: "red",
    });
    return <Navigate to="/staff/users" />;
  }

  return (
    <Container>
      <Group>
        <Title>{player.name}</Title>
        <UserBanStatus steamId={steamid} />
      </Group>
      <Divider mb={"xs"} />
      <div>
        <Flex>
          <Text weight={"bolder"}>Name:&nbsp;</Text>
          <Text>{player.name}</Text>
        </Flex>
        <Flex>
          <Text weight={"bolder"}>SteamId:&nbsp;</Text>
          <Text>{player.steamId}</Text>
        </Flex>
        <Flex>
          <Text weight={"bolder"}>First Joined:&nbsp;</Text>
          <Text>{displayDate(player.created_at)}</Text>
        </Flex>
        <Flex>
          <Text weight={"bolder"}>Last Update:&nbsp;</Text>
          <Text>{displayDate(player.last_updated)}</Text>
        </Flex>
      </div>
    </Container>
  );
};
