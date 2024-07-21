import axios, { AxiosResponse } from "axios";

export type ChallengeVerifier = {
  codeChallenge: string;
  codeVerifier: string;
};

export type LoginReq = {
  clientId: string;
  redirect: string;
  codeChallenge: string;
};

export type AuthorizeReq = {
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
} & CommonResp;

export type RefreshReq = {
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
  constructor(
    private oAuthEndpoint: string,
    private authorizeEndpoint: string,
    private refreshEndpoint: string
  ) {}

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
    window.location.replace(
      `${this.oAuthEndpoint}/login?clientId=${clientId}&codeChallenge=${codeChallenge}&redirect=${redirect}`
    );
  }

  public authorize(
    codeVerifier: string
  ): Promise<AxiosResponse<AuthorizeResp>> {
    const urlParams = new URLSearchParams(window.location.search);
    const req: AuthorizeReq = {
      codeVerifier: codeVerifier,
      authorizationCode: urlParams.get("authorizationCode") ?? "",
    };
    return axios.post(this.authorizeEndpoint, req);
  }

  public refresh(refreshToken: string): Promise<AxiosResponse<RefreshResp>> {
    const req: RefreshReq = {
      refreshToken: refreshToken,
    };
    return axios.post(this.refreshEndpoint, req);
  }

  public verify(accessToken: string): Promise<AxiosResponse<CommonResp>> {
    return axios.post(`${this.oAuthEndpoint}/api/auth/verify`, {
      accessToken: accessToken,
    });
  }

  public getPublicKey(): Promise<AxiosResponse<PublicJwk>> {
    return axios.get(`${this.oAuthEndpoint}/api/setup/public-key`);
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
