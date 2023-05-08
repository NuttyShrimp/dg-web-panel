import { Divider, Flex, Loader, Table, Text, Title } from "@mantine/core";
import { displayDate } from "@src/helpers/time";
import { characterState } from "@src/stores/character/state";
import { useCharacterActions } from "@src/stores/character/useCharacterActions";
import { FC, useEffect } from "react";
import { Link } from "react-router-dom";
import { useRecoilState, useRecoilValue } from "recoil";

export const CharacterInfo: FC<{ cid: number }> = ({ cid }) => {
  const characterInfo = useRecoilValue(characterState.selected);
  const [characterReputation, setCharacterReputation] = useRecoilState(characterState.reputation);
  const { fetchCharReputation } = useCharacterActions();

  useEffect(() => {
    let ignore = false;
    if (!characterReputation) {
      const doFetch = async () => {
        const data = await fetchCharReputation(cid);
        if (!ignore) {
          setCharacterReputation(data);
        }
      };
      doFetch();
    }
    return () => {
      ignore = true;
    };
  }, []);

  if (!characterInfo) {
    return (
      <div>
        <Loader size="lg" />
        <Text>Loading...</Text>
      </div>
    );
  }

  return (
    <div>
      <Title size="h2">
        {characterInfo?.info.firstname} {characterInfo?.info.lastname}
      </Title>
      <Divider mb={"xs"} />
      <Title size="h3">Character info</Title>
      <Flex>
        <Text weight={"bolder"}>SteamId:&nbsp;</Text>
        <Link to={`/staff/users/${characterInfo?.user.steamId}`}>
          <Text>{characterInfo?.user.steamId}</Text>
        </Link>
      </Flex>
      <Flex>
        <Text weight={"bolder"}>Citizenid:&nbsp;</Text>
        <Text>{characterInfo?.citizenid}</Text>
      </Flex>
      <Flex>
        <Text weight={"bolder"}>Birthdate:&nbsp;</Text>
        <Text>{characterInfo?.info.birthdate}</Text>
      </Flex>
      <Flex>
        <Text weight={"bolder"}>Phone Nr.:&nbsp;</Text>
        <Text>{characterInfo?.info.phone}</Text>
      </Flex>
      <Flex>
        <Text weight={"bolder"}>Gender:&nbsp;</Text>
        <Text>{characterInfo?.info.gender ? "Woman" : "Man"}</Text>
      </Flex>
      <Flex>
        <Text weight={"bolder"}>Created at:&nbsp;</Text>
        <Text>{displayDate(characterInfo.created_at)}</Text>
      </Flex>
      <Flex>
        <Text weight={"bolder"}>Last updated:&nbsp;</Text>
        <Text>{displayDate(characterInfo.last_updated)}</Text>
      </Flex>
      <Divider my={"sm"} />
      <Title size="h3">Metadata</Title>
      <Table>
        <thead>
          <tr>
            <th style={{ width: `30%` }}>Key</th>
            <th>Value</th>
          </tr>
        </thead>
        <tbody>
          {Object.entries(JSON.parse(characterInfo.data.metadata)).map(([k, v]) => (
            <tr key={`character-metadata-${k}`}>
              <td>{k}</td>
              <td>{JSON.stringify(v)}</td>
            </tr>
          ))}
        </tbody>
      </Table>
      {characterReputation && (
        <>
          <Divider my={"sm"} />
          <Title size="h3">Reputation</Title>
          <Table>
            <thead>
              <tr>
                <th style={{ width: `30%` }}>Place</th>
                <th>Rep</th>
              </tr>
            </thead>
            <tbody>
              {Object.entries(characterReputation).map(([k, v]) => (
                <tr key={`character-reputation-${k}`}>
                  <td>{k}</td>
                  <td>{JSON.stringify(v)}</td>
                </tr>
              ))}
            </tbody>
          </Table>
        </>
      )}
    </div>
  );
};
