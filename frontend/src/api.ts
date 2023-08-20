import axios from "axios";
import { Dispatch, SetStateAction } from "react";

const BASE = "http://localhost:20476";

export type Client = {
  id: string;
  secret: string;
  redirects: string;
  accessTokenAge: number;
  refreshTokenAge: number;
};

export type CommonResp = {
  ok: boolean;
  msg: string;
};

export type ClientListResp = { clients: Client[] } & CommonResp;

export async function clientList(setWarn: Dispatch<SetStateAction<string>>) {
  return axiosGet<ClientListResp>("/setup/client-list", setWarn);
}

export type ClientCreateReq = {
  adminPassword: string;
} & Client;

export async function clientCreate(
  req: ClientCreateReq,
  setWarn: Dispatch<SetStateAction<string>>,
) {
  return axiosPost<CommonResp>("/setup/client-create", req, setWarn);
}

export type ClientDeleteReq = {
  id: string;
  adminPassword: string;
};

export async function clientDelete({
  req,
  setWarn,
}: {
  req: ClientDeleteReq;
  setWarn: Dispatch<SetStateAction<string>>;
}) {
  return axiosDelete<CommonResp>("/setup/client-delete", req, setWarn);
}

// export async function login() {}

type Response<T> = {
  [K: string]: T;
} & CommonResp;

async function axiosGet<T>(
  url: string,
  setWarn: Dispatch<SetStateAction<string>>,
) {
  try {
    const resp = await axios.get(`${BASE}${url}`);
    const data = resp.data as Response<T>;
    if (!data.ok) {
      setWarn(data.msg);
      return null;
    } else {
      return data;
    }
  } catch (err: any) {
    setWarn(err.toString());
  }
}

async function axiosPost<T>(
  url: string,
  body: any,
  setWarn: Dispatch<SetStateAction<string>>,
) {
  try {
    const resp = await axios.post(`${BASE}${url}`, body);
    const data = resp.data as Response<T>;
    if (!data.ok) {
      setWarn(data.msg);
      return null;
    } else {
      return data;
    }
  } catch (err: any) {
    setWarn(err.toString());
  }
}

async function axiosDelete<T>(
  url: string,
  body: any,
  setWarn: Dispatch<SetStateAction<string>>,
) {
  try {
    const resp = await axios.delete(`${BASE}${url}`, body);
    const data = resp.data as Response<T>;
    if (!data.ok) {
      setWarn(data.msg);
      return null;
    } else {
      return data;
    }
  } catch (err: any) {
    setWarn(err.toString());
  }
}
