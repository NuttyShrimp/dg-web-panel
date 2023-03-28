declare namespace Logs {
  interface Response {
    logs: Log[];
    total: number;
  }
  interface Log {
    short_message: string;
    _id: string;
    _logtype: string;
    full_message: Record<string, string>;
    // unix timestamp in SECONDS
    timestamp: number;
  }
}
