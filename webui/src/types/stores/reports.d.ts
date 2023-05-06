declare namespace ReportState {
  type Filter = {
    open: boolean;
    closed: boolean;
    search: string;
  };
  type Member = {
    steamId: string;
    name: string;
  };
  type Report = {
    id: number;
    title: string;
    open: boolean;
    /**
     * ISO-formatted date string
     */
    createdAt: string;
    /**
     * ISO-formatted date string
     */
    updatedAt: string;
    members?: Member[];
  };
  type Message = {
    id: number;
    message: string;
    type: "text" | "image";
    createdAt: string;
    updatedAt: string;
    deletedAt: null;
    sender: MessageSender;
  };
  type MessageSender = {
    username: string;
    avatarUrl: string;
    roles: string[];
    steamId?: string;
  };
}
