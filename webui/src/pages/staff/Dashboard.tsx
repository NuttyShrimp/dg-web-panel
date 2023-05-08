import { Container, Divider, Grid, Stack } from "@mantine/core";
import { DnDList } from "@src/components/DnDList";
import { NumberCard } from "@src/components/NumberCard";
import { SimpleTimeline } from "@src/components/SimpleTimeline";
import { NoteList } from "@src/components/StaffNotes/NoteList";
import { axiosInstance } from "@src/helpers/axiosInstance";
import { useEffect, useState } from "react";

export const StaffDashboard = () => {
  const [info, setInfo] = useState<Staff.DashboardInfo>({
    queue: [],
    joinEvents: [],
    activePlayers: 0,
    queuedPlayers: 0,
  });

  const fetchInfo = async () => {
    try {
      const res = await axiosInstance.get<Staff.DashboardInfo>("/staff/dashboard");
      if (res.status !== 200) return;
      setInfo(res.data);
    } catch (e) {
      console.error("Failed to fetch dashboard info", e);
    }
  };

  useEffect(() => {
    fetchInfo();
  }, []);

  return (
    <Container my="sm" size="xl">
      <Grid gutter={"sm"}>
        <Grid.Col xs={4}>
          <p>Join Events</p>
          <Divider pb="sm" />
          <SimpleTimeline list={info?.joinEvents ?? []} />
        </Grid.Col>
        <Grid.Col xs={4}>
          <Stack>
            <NumberCard title="Active Players" count={info.activePlayers} />
            <NoteList />
          </Stack>
        </Grid.Col>
        <Grid.Col xs={4}>
          <Stack>
            <NumberCard title="Players in Queue" count={info.queuedPlayers} />
            <DnDList
              title={"Queue"}
              elements={info.queue.map(ply => `${ply.name} (${ply.identifiers.steam})`)}
              maxWidth
              emptyListHolder="No players in queue"
            />
          </Stack>
        </Grid.Col>
      </Grid>
    </Container>
  );
};
