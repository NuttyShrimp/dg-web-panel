import { useQuery } from "@tanstack/react-query";
import { basicGet } from "@src/lib/actions/basicReq";
import { Button, Card, Divider, Flex, Image, Text, Title } from "@mantine/core";
import { Link } from "@src/components/Router/Link";
import { LoadingSpinner } from "@src/components/LoadingSpinner";
import { useMemo } from "react";
import { queryClient } from "@src/helpers/queryClient";
import { axiosInstance } from "@src/helpers/axiosInstance";

const FlyerCard = ({ flyer, canApprove }: { flyer: Dev.Flyer; canApprove?: boolean }) => (
  <Card shadow="sm" m="xs" p="lg" radius="md" withBorder>
    <Card.Section>
      <Image src={flyer.link} height={200} alt={flyer.link} fit="contain" />
    </Card.Section>

    <Text size="sm" mt={"sm"}>
      <Link to={`/staff/users/${flyer.character.steamId}`} target="_blank">
        {flyer.character.user.name}
      </Link>
      {" - "}
      <Link to={`/staff/characters/${flyer.character.citizenid}`} target="_blank">
        {flyer.character.info.firstname} {flyer.character.info.lastname} ({flyer.character.citizenid})
      </Link>
    </Text>

    <Flex justify={"space-between"}>
      {canApprove && (
        <Button
          variant="light"
          color="green"
          mt="md"
          radius="md"
          onClick={async () => {
            await axiosInstance.post(`/cfx/flyers/${flyer.id}`);
            queryClient.invalidateQueries(["flyers"]);
          }}
        >
          Approve
        </Button>
      )}
      <Button
        variant="light"
        color="red"
        mt="md"
        radius="md"
        onClick={async () => {
          await axiosInstance.delete(`/cfx/flyers/${flyer.id}`);
          queryClient.invalidateQueries(["flyers"]);
        }}
      >
        Remove
      </Button>
    </Flex>
  </Card>
);

const RetrievedFlyers = () => {
  const {
    data: items,
    isLoading,
    isError,
    error,
  } = useQuery<Inventory.Item[], Error>({
    queryKey: ["flyers", "retrieved"],
    queryFn: () => basicGet<Inventory.Item[]>("/cfx/flyers/retrieved"),
  });

  if (isError) {
    return (
      <div>
        <h3>Error</h3>
        <div>An unexpected error occurred {error.message}</div>
        {import.meta.env.DEV && <pre>{JSON.stringify(error, null, 2)}</pre>}
      </div>
    );
  }

  return (
    <div>
      <Title order={2} mb={"xs"}>
        Existing Flyers
      </Title>
      {isLoading && <LoadingSpinner />}
      <Flex wrap={"wrap"}>
        {items?.map(item => (
          <Card shadow="sm" m="xs" p="lg" radius="md" withBorder key={item.id}>
            <Card.Section>
              <Image src={item.metadata.link} height={200} alt={item.metadata.link} fit="contain" />
            </Card.Section>

            <Text size="sm" mt={"sm"}>
              {item.inventory}
            </Text>

            <Flex justify={"space-between"}>
              <Button
                variant="light"
                color="red"
                mt="md"
                radius="md"
                onClick={async () => {
                  await axiosInstance.delete(`/cfx/inventory/${item.id}`);
                  queryClient.invalidateQueries(["flyers", "retrieved"]);
                }}
              >
                Remove
              </Button>
            </Flex>
          </Card>
        ))}
      </Flex>
    </div>
  );
};

export const FlyerRequestPage = () => {
  const {
    data: flyers,
    isLoading,
    isError,
    error,
  } = useQuery<Dev.Flyer[], Error>({
    queryKey: ["flyers"],
    queryFn: () => basicGet<Dev.Flyer[]>("/cfx/flyers/"),
  });

  const approvedFlyers = useMemo(() => flyers?.filter(flyer => flyer.approved), [flyers]);
  const pendingFlyers = useMemo(() => flyers?.filter(flyer => !flyer.approved), [flyers]);

  if (isError) {
    return (
      <div>
        <h3>Error</h3>
        <div>An unexpected error occurred {error.message}</div>
        {import.meta.env.DEV && <pre>{JSON.stringify(error, null, 2)}</pre>}
      </div>
    );
  }

  if (isLoading) {
    return <LoadingSpinner />;
  }

  return (
    <div>
      <Title order={2} mb={"xs"}>
        Pending Approval
      </Title>
      <Flex wrap={"wrap"}>
        {pendingFlyers?.map(flyer => (
          <FlyerCard flyer={flyer} key={flyer.id} canApprove />
        ))}
      </Flex>
      <Divider mt={"xs"} />
      <Title order={2} mb={"xs"}>
        Approved
      </Title>
      <Flex wrap={"wrap"}>
        {approvedFlyers?.map(flyer => (
          <FlyerCard flyer={flyer} key={flyer.id} />
        ))}
      </Flex>
      <Divider mt={"xs"} />
      <RetrievedFlyers />
    </div>
  );
};
