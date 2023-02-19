import { Alert } from "@mantine/core";
import { AlertFillIcon } from "@primer/octicons-react";
import { basicGet } from "@src/lib/actions/basicReq";
import { useQuery } from "@tanstack/react-query";
import { FC, PropsWithChildren } from "react";

export const UpdateAnnounceProvider: FC<PropsWithChildren<{}>> = props => {
  const { isSuccess, data } = useQuery<boolean, Error>({
    queryKey: ["update-state"],
    queryFn: () => basicGet<boolean>("/state/schedule/update"),
  });

  return (
    <>
      {isSuccess && data && (
        <Alert
          color="dg-prim"
          title="Scheduled update"
          variant="filled"
          icon={<AlertFillIcon size={14} />}
          style={{
            position: "absolute",
            top: "1rem",
            left: 0,
            right: 0,
            margin: "0 auto",
            width: "400px",
          }}
        >
          Pasop! Binnen minder dan 1 minuut zal er een update naar de website gepushed worden. Om het verlies van werk
          tegen te gaan raden we je aan om je werk op te slaan!
        </Alert>
      )}
      {props.children}
    </>
  );
};
