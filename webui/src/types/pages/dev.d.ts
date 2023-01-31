declare namespace Dev {
  interface APIKey {
    key: string;
    comment: string;
    expiry: string;
    userId: number;
    User: {
      Username: string;
    };
  }
  interface CacheControlEntry {
    label: string;
    endpoint: string;
  }
}
