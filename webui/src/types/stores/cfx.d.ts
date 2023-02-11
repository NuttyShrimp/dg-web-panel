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
  interface Penalty {
    id: number;
    steamId: string;
    penalty: "ban" | "warn" | "kick";
    reason: string;
    points: number;
    length: number;
    date: string;
    automated: string;
  }
  namespace Business {
    interface Type {
      id: number;
      name: string;
    }
    interface Entry {
      name: string;
      label: string;
      id: number;
      bankAccountId: string;
      type: Type;
    }
    interface Log {
      id: number;
      type: string;
      action: string;
      businessId: number;
      character: Omit<Character, "info">;
    }
    interface Role {
      id: number;
      name: string;
      permissions: string[];
    }
    interface Employee {
      id: number;
      isOwner: bool;
      character: Character;
      role: Role;
    }
  }
}
