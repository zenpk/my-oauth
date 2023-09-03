import { Dispatch, SetStateAction } from "react";
import { axiosDelete, axiosGet, axiosPost, CommonResp } from "./basic.ts";

export type Client = {
  id: string;
  secret: string;
  redirects: string;
  accessTokenAge: number;
  refreshTokenAge: number;
};

export type ClientListResp = { clients: Client[] } & CommonResp;

export async function clientListApi(setWarn: Dispatch<SetStateAction<string>>) {
  return axiosGet<ClientListResp>("/setup/client-list", setWarn);
}

export type ClientCreateReq = {
  adminPassword: string;
} & Client;

export async function clientCreateApi(
  req: ClientCreateReq,
  setWarn: Dispatch<SetStateAction<string>>,
) {
  return axiosPost<CommonResp>("/setup/client-create", req, setWarn);
}

export type ClientDeleteReq = {
  id: string;
  adminPassword: string;
};

export async function clientDeleteApi(
  req: ClientDeleteReq,
  setWarn: Dispatch<SetStateAction<string>>,
) {
  return axiosDelete<CommonResp>("/setup/client-delete", req, setWarn);
}

export type RegisterReq = {
  username: string;
  password: string;
  invitationCode: string;
};

export async function registerApi(
  req: RegisterReq,
  setWarn: Dispatch<SetStateAction<string>>,
) {
  return axiosPost<CommonResp>("/setup/register", req, setWarn);
}
