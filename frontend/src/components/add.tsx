import {
  TextInput,
  TextInputProps,
  ActionIcon,
  useMantineTheme,
} from "@mantine/core";
import { IconCheck, IconPlus } from "@tabler/icons-react";
import { useState } from "react";

interface AddProps extends TextInputProps {
  onAdd: (title: string) => void;
}

export function Add(props: AddProps) {
  const theme = useMantineTheme();
  const [value, setValue] = useState("");
  const { onAdd } = props;

  return (
    <TextInput
      icon={<IconCheck size="1.1rem" stroke={1.5} />}
      radius="xl"
      size="md"
      onChange={(event) => {
        setValue(event.currentTarget.value);
      }}
      rightSection={
        <ActionIcon
          onClick={() => onAdd(value)}
          size={32}
          radius="xl"
          color={theme.primaryColor}
          variant="filled"
        >
          <IconPlus size="1.1rem" stroke={1.5} />
        </ActionIcon>
      }
      placeholder="Todo Title"
      rightSectionWidth={42}
    />
  );
}
