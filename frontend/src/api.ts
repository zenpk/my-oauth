import axios from "axios";

export type Client = {
  id: string;
  secret: string;
  redirects: string;
  accessTokenAge: number;
  refreshTokenAge: number;
};

const URL = "http://localhost:20476";

export async function clientList() {
  return axios.get(`${URL}/setup/client-list`);
}

export async function clientCreate() {}

export async function clientDelete({
  id,
  adminPassword,
}: {
  id: string;
  adminPassword: string;
}) {}

export async function login() {}
