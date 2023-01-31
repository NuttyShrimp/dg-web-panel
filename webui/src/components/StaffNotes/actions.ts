import { openConfirmModal } from "@mantine/modals";
import { axiosInstance } from "@src/helpers/axiosInstance";
import { noteState } from "@src/stores/notes/state";
import { JSONContent } from "@tiptap/core";
import { useCallback } from "react";
import { useSetRecoilState } from "recoil";

export const fetchNotes = async () => {
  try {
    const res = await axiosInstance.get<Notes.Note[]>("/staff/notes");
    return res.data;
  } catch (e) {
    console.error(e);
    return [];
  }
};

export const createNote = async (note: JSONContent) => {
  try {
    await axiosInstance.post("/staff/notes", {
      note: JSON.stringify(note),
    });
  } catch (e) {
    console.error(e);
  }
};

export const saveNote = async (noteId: number, note: JSONContent) => {
  try {
    await axiosInstance.post(`/staff/notes/${noteId}`, {
      note: JSON.stringify(note),
    });
  } catch (e) {
    console.error(e);
  }
};

export const deleteNote = async (noteId: number) => {
  try {
    await axiosInstance.delete(`/staff/notes/${noteId}`);
  } catch (e) {
    console.error(e);
  }
};

export const useRefreshNotes = () => {
  const setNotes = useSetRecoilState(noteState.notes);

  const refresh = useCallback(async () => {
    const newNotes = await fetchNotes();
    setNotes(newNotes);
  }, [setNotes]);

  return {
    refresh,
  };
};
