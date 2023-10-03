import { Flex, Text, Title } from "@mantine/core";
import { FC } from "react";
import { Link } from "../Router/Link";

export const BusinessInfo: FC<{ info: CfxState.Business.Entry }> = ({ info }) => {
  return (
    <div>
      <Title order={3}>Info</Title>
      <Flex>
        <Text fw={"bolder"}>Name:&nbsp;</Text>
        <Text>{info.name}</Text>
      </Flex>
      <Flex>
        <Text fw={"bolder"}>Label:&nbsp;</Text>
        <Text>{info.label}</Text>
      </Flex>
      <Flex>
        <Text fw={"bolder"}>Id:&nbsp;</Text>
        <Text>{info.id}</Text>
      </Flex>
      <Flex>
        <Text fw={"bolder"}>Type:&nbsp;</Text>
        <Text>{info.type.name}</Text>
      </Flex>
      <Flex>
        <Text fw={"bolder"}>Bank Account:&nbsp;</Text>
        <Link to={`/staff/bank/${info.bankAccountId}`}>
          <Text>{info.bankAccountId}</Text>
        </Link>
      </Flex>
    </div>
  );
};
