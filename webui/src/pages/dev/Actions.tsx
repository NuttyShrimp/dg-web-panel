import { Button, Container } from "@mantine/core";
import { openModal } from "@mantine/modals";
import { CreateBusinessModal } from "@src/components/Business/modals/createBusiness";
import { CreateVehiclemodal } from "@src/components/Vehicles/modals/createvehicle";

export const DevActionPage = () => {
  const createBusiness = () => {
    openModal({
      title: "Create Business",
      children: <CreateBusinessModal />,
    });
  };

  // TODO: cleanup + typos
  const giveVehicle = () => {
    openModal({
      title: "Give Vehicle",
      children: <CreateVehiclemodal />,
    });
  };

  return (
    <Container>
      <Button onClick={createBusiness}>Create business</Button>
      <Button onClick={giveVehicle}>Give Vehicle</Button>
    </Container>
  );
};
