interface QueryType {
  types: string[];
  inputs: string[];
  prefix?: string;
}

interface MenuEntry {
  label: string;
  actions: Record<string, keyof typeof queries>;
}

export const menus: MenuEntry[] = [
  {
    label: "Admin",
    actions: {
      adminActions: "Actions",
      adminOpenClose: "Open/Close menu",
      adminDevmode: "Toggled devmode",
      coreJoinLeave: "Join/Left",
      policeDroppedWhileCuffed: "Dropped while cuffed",
      charSpawned: "Char spawn locations",
    },
  },
  {
    label: "Anticheat",
    actions: {
      actriggers: "Triggered Anticheat",
    },
  },
  {
    label: "Chat",
    actions: {
      chatCmds: "Messages",
    },
  },
  {
    label: "Weather",
    actions: {
      blackout: "Blackout",
    },
  },
  {
    label: "Business",
    actions: {
      blackout: "Actions",
    },
  },
  {
    label: "Criminal",
    actions: {
      criminalWeed: "Weed",
      criminalFence: "Fence",
      criminalCornersell: "Cornersell",
      criminalHeists: "Heists",
      criminalStoreRobs: "Store robbery",
      criminalHouseRobs: "House robbery",
    },
  },
  {
    label: "Farming",
    actions: {
      farmingView: "Viewed plants",
      farmingFeed: "Feed plant",
      farmingWater: "Water plant",
      farmingHarvest: "Harvested plant",
    },
  },
  {
    label: "Financials",
    actions: {
      financialsTrans: "Transacties",
      financialsDeposit: "Deposit",
      financialsWithdraw: "Withdraw",
      financialsCash: "Cash",
    },
  },
  {
    label: "Hospital",
    actions: {
      hospitalStatus: "Status veranderingen",
    },
  },
  {
    label: "Inventory",
    actions: {
      crafting: "Crafting",
      usedItems: "Used/Bought Item",
      destroyedItems: "Destroyed items",
      dropActions: "Move from/to ground",
      plysRobbed: "Robbed player",
    },
  },
  {
    label: "Jobs",
    actions: {
      jobsGroupCreations: "Group creation/disband",
      jobsDutyToggle: "Duty toggle",
    },
  },
  {
    label: "Mechanic",
    actions: {
      mechRepairs: "Repairs",
      mechUpgrades: "Upgrades",
    },
  },
  {
    label: "Scenes",
    actions: {
      sceneActions: "Create/Delete",
    },
  },
  {
    label: "Weapons",
    actions: {
      weaponEquip: "(Un)equip,",
    },
  },
  {
    label: "LogList",
    actions: {
      loglist: "Used",
    },
  },
];

export const queries: Record<string, QueryType> = {
  adminActions: {
    types: [
      "admin:menu:open",
      "admin:menu:close",
      "admin:menu:action",
      "admin:announce",
      "admin:whitelist:refresh",
      "admin:selector:action",
    ],
    inputs: ["steamId"],
  },
  adminOpenClose: {
    types: ["admin:menu:open", "admin:menu:close"],
    inputs: ["steamId"],
  },
  adminDevmode: {
    types: ["admin:devmode"],
    inputs: ["steamId"],
  },
  actriggers: {
    types: ["anticheat:flag", "anticheat:damageModifier", "anticheat:explosion"],
    inputs: ["steamId"],
  },
  chatCmds: {
    types: ["chat:ooc:local", "chat:ooc:global", "chat:me", "chat:message"],
    inputs: ["steamId", "citizenid"],
  },
  blackout: {
    types: ["chat:ooc:local", "chat:ooc:global", "chat:me", "chat:message"],
    inputs: ["citizenid"],
  },
  businessActions: {
    types: [
      "business:hire",
      "business:fire",
      "business:changeBankPermission",
      "business:assignRole",
      "business:createRole",
      "business:updateRole",
      "business:deleteRole",
    ],
    inputs: ["citizenid", "business"],
  },
  criminalWeed: {
    types: ["weed:feed", "weed:destroy", "weed:remove", "weed:cut", "weed:planted"],
    inputs: ["citizenid"],
  },
  criminalFence: {
    types: ["criminal:fence:sell", "criminal:fence:take"],
    inputs: ["citizenid"],
  },
  criminalCornersell: {
    types: ["cornersell:start", "cornersell:sell"],
    inputs: ["citizenid"],
  },
  criminalHeists: {
    types: [
      "heists:fleeca:power",
      "heists:hack:start",
      "heists:laptop:failed",
      "heists:laptop:success",
      "heists:laptop:buy",
      "heists:laptop:pickup",
      "heists:trolley:loot",
      "heists:trolley:specialLoot",
    ],
    inputs: ["citizenid"],
  },
  criminalStoreRobs: {
    types: [
      "storerobbery:registers:stateChange",
      "storerobbery:safes:hack",
      "storerobbery:safe:rob",
      "storerobbery:safes:canceledHack",
    ],
    inputs: ["citizenid"],
  },
  criminalHouseRobs: {
    types: [
      "houserobbery:signin:logout",
      "houserobbery:signin:login",
      "houserobbery:door:unlock",
      "houserobery:door:lock",
      "houserobbery:house:enter",
      "houserobbery:house:leave",
      "houserobbery:house:loot",
      "houserobbery:job:findtimeout",
      "houserobbery:job:start",
      "houserobbery:job:finish",
    ],
    inputs: ["citizenid"],
  },
  hospitalStatus: {
    types: [
      "hospital:down:died",
      "hospital:down:stageChanged",
      "hospital:down:respawnToHospital",
      "hospital:down:respawnToBed",
      "hospital:revive",
    ],
    inputs: ["citizenid", "steamId"],
  },
  sceneActions: {
    types: ["scenes:created", "scenes:deleted"],
    inputs: ["citizenid", "steamId"],
  },
  coreJoinLeave: {
    types: ["core:left", "core:crashed", "core:joined", "core:created", "core:select"],
    inputs: ["citizenid", "steamId"],
  },
  crafting: {
    types: ["materials:crafting:crafted"],
    inputs: ["citizenid", "steamId"],
  },
  weaponEquip: {
    types: ["weapons:equip", "weapons:unequip"],
    inputs: ["citizenid", "steamId"],
  },
  mechRepairs: {
    types: ["vehicles:status:updatePart"],
    inputs: ["citizenid", "steamId", "vin"],
  },
  mechUpgrades: {
    types: ["vehicles:upgrades:performace", "vehicles:upgrades:cosmetic"],
    inputs: ["vin"],
  },
  usedItems: {
    types: ["inventory:item:used", "inventory:item:bought"],
    inputs: ["citizenid", "steamId", "itemName"],
  },
  destroyedItems: {
    types: ["inventory:item:destroyed"],
    inputs: ["citizenid", "steamId", "itemName"],
  },
  dropActions: {
    types: ["inventory:item:moved"],
    inputs: ["citizenid", "steamId", "itemName"],
    prefix: "(newInventory:/drop__.*/ OR oldInventory:/drop__.*/)",
  },
  farmingView: {
    types: ["farming:plant:view"],
    inputs: ["citizenid", "steamId"],
  },
  farmingHarvest: {
    types: ["farming:plant:harvest"],
    inputs: ["citizenid", "steamId"],
  },
  farmingFeed: {
    types: ["farming:plant:feed"],
    inputs: ["citizenid", "steamId"],
  },
  farmingWater: {
    types: ["farming:plant:water"],
    inputs: ["citizenid", "steamId"],
  },
  jobsGroupCreations: {
    types: ["jobs:groupmanger:create", "jobs:groupmanger:disband"],
    inputs: ["citizenid", "steamId"],
  },
  jobsDutyToggle: {
    types: ["jobs:whitelist:signin:success", "jobs:whitelist:signout:success"],
    inputs: ["citizenid", "steamId"],
  },
  plysRobbed: {
    types: ["police:interactions:robbedPlayer"],
    inputs: ["citizenid", "steamId"],
  },
  financialsTrans: {
    types: ["financials:transfer:success", "financials:purchase:success", "financials:paycheck:success"],
    inputs: ["accountCid", "originAccountId", "targetAccountId"],
  },
  financialsDeposit: {
    types: ["financials:deposit:success"],
    inputs: ["accountCid", "originAccountId"],
  },
  financialsWithdraw: {
    types: ["financials:withdraw:success"],
    inputs: ["accountCid", "originAccountId"],
  },
  financialsCash: {
    types: ["cash:remove", "cash:add"],
    inputs: ["citizenid"],
  },
  loglist: {
    types: ["logtype:idlist:use"],
    inputs: ["citizenid", "steamId"],
  },
  policeDroppedWhileCuffed: {
    types: ["police:interactions:droppedWithCuffs"],
    inputs: ["citizenid"],
  },
  charSpawned: {
    types: ["chars:spawn"],
    inputs: ["citizenid", "steamId"],
  },
};

export const fieldsToQuery: Record<string, string[]> = {
  steamId: ["plyInfo_steamId", "steamId"],
  citizenid: ["plyInfo_cid", "cid"],
  accountCid: ["cid", "acceptor_cid"],
  vin: ["vin"],
  itemName: ["itemName"],
  targetAccountId: ["targetAccountId"],
  originAccountId: ["account"],
};
