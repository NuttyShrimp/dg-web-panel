import { Container, Text } from "@mantine/core";
import { useParams } from "react-router-dom";

export const UserPage = () => {
  const { steamid } = useParams();
  return (
    <Container>
      <Text>{steamid}</Text>
    </Container>
  );
};
