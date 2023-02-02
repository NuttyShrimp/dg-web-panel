import { cfxState } from "@src/stores/cfx/state";
import { FC } from "react";
import { useNavigate } from "react-router-dom";
import { useSetRecoilState } from "recoil";
import { UserSelect } from "../Select";

export const SearchAndSelect: FC<{ cid?: string }> = ({ cid }) => {
  const selectPlayer = useSetRecoilState(cfxState.player);
  const navigate = useNavigate();

  return (
    <UserSelect
      steamId={cid}
      onChange={val => {
        selectPlayer(val);
        navigate(`/staff/users/${val}`);
      }}
    />
  );
};
