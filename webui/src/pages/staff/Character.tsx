import { Center, Container, LoadingOverlay, Tabs } from "@mantine/core";
import { InfoIcon } from "@primer/octicons-react";
import { BankInfo } from "@src/components/Characters/Bank";
import { CharacterInfo } from "@src/components/Characters/Info";
import { RealEstateList } from "@src/components/Characters/RealEstate/RealEstateList";
import { SelectCharacter } from "@src/components/Characters/Select";
import { VehicleInfo } from "@src/components/Characters/Vehicles";
import { FontAwesomeIcon } from "@src/components/Icon";
import { useCfxPlayer } from "@src/stores/cfx/hooks/useCfxPlayer";
import { characterState } from "@src/stores/character/state";
import { useCharacterActions } from "@src/stores/character/useCharacterActions";
import { useEffect, useState } from "react";
import { Navigate, useNavigate, useParams } from "react-router-dom";
import { useRecoilState } from "recoil";

export const CharacterPage = () => {
  const { validateCid } = useCfxPlayer();
  const { cid } = useParams();
  const [selectedCid, setSelectedCid] = useRecoilState(characterState.cid);
  // 0 =loading, 1 = succ, 2 = failed
  const [cidState, setCidState] = useState(0);
  const { resetStores } = useCharacterActions();
  const navigate = useNavigate();

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
    return <LoadingOverlay visible overlayProps={{ blur: 7 }} />;
  }
  if (cidState === 2) {
    return <Navigate to="/errors/404" replace />;
  }
  return (
    <Container>
      <Center>
        <SelectCharacter
          cid={cid}
          onChange={cid => {
            resetStores();
            setSelectedCid(cid);
            navigate(`/staff/characters/${cid}`);
          }}
        />
      </Center>
      <Tabs variant="pills" defaultValue={"info"} keepMounted={false} pt={"sm"}>
        <Tabs.List mb={"xs"}>
          <Tabs.Tab value="info" leftSection={<InfoIcon size={14} />}>
            Info
          </Tabs.Tab>
          <Tabs.Tab value="bank" leftSection={<FontAwesomeIcon icon={"building-columns"} size={"sm"} />}>
            Bank
          </Tabs.Tab>
          <Tabs.Tab value="vehicles" leftSection={<FontAwesomeIcon icon={"car"} size={"sm"} />}>
            Vehicles
          </Tabs.Tab>
          <Tabs.Tab value="inventory" leftSection={<FontAwesomeIcon icon={"backpack"} size={"sm"} />}>
            Inventory
          </Tabs.Tab>
          <Tabs.Tab value="real_estate" leftSection={<FontAwesomeIcon icon={"house"} size={"sm"} />}>
            Real Estate
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

        <Tabs.Panel value="real_estate">
          <RealEstateList cid={Number(cid)} />
        </Tabs.Panel>
      </Tabs>
    </Container>
  );
};
