import { Dispatch, SetStateAction } from "react";
import { axiosPost, CommonResp } from "./basic.ts";

export type LoginReq = {
  username: string;
  password: string;
  clientId: string;
  clientSecret: string;
  codeChallenge: string;
  redirect: string;
};

export type LoginResp = { authorizationCode: string } & CommonResp;

export async function loginApi(
  req: LoginReq,
  setWarn: Dispatch<SetStateAction<string>>,
) {
  return axiosPost<LoginResp>("/auth/login", req, setWarn);
}
