import axios from "axios";

const URL = "http://localhost:20476";

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

export async function clientList() {
  return axios.get(`${URL}/setup/client-list`);
}

export type ClientCreateReq = {
  adminPassword: string;
} & Client;

export async function clientCreate(client: ClientCreateReq) {
  return axios.post(`${URL}/setup/client-create`, client);
}

export async function clientDelete({
  id,
  adminPassword,
}: {
  id: string;
  adminPassword: string;
}) {}

export async function login() {}
