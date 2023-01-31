import { Button, Center, Container, Loader, Menu, Table, Text } from "@mantine/core";
import { ArrowSwitchIcon } from "@primer/octicons-react";
import { axiosInstance } from "@src/helpers/axiosInstance";
import { characterState } from "@src/stores/character/state";
import { FC, useEffect } from "react";
import { useRecoilState } from "recoil";

export const VehicleInfo: FC<{ cid: number }> = ({ cid }) => {
  const [charVehicles, setCharVehicles] = useRecoilState(characterState.vehicles);

  const fetchData = async () => {
    try {
      const res = await axiosInstance.get<CharacterState.Vehicle[]>(`/character/vehicles/cid/${cid}`);
      if (res.status !== 200) return;
      return res.data;
    } catch (e) {
      console.error(e);
    }
  };

  useEffect(() => {
    let ignore = false;
    if (!charVehicles) {
      const doFetch = async () => {
        const data = await fetchData();
        if (data && !ignore) {
          setCharVehicles(data);
        }
      };
      doFetch();
    }
    return () => {
      ignore = true;
    };
  }, []);

  if (!charVehicles) {
    return (
      <div>
        <Loader size="lg" />
        <Text>Loading...</Text>
      </div>
    );
  }

  if (charVehicles.length === 0) {
    return (
      <Container>
        <Center>
          <Text>Deze persoon bezit geen voertuigen</Text>
        </Center>
      </Container>
    );
  }

  // List including some quick actions:
  //  - transfer ownership
  //  - Change State
  //  - Repair
  //  - Delete if exists
  return (
    <Container>
      <Table>
        <thead>
          <tr>
            <th>Model</th>
            <th>Vin</th>
            <th>Plate</th>
            <th>Fake plate</th>
            <th>State</th>
            <th>Garage Id</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {charVehicles.map(veh => (
            <tr key={veh.vin}>
              <td>{veh.model}</td>
              <td>{veh.vin}</td>
              <td>{veh.plate}</td>
              <td>{veh.fakeplate}</td>
              <td>{veh.state}</td>
              <td>{veh.garageId}</td>
              <td>
                <Menu>
                  <Menu.Target>
                    <Button>Actions</Button>
                  </Menu.Target>
                  <Menu.Dropdown>
                    <Menu.Item icon={<ArrowSwitchIcon size={14} />}>Transfer Ownership</Menu.Item>
                  </Menu.Dropdown>
                </Menu>
              </td>
            </tr>
          ))}
        </tbody>
      </Table>
    </Container>
  );
};
