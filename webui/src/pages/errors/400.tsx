import { Centerbox } from "@components/CenterBox/centerbox";
import { Text } from "@mantine/core";

export const E403 = () => {
  return (
    <Centerbox title={"403"}>
      <Text>Well, it seems like you don&apos;t have access to this page</Text>
    </Centerbox>
  );
};

export const E404 = () => {
  return (
    <Centerbox title={"404"}>
      <Text>We could not find what you were searching for</Text>
    </Centerbox>
  );
};
