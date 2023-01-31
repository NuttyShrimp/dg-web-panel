declare namespace Logs {
  interface Log {
    short_message: string;
    _id: string;
    _logtype: string;
    full_message: Record<string, string>;
    // unix timestamp in SECONDS
    timestamp: number;
  }
}
