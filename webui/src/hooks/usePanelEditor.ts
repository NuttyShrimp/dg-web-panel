import { EditorOptions, useEditor } from "@tiptap/react";
import Highlight from "@tiptap/extension-highlight";
import Placeholder from "@tiptap/extension-placeholder";
import Typography from "@tiptap/extension-typography";
import { Blockquote } from "@tiptap/extension-blockquote";
import { Bold } from "@tiptap/extension-bold";
import { BulletList } from "@tiptap/extension-bullet-list";
import { Code } from "@tiptap/extension-code";
import { CodeBlock } from "@tiptap/extension-code-block";
import { Document } from "@tiptap/extension-document";
import { Dropcursor } from "@tiptap/extension-dropcursor";
import { Gapcursor } from "@tiptap/extension-gapcursor";
import { HardBreak } from "@tiptap/extension-hard-break";
import { History } from "@tiptap/extension-history";
import { HorizontalRule } from "@tiptap/extension-horizontal-rule";
import { Italic } from "@tiptap/extension-italic";
import { ListItem } from "@tiptap/extension-list-item";
import { OrderedList } from "@tiptap/extension-ordered-list";
import { Paragraph } from "@tiptap/extension-paragraph";
import { Strike } from "@tiptap/extension-strike";
import { Text } from "@tiptap/extension-text";
import Heading from "@src/components/CommentEditor/extensions/heading";

export const usePanelEditor = (props: Partial<EditorOptions>) => {
  return useEditor({
    ...props,
    extensions: [
      Blockquote,
      Bold,
      BulletList,
      Code,
      CodeBlock,
      Document,
      Dropcursor,
      Gapcursor,
      HardBreak,
      History,
      HorizontalRule,
      Italic,
      ListItem,
      OrderedList,
      Paragraph,
      Strike,
      Text,
      Heading,
      Highlight,
      Typography,
      Placeholder.configure({
        placeholder: "Write something ...",
      }),
    ],
  });
};
