import { useSearchParams } from "react-router-dom";
import { useRef, useState } from "react";
import { Button } from "./components/Button.tsx";
import { Input } from "./components/Input.tsx";
import { registerApi, RegisterReq } from "./apis/setup.ts";

export function Register() {
  const [info, setInfo] = useState("");
  const [warn, setWarn] = useState("");
  const usernameRef = useRef<HTMLInputElement | null>(null);
  const passwordRef = useRef<HTMLInputElement | null>(null);
  const [searchParams] = useSearchParams();

  function register() {
    const invitationCode = searchParams.get("invitationCode");
    if (
      !(
        invitationCode &&
        usernameRef.current &&
        usernameRef.current.value &&
        passwordRef.current &&
        passwordRef.current.value
      )
    ) {
      setWarn("Some information is missing");
      return;
    }
    const req: RegisterReq = {
      username: usernameRef.current.value,
      password: passwordRef.current.value,
      invitationCode: invitationCode,
    };
    registerApi(req, setWarn).then((resp) => {
      if (resp) {
        setInfo("Register succeeded! You can now close the window.");
      }
    });
  }

  return (
    <div className={"card"}>
      <h1>Register MyOAuth</h1>
      {warn && <span className={"warn"}>{warn}</span>}
      {info && <span className={"warn black"}>{info}</span>}
      {!info && (
        <>
          <Input
            label={"Username"}
            inputType={"text"}
            myRef={usernameRef}
            enter={register}
          />
          <Input
            label={"Password"}
            inputType={"password"}
            myRef={passwordRef}
            enter={register}
          />
          <Button text={"Go"} click={register} className={"full-width mt-1"} />
        </>
      )}
    </div>
  );
}
