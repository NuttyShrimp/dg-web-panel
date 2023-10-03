import { DEFAULT_THEME, type MantineTheme } from "@mantine/core";
import { theme } from "@src/styles/theme";
import HeadingTT from "@tiptap/extension-heading";
import { mergeAttributes } from "@tiptap/react";

const Heading = HeadingTT.extend({
  renderHTML({ HTMLAttributes, node }) {
    const combinedTheme = { ...DEFAULT_THEME.headings, ...theme.headings };
    const hasLevel = HeadingTT.options.levels.includes(node.attrs.level);
    const level = hasLevel ? node.attrs.level : this.options.levels[0];
    const lvlStyles = combinedTheme.sizes?.[`h${level}` as keyof MantineTheme["headings"]["sizes"]];

    return [
      `h${level}`,
      mergeAttributes(HeadingTT.options.HTMLAttributes, HTMLAttributes, {
        style: Object.entries({
          "font-family": combinedTheme.fontFamily,
          "font-weight": combinedTheme.fontWeight,
          "font-size": `${lvlStyles?.fontSize}px`,
          "line-height": lvlStyles?.lineHeight,
        }).reduce((styleStr, [k, v]) => styleStr + `${k}:${v};`, ""),
      }),
      0,
    ];
  },
});

export default Heading;
