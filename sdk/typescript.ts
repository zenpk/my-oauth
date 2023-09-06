import axios from "axios";


export type ChallengeVerifier = {
    codeChallenge: string;
    codeVerifier: string;
}

export type LoginInfo = {
    clientId: string;
    clientSecret: string;
    redirect: string;
    codeChallenge: string;
}

export class MyOAuthSdk {
    endpoint: string;

    constructor(endpoint: string) {
        this.endpoint = endpoint;
    }

    async genChallengeVerifier(len: number) {
        const bytes = new Uint8Array(len);
        for (let i = 0; i < len; i++) {
            bytes[i] = Math.floor(Math.random() * 256);
        }

        const challenge = btoa(String.fromCharCode.apply(null, bytes))
            .replace('+', '-')
            .replace('/', '_')
            .replace(/=+$/, '');

        const hashBuffer = await crypto.subtle.digest("SHA-256", bytes);
        const hashArray = new Uint8Array(hashBuffer);
        const verifier = btoa(String.fromCharCode.apply(null, hashArray))
            .replace('+', '-')
            .replace('/', '_')
            .replace(/=+$/, '');

        return {string1: challenge, string2: verifier};
    }

    redirectLogin(info: LoginInfo) {
        const clientId = this.urlEncode(info.clientId);
        const clientSecret = this.urlEncode(info.clientSecret);
        const redirect = this.urlEncode(info.redirect);
        const codeChallenge = this.urlEncode(info.codeChallenge);
        window.location.replace(`${this.endpoint}/login?clientId=${clientId}&cliendSecret=${clientSecret}&codeChallenge=${codeChallenge}&redirect=${redirect}`);
    }


    urlEncode(src: string) {
    }

}
