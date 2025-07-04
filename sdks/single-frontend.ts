import axios, { AxiosResponse } from "axios";

export type ChallengeVerifier = {
  codeChallenge: string;
  codeVerifier: string;
};

export type LoginReq = {
  clientId: string;
  redirect: string;
  codeChallenge: string;
  context: string;
};

export type AuthorizeReq = {
  clientId: string;
  clientSecret: string;
  codeVerifier: string;
  authorizationCode: string;
};

export type CommonResp = {
  ok: boolean;
  msg: string;
};

export type AuthorizeResp = {
  accessToken: string;
  refreshToken: string;
  context: string;
} & CommonResp;

export type RefreshReq = {
  clientId: string;
  clientSecret: string;
  refreshToken: string;
};

export type RefreshResp = {
  accessToken: string;
} & CommonResp;

export type PublicJwk = {
  kty: string;
  e: string;
  use: string;
  kid: string;
  alg: string;
  n: string;
};

export class MyOAuthSdk {
  constructor(private endpoint: string) {}

  public async genChallengeVerifier(len: number) {
    const bytes = new Uint8Array(len);
    crypto.getRandomValues(bytes);

    const verifier = this.arrayToBase64Url(bytes);

    const encoder = new TextEncoder();
    const data = encoder.encode(verifier);
    const hashBuffer = await crypto.subtle.digest("SHA-256", data);
    const hashArray = new Uint8Array(hashBuffer);
    const challenge = this.arrayToBase64Url(hashArray);

    const challengeVerifier: ChallengeVerifier = {
      codeChallenge: challenge,
      codeVerifier: verifier,
    };
    return challengeVerifier;
  }

  public redirectLogin(req: LoginReq) {
    const clientId = encodeURIComponent(req.clientId);
    const redirect = encodeURIComponent(req.redirect);
    const codeChallenge = encodeURIComponent(req.codeChallenge);
    const context = encodeURIComponent(req.context);
    window.location.replace(
      `${this.endpoint}/login?clientId=${clientId}&codeChallenge=${codeChallenge}&redirect=${redirect}&context=${context}`
    );
  }

  public authorize(req: AuthorizeReq): Promise<AxiosResponse<AuthorizeResp>> {
    const urlParams = new URLSearchParams(window.location.search);
    req.authorizationCode = urlParams.get("authorizationCode") ?? "";
    return axios.post(`${this.endpoint}/api/auth/authorize`, req);
  }

  public refresh(req: RefreshReq): Promise<AxiosResponse<RefreshResp>> {
    return axios.post(`${this.endpoint}/api/auth/refresh`, req);
  }

  public verify(accessToken: string): Promise<AxiosResponse<CommonResp>> {
    return axios.post(`${this.endpoint}/api/auth/verify`, {
      accessToken: accessToken,
    });
  }

  public getPublicKey(): Promise<AxiosResponse<PublicJwk>> {
    return axios.get(`${this.endpoint}/api/setup/public-key`);
  }

  public arrayToBase64Url(array: Uint8Array) {
    let src = "";
    array.forEach((num) => {
      src += String.fromCharCode(num);
    });
    return this.stringToBase64Url(src);
  }

  public stringToBase64Url(src: string) {
    return btoa(src).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");
  }
}
