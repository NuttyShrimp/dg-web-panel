declare namespace CfxState {
  interface Player {
    name: string;
    steamId: string;
    created_at: string;
    discordId: string;
    last_updated: string;
    points: {
      points: number;
      updated_at: string;
    };
  }
  interface Character {
    citizenid: number;
    last_updated: string;
    created_at: string;
    steamId: string;
    // Only in active char list
    serverId?: number;
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

  namespace RealEstate {
    interface Location {
      id: number;
      name: string;
      garage: string;
      clothing: string;
      stash: string;
      logout: string;
      access: Access[];
    }

    interface Access {
      locationId: number;
      owner: boolean;
      citizenId: number;
      character: Character;
    }
  }
}
