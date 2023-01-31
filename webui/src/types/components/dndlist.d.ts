declare namespace DnDList {
  interface Props {
    title: string;
    elements: string[];
    /**
     * Text displayed when the DND list is empty
     */
    emptyListHolder?: string;
    maxWidth?: boolean;
  }
}
