export const classInfo = {
  A: {
    length: 1,
    points: 10,
  },
  B: {
    length: 3,
    points: 10,
  },
  C: {
    length: 7,
    points: 15,
  },
  D: {
    length: -1,
    points: 30,
  },
};
export const reasons: Record<string, keyof typeof classInfo> = {
  "Uit karakter gaan op een fout moment": "A",
  "Misbruik commands": "A",
  "Overtreding regels rond server restarts": "A",
  Lootboxing: "A",
  "Powergaming / FailRP": "B",
  "Overtreding New Life Rule": "B",
  "Misbruik van het meerdere character systeem": "B",
  "No Value of Life": "B",
  "Overschreiding regels rond criminele activiteiten": "B",
  "Imiteren van overheids instanties": "B",
  ERP: "B",
  "Combat Loggen": "B",
  "RDM / VDM": "B",
  Cheaten: "D",
  Metagaming: "D",
  Bugabuse: "A",
};
