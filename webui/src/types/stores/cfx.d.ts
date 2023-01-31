declare namespace CfxState {
  interface Player {
    name: string;
    steamId: string;
    created_at: string;
    discord: string;
    last_updated: string;
  }
  interface Character {
    citizenid: number;
    last_updated: string;
    created_at: string;
    steamId: string;
    user: {
      steamId: string;
      name: string;
      license: string;
      discord: string;
      last_updated: string;
      created_at: string;
    };
    data: {
      citizenid: number;
      position: string;
      metadata: string;
      last_updated: string;
      created_at: string;
    };
    info: {
      citizenid: number;
      firstname: string;
      lastname: string;
      birthdate: string;
      gender: 0;
      nationality: string;
      phone: string;
      cash: number;
      last_updated: string;
      created_at: string;
    };
  }
}
