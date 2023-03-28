import { Button, TextInput } from "@mantine/core";
import { useEffect, useState } from "react";
import { SelectCharacter } from "../../Characters/Select";
import { UserSelect } from "../../Users/Select";

import "./style.scss";

interface InputsModalProps {
  fields: string[];
  onSubmit: (answers: Record<string, string>) => void;
}

const fieldToSelect: Record<string, any> = {
  citizenid: SelectCharacter,
  steamId: UserSelect,
};

export const InputsModal = (props: InputsModalProps) => {
  const [fields, setFields] = useState<Record<string, string>>({});

  useEffect(() => {
    const nFields: Record<string, string> = {};
    props.fields.forEach(f => {
      nFields[f] = "";
    });
    setFields(nFields);
  }, [props.fields]);

  const updateField = (field: string, value: string) => {
    if (fields[field] === undefined) return;
    setFields({ ...fields, [field]: value });
  };

  return (
    <div className="input-model-content">
      {props.fields.map(f =>
        fieldToSelect?.[f] ? (
          fieldToSelect[f]({
            key: f,
            value: fields[f],
            onChange: (e: any) => updateField(f, String(e)),
          })
        ) : (
          <TextInput key={f} value={fields[f]} onChange={e => updateField(f, e.currentTarget.value)} />
        )
      )}
      <Button mt={"xs"} onClick={() => props.onSubmit(fields)}>
        Submit
      </Button>
    </div>
  );
};
