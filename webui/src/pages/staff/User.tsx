import { Container, Divider, Flex, Group, Tabs, Text, Title } from "@mantine/core";
import { showNotification } from "@mantine/notifications";
import { InfoIcon } from "@primer/octicons-react";
import { FontAwesomeIcon } from "@src/components/Icon";
import { UserActionMenu } from "@src/components/Users/ActionMenu";
import { UserBanStatus } from "@src/components/Users/BanStatus";
import { UserCharList } from "@src/components/Users/CharList";
import { UserPenaltyList } from "@src/components/Users/PenaltyList";
import { displayDate, displayTimeDate } from "@src/helpers/time";
import { cfxState } from "@src/stores/cfx/state";
import dayjs from "dayjs";
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
      <Tabs variant="pills" defaultValue={"info"} keepMounted={false} pt={"sm"}>
        <Tabs.List mb={"xs"}>
          <Tabs.Tab value="info" icon={<InfoIcon size={14} />}>
            Info
          </Tabs.Tab>
          <Tabs.Tab value="chars" icon={<FontAwesomeIcon icon="users" size={"sm"} />}>
            Characters
          </Tabs.Tab>
          <Tabs.Tab value="penalties" icon={<FontAwesomeIcon icon="hammer-war" size={"sm"} />}>
            Penalties
          </Tabs.Tab>
          <UserActionMenu steamId={steamid} />
        </Tabs.List>
        <Tabs.Panel value="info">
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
              <Text weight={"bolder"}>DiscordId:&nbsp;</Text>
              <Text>{player.discord}</Text>
            </Flex>
            <Flex>
              <Text weight={"bolder"}>First Joined:&nbsp;</Text>
              <Text>{displayDate(player.created_at)}</Text>
            </Flex>
            <Flex>
              <Text weight={"bolder"}>Last Update:&nbsp;</Text>
              <Text>{displayDate(player.last_updated)}</Text>
            </Flex>
            <Flex>
              <Text weight={"bolder"}>Points:&nbsp;</Text>
              <Text>{player.points.points}</Text>
            </Flex>
            {player.points.points > 0 && (
              <Flex>
                <Text weight={"bolder"}>Points reset-day:&nbsp;</Text>
                <Text>
                  {displayTimeDate(dayjs(player.points.updated_at).add(player.points.points, "d").toString())}
                </Text>
              </Flex>
            )}
          </div>
        </Tabs.Panel>
        <Tabs.Panel value="chars">
          <UserCharList steamId={steamid} />
        </Tabs.Panel>
        <Tabs.Panel value="penalties">
          <UserPenaltyList steamId={steamid} />
        </Tabs.Panel>
      </Tabs>
    </Container>
  );
};
