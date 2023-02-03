import { Center, Loader, Stack, Text } from "@mantine/core";

export const LoadingSpinner = () => {
  return (
    <Center>
      <Stack spacing={"sm"}>
        <Center>
          <Loader />
        </Center>
        <Text>Loading...</Text>
      </Stack>
    </Center>
  );
};
