import { fetchNotes } from "@src/components/StaffNotes/actions";
import { atom, selector } from "recoil";

export const noteState = {
  notes: atom<Notes.Note[]>({
    key: "notes-list",
    default: selector({
      key: "notes-list-default",
      get: fetchNotes,
    }),
  }),
  activeNote: atom<number>({
    key: "notes-selected",
    default: 0,
  }),
};
