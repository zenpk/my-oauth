import { Input } from "./components/Input.tsx";
import { useRef, useState } from "react";
import { Button } from "./components/Button.tsx";
import { useSearchParams } from "react-router-dom";
import { loginApi, LoginReq } from "./apis/auth.ts";

export function Login() {
  const [warn, setWarn] = useState("");
  const usernameRef = useRef<HTMLInputElement | null>(null);
  const passwordRef = useRef<HTMLInputElement | null>(null);
  const [searchParams] = useSearchParams();

  function login() {
    const clientId = searchParams.get("clientId");
    const codeChallenge = searchParams.get("codeChallenge");
    const redirect = searchParams.get("redirect");
    if (
      !(
        clientId &&
        codeChallenge &&
        redirect &&
        usernameRef.current &&
        usernameRef.current.value &&
        passwordRef.current &&
        passwordRef.current.value
      )
    ) {
      setWarn("Some information is missing");
      return;
    }
    const req: LoginReq = {
      username: usernameRef.current.value,
      password: passwordRef.current.value,
      clientId: decodeURIComponent(clientId),
      codeChallenge: decodeURIComponent(codeChallenge),
      redirect: decodeURIComponent(redirect),
    };
    loginApi(req, setWarn).then((resp) => {
      if (resp) {
        window.location.replace(
          `${decodeURIComponent(redirect)}?authorizationCode=${
            resp.authorizationCode
          }`,
        );
      }
    });
  }

  return (
    <div className={"card"}>
      <h1>Login with MyOAuth</h1>
      {warn && <span className={"warn"}>{warn}</span>}
      <Input
        label={"Username"}
        inputType={"text"}
        myRef={usernameRef}
        enter={login}
      />
      <Input
        label={"Password"}
        inputType={"password"}
        myRef={passwordRef}
        enter={login}
      />
      <Button text={"Go"} click={login} className={"full-width mt-1"} />
    </div>
  );
}
