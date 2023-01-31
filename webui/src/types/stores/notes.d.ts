declare namespace Notes {
  interface Note {
    id: number;
    note: string;
    user: {
      Username: string;
    };
    createdAt: string;
    updatedAt: string;
  }
}
