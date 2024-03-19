import axios from "axios";
import { Dispatch, SetStateAction } from "react";

const BASE = import.meta.env.VITE_BACKEND_URL;

export type CommonResp = {
  ok: boolean;
  msg: string;
};
export async function axiosGet<T extends CommonResp>(
  url: string,
  setWarn: Dispatch<SetStateAction<string>>,
) {
  try {
    const resp = await axios.get(`${BASE}${url}`);
    const data = resp.data as T;
    if (!data.ok) {
      setWarn(data.msg);
      return null;
    } else {
      setWarn("");
      return data;
    }
  } catch (err: any) {
    setWarn(err.toString());
    return null;
  }
}

export async function axiosPost<T extends CommonResp>(
  url: string,
  body: any,
  setWarn: Dispatch<SetStateAction<string>>,
) {
  try {
    const resp = await axios.post(`${BASE}${url}`, body);
    const data = resp.data as T;
    if (!data.ok) {
      setWarn(data.msg);
      return null;
    } else {
      setWarn("");
      return data;
    }
  } catch (err: any) {
    setWarn(err.toString());
    return null;
  }
}

export async function axiosDelete<T extends CommonResp>(
  url: string,
  body: any,
  setWarn: Dispatch<SetStateAction<string>>,
) {
  try {
    const resp = await axios.delete(`${BASE}${url}`, { data: body });
    const data = resp.data as T;
    if (!data.ok) {
      setWarn(data.msg);
      return null;
    } else {
      setWarn("");
      return data;
    }
  } catch (err: any) {
    setWarn(err.toString());
    return null;
  }
}
