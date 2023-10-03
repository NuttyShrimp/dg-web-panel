import { Center, Container, Stack, Text } from "@mantine/core";
import { List } from "@src/components/List";
import { Link } from "@src/components/Router/Link";
import { useCfxBusiness } from "@src/stores/cfx/hooks/useCfxBusiness";
import { cfxState } from "@src/stores/cfx/state";
import { useEffect } from "react";
import { useRecoilValue } from "recoil";

export const BusinessList = () => {
  const businesses = useRecoilValue(cfxState.businesses);
  const { fetchAll } = useCfxBusiness();

  useEffect(() => {
    fetchAll();
  }, []);

  return (
    <Container>
      <Center>
        <List highlightHover>
          {businesses.map(b => (
            <List.Entry key={b.id}>
              <Link to={`/staff/business/${b.id}`} style={{ width: "100%" }} noColor>
                <Stack gap={5}>
                  <Text fw={"bolder"}>
                    {b.label} | {b.name} | {b.id}
                  </Text>
                  <Text size={"sm"}>type: {b.type.name}</Text>
                </Stack>
              </Link>
            </List.Entry>
          ))}
        </List>
      </Center>
    </Container>
  );
};
