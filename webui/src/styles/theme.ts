import { createTheme, DefaultMantineColor, MantineColorsTuple } from "@mantine/core";

export const theme = createTheme({
  colors: {
    "dg-prim": [
      "#C6CAEA",
      "#A9AFE0",
      "#8E96D7",
      "#767FCF",
      "#5E69C6",
      "#4955BE",
      "#3F4AAE",
      "#39439C",
      "#343D8C",
      "#2F377D",
    ],
    "dg-prim-dark": [
      "#E2E7EF",
      "#B5BFD2",
      "#929FB6",
      "#77849D",
      "#646F84",
      "#4E596E",
      "#3D475B",
      "#2F394C",
      "#242D40",
      "#131D2E",
    ],
    "dg-sec": [
      "#FFF8E3",
      "#FFE297",
      "#F9CB64",
      "#E6B341",
      "#E8A30A",
      "#D18D00",
      "#B37700",
      "#936100",
      "#785000",
      "#634100",
    ],
    "dg-tert": [
      "#FAE6EA",
      "#F3BDC9",
      "#EE97AB",
      "#EA7490",
      "#E85476",
      "#E1385F",
      "#D7234D",
      "#BD2245",
      "#A6203F",
      "#801D34",
    ],
  },
  primaryColor: "dg-prim",
  headings: { fontFamily: "Greycliff CF, sans serif" },
});

type ExtendedCustomColors = "dg-prim" | "dg-prim-dark" | "dg-sec" | "dg-tert" | DefaultMantineColor;

declare module "@mantine/core" {
  export interface MantineThemeColorsOverride {
    colors: Record<ExtendedCustomColors, MantineColorsTuple>;
  }
}
