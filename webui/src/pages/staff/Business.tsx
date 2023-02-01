import { Container, Tabs, Title } from "@mantine/core";
import { InfoIcon } from "@primer/octicons-react";
import { BusinessEmployees } from "@src/components/Business/Employees";
import { BusinessInfo } from "@src/components/Business/Info";
import { BusinessLogs } from "@src/components/Business/Logs";
import { FontAwesomeIcon } from "@src/components/Icon";
import { useCfxBusiness } from "@src/stores/cfx/hooks/useCfxBusiness";
import { Link, useParams } from "react-router-dom";

export const Business = () => {
  const { id } = useParams();
  const { getInfo } = useCfxBusiness();

  const business = getInfo(id ?? "0");

  if (id === "" || !business) {
    <Link to={"/errors/404"} />;
  }

  return (
    <Container>
      <Title order={2}>
        {business?.label} | {business?.id}
      </Title>
      <Tabs variant="pills" defaultValue={"info"} keepMounted={false} pt={"sm"}>
        <Tabs.List mb={"xs"}>
          <Tabs.Tab value="info" icon={<InfoIcon size={14} />}>
            Info
          </Tabs.Tab>
          <Tabs.Tab value="employees" icon={<FontAwesomeIcon icon={"users"} size={"sm"} />}>
            Employees
          </Tabs.Tab>
          <Tabs.Tab value="logs" icon={<FontAwesomeIcon icon={"book"} size={"sm"} />}>
            Logs
          </Tabs.Tab>
        </Tabs.List>

        <Tabs.Panel value="info">{business && <BusinessInfo info={business} />}</Tabs.Panel>
        <Tabs.Panel value="employees">
          <BusinessEmployees id={Number(id)} />
        </Tabs.Panel>
        <Tabs.Panel value="logs">
          <BusinessLogs id={Number(id)} />
        </Tabs.Panel>
      </Tabs>
    </Container>
  );
};
