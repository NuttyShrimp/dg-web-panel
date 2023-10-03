import { ActionIcon, Button, Card, Divider, Paper, ScrollArea, Stack, Text } from "@mantine/core";
import { openConfirmModal } from "@mantine/modals";
import { PlusIcon, TrashIcon } from "@primer/octicons-react";
import { displayTimeDate } from "@src/helpers/time";
import { usePanelEditor } from "@src/hooks/usePanelEditor";
import { noteState } from "@src/stores/notes/state";
import { EditorContent } from "@tiptap/react";
import { FC, useEffect } from "react";
import { useRecoilValue, useSetRecoilState } from "recoil";
import { createNote, deleteNote, saveNote, useRefreshNotes } from "./actions";

import "./style.scss";

const NoteEntry: FC<{ note: Notes.Note; selected: boolean }> = ({ note, selected }) => {
  const setActiveNote = useSetRecoilState(noteState.activeNote);
  const { refresh } = useRefreshNotes();
  const editor = usePanelEditor({
    content: JSON.parse(note.note) ?? "",
    editable: false,
    onBlur: ({ editor }) => saveNote(note.id, editor.getJSON()),
  });

  useEffect(() => {
    editor?.setEditable(selected);
  }, [selected, editor]);

  const openDelNoteModal = () => {
    openConfirmModal({
      color: "red",
      title: "Delete note",
      children: <Text>Are you sure you want to delete this note created by {note.user.Username}</Text>,
      labels: { cancel: "Cancel", confirm: "Delete" },
      onCancel: () => {
        // empty
      },
      onConfirm: async () => {
        await deleteNote(note.id);
        refresh();
      },
    });
  };

  return (
    <Card shadow="xs" radius="md" p="sm" onClick={() => setActiveNote(note.id)} onBlur={() => setActiveNote(0)}>
      <EditorContent editor={editor} />
      <Divider pb={"xs"} />
      <Text>
        created by {note.user.Username} - {displayTimeDate(note.createdAt)}
        <br />
        last updated at {displayTimeDate(note.updatedAt)}
      </Text>
      {selected && (
        <div className="note-delete">
          <ActionIcon color={"red"} onClick={() => openDelNoteModal()}>
            <TrashIcon />
          </ActionIcon>
        </div>
      )}
    </Card>
  );
};

export const NoteList = () => {
  const notes = useRecoilValue(noteState.notes);
  const selNote = useRecoilValue(noteState.activeNote);
  const newNoteEditor = usePanelEditor({});
  const { refresh } = useRefreshNotes();

  const openNewNoteModal = () => {
    openConfirmModal({
      title: "Create a new staff note",
      children: (
        <Paper withBorder p="xs">
          <EditorContent editor={newNoteEditor} />
        </Paper>
      ),
      labels: { confirm: "Create", cancel: "Cancel" },
      onCancel: () => {
        /*empty*/
      },
      onConfirm: async () => {
        await createNote(newNoteEditor?.getJSON() ?? {});
        await refresh();
      },
    });
  };

  return (
    <Paper withBorder radius="md" p="sm">
      <ScrollArea mah={"70vh"}>
        <Stack>
          {notes.map(n => (
            <NoteEntry key={n.id} note={n} selected={selNote === n.id} />
          ))}
          <Button onClick={openNewNoteModal} leftSection={<PlusIcon />}>
            Add
          </Button>
        </Stack>
      </ScrollArea>
    </Paper>
  );
};
