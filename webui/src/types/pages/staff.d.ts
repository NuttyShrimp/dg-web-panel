declare namespace Staff {
  interface QueuedPlayer {
    source: number;
    identifiers: Record<string, string>;
    name: string;
  }
  interface DashboardInfo {
    activePlayers: number;
    queue: QueuedPlayer[];
    queuedPlayers: number;
    joinEvents: SimpleTimeline.Entry[];
  }
}
