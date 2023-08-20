declare namespace Inventory {
  type Item = {
    id: string;
    name: string;
    inventory: string;
    position: string;
    hotkey: number;
    metadata: Record<string, any>;
    destroyDate: number;
  };
}
