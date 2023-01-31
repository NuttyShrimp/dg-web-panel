declare namespace AuthState {
  type UserInfo = {
    username: string;
    avatarUrl: string;
    roles: string[];
    error?: string;
  };
}
