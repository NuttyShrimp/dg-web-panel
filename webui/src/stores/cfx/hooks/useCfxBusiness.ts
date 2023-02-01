import { axiosInstance } from "@src/helpers/axiosInstance";
import { useCallback } from "react";
import { useRecoilState } from "recoil";
import { cfxState } from "../state";

export const useCfxBusiness = () => {
  const [businesses, setBusinesses] = useRecoilState(cfxState.businesses);

  const fetchAll = async (filter?: { cid: string }) => {
    try {
      const urlParams = new URLSearchParams();
      if (filter) {
        Object.keys(filter).forEach(fk => {
          urlParams.set(fk, filter[fk as keyof typeof filter]);
        });
      }
      const res = await axiosInstance.get<CfxState.Business.Entry[]>(`/staff/business/all?${urlParams.toString()}`);
      if (res.status !== 200) return;
      setBusinesses(res.data);
    } catch (e) {
      console.error(e);
    }
  };

  const getInfo = useCallback(
    (id: string) => {
      const business = businesses.find(b => b.id === Number(id));
      if (!business) fetchAll();
      return business;
    },
    [businesses]
  );

  const fetchLogs = async (id: number, page: number) => {
    try {
      const res = await axiosInstance.get<CfxState.Business.Log[]>(`/staff/business/${id}/logs?page=${page}`);
      return res.data;
    } catch (e) {
      console.error(e);
      return [];
    }
  };

  const deleteBusiness = async (id: number) => {
    try {
      await axiosInstance.delete(`/staff/business/${id}`);
      setBusinesses(bs => bs.filter(b => b.id !== id));
    } catch (e) {
      console.error(e);
    }
  };

  const changeOwner = async (id: number, cid: number) => {
    try {
      await axiosInstance.post(`/staff/business/${id}/owner`, {
        cid,
      });
    } catch (e) {
      console.error(e);
    }
  };

  return {
    fetchAll,
    fetchLogs,
    getInfo,
    deleteBusiness,
    changeOwner,
  };
};
