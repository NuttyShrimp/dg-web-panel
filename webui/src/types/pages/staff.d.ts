declare namespace Staff {
  interface DashboardInfo {
    activePlayers: number;
    queue: any[];
    queuedPlayers: number;
    joinEvents: SimpleTimeline.Entry[];
  }
}
