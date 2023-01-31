import { Center, Container, LoadingOverlay, Tabs } from "@mantine/core";
import { InfoIcon } from "@primer/octicons-react";
import { BankInfo } from "@src/components/Characters/Bank";
import { CharacterInfo } from "@src/components/Characters/Info";
import { SearchAndSelect } from "@src/components/Characters/SearchAndSelect";
import { VehicleInfo } from "@src/components/Characters/Vehicles";
import { FontAwesomeIcon } from "@src/components/Icon";
import { useCfxActions } from "@src/stores/cfx/useCfxActions";
import { characterState } from "@src/stores/character/state";
import { useEffect, useState } from "react";
import { Navigate, useParams } from "react-router-dom";
import { useRecoilState } from "recoil";

export const CharacterPage = () => {
  const { validateCid } = useCfxActions();
  const { cid } = useParams();
  const [selectedCid, setSelectedCid] = useRecoilState(characterState.cid);
  // 0 =loading, 1 = succ, 2 = failed
  const [cidState, setCidState] = useState(0);

  useEffect(() => {
    if (cid && selectedCid !== Number(cid)) {
      setSelectedCid(Number(cid));
    }
  });

  useEffect(() => {
    let ignore = false;
    const fetch = async () => {
      const valid = await validateCid(parseInt(cid ?? "0"));
      if (ignore) return;
      setCidState(valid ? 1 : 2);
    };
    fetch();
    return () => {
      ignore = true;
    };
  }, []);

  if (!cid || !parseInt(cid)) {
    return <Navigate to="/errors/404" replace />;
  }
  if (cidState === 0) {
    return <LoadingOverlay visible overlayBlur={7} />;
  }
  if (cidState === 2) {
    return <Navigate to="/errors/404" replace />;
  }
  return (
    <Container>
      <Center>
        <SearchAndSelect cid={cid} />
      </Center>
      <Tabs variant="pills" defaultValue={"info"} keepMounted={false} pt={"sm"}>
        <Tabs.List mb={"xs"}>
          <Tabs.Tab value="info" icon={<InfoIcon size={14} />}>
            Info
          </Tabs.Tab>
          <Tabs.Tab value="bank" icon={<FontAwesomeIcon icon={"building-columns"} size={"sm"} />}>
            Bank
          </Tabs.Tab>
          <Tabs.Tab value="vehicles" icon={<FontAwesomeIcon icon={"car"} size={"sm"} />}>
            Vehicles
          </Tabs.Tab>
          <Tabs.Tab value="inventory" icon={<FontAwesomeIcon icon={"backpack"} size={"sm"} />}>
            Inventory
          </Tabs.Tab>
        </Tabs.List>

        <Tabs.Panel value="info">
          <CharacterInfo cid={Number(cid)} />
        </Tabs.Panel>

        <Tabs.Panel value="bank">
          <BankInfo cid={Number(cid)} />
        </Tabs.Panel>

        <Tabs.Panel value="vehicles">
          <VehicleInfo cid={Number(cid)} />
        </Tabs.Panel>
      </Tabs>
    </Container>
  );
};
