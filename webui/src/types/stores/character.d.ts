declare namespace CharacterState {
  interface Data {
    steamid: string;
    cid: uint;
    firstname: string;
    lastname: string;
    birthdate: string;
    gender: string;
    nationality: string;
    phone: string;
    metadata: Record<string, any>;
    created_at: number;
    last_updated: number;
  }
  interface Bank {
    accountId: string;
    name: string;
    type: string;
    balance: number;
    permissions: {
      deposit: boolean;
      withdraw: boolean;
      transfer: boolean;
      transactions: boolean;
    };
  }
  interface Vehicle {
    vin: string;
    model: string;
    plate: string;
    fakeplate: string;
    state: string;
    garageId: string;
  }
}
