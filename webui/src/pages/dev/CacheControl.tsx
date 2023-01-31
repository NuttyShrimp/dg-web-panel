import { Button, Container, Flex, Text, Title } from "@mantine/core";
import { AlertIcon } from "@primer/octicons-react";
import { cacheControlEntries } from "@src/data/cacheControl";
import { axiosInstance } from "@src/helpers/axiosInstance";

export const CacheControl = () => {
  const handleCacheReset = (endpoint: string) => {
    axiosInstance.post(endpoint);
  };

  return (
    <Container>
      <Flex justify="center" align="center" direction="column">
        <Title size="h4">Control the server caches</Title>
        <Text>
          <AlertIcon /> Think twice before resetting a cache. It can lead to the database or cfx server hanging stuck
          because of the amount of data that is being transfered
        </Text>
      </Flex>

      <Flex mih={50} gap="md" justify="flex-start" align="flex-start" direction="row" wrap="wrap">
        {cacheControlEntries.map((entry, idx) => (
          <Button key={`cache-control-btn-${idx}`} onClick={() => handleCacheReset(entry.endpoint)}>
            {entry.label}
          </Button>
        ))}
      </Flex>
    </Container>
  );
};
