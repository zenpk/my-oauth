import { Dispatch, SetStateAction, useEffect, useRef, useState } from "react";
import { Input } from "./components/Input.tsx";
import { Button } from "./components/Button.tsx";
import {
  Client,
  clientCreate,
  ClientCreateReq,
  clientList,
  ClientListResp,
} from "./api.ts";

export function Admin() {
  const [adminPassword, setAdminPassword] = useState("");
  const [warn, setWarn] = useState("");
  const [showAddForm, setShowAddForm] = useState(false);
  const [clients, setClients] = useState<Client[]>([]);
  const adminPasswordRef = useRef<HTMLInputElement | null>(null);

  useEffect(() => {
    if (adminPassword !== "" && !showAddForm) {
      clientList(setWarn).then((resp) => {
        if (resp) {
          const data = resp as ClientListResp;
          setClients(data.clients);
        }
      });
    }
  }, [adminPassword, showAddForm]);

  function saveAdminPassword() {
    if (adminPasswordRef.current && adminPasswordRef.current.value) {
      setAdminPassword(adminPasswordRef.current.value);
    } else {
      setWarn("Save admin password failed");
    }
  }

  return (
    <div className={"card"}>
      <h1>Admin Page</h1>
      {warn && <span className={"warn"}>{warn}</span>}
      {!adminPassword && (
        <div>
          <Input
            label={"Admin Password"}
            inputType={"text"}
            myRef={adminPasswordRef}
            enter={saveAdminPassword}
            buttonText={"Save"}
          />
        </div>
      )}
      {adminPassword && (
        <div className={"flex-basic-column"}>
          <h2>Client list</h2>
          {!showAddForm && (
            <Button
              text={"Add"}
              click={() => {
                setShowAddForm(true);
              }}
            />
          )}
          {showAddForm && (
            <AddForm
              setShowAddForm={setShowAddForm}
              setWarn={setWarn}
              adminPassword={adminPassword}
            />
          )}
          <table>
            <thead>
              <tr>
                <th>Id</th>
                <th>Access Token Age</th>
                <th>Refresh Token Age</th>
                <th>Redirects</th>
                <th>Action</th>
              </tr>
            </thead>
            <tbody>
              {clients.map((client) => {
                return (
                  <tr>
                    <td>{client.id}</td>
                    <td>{client.accessTokenAge}</td>
                    <td>{client.refreshTokenAge}</td>
                    <td>{client.redirects}</td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}

function AddForm({
  setShowAddForm,
  setWarn,
  adminPassword,
}: {
  setShowAddForm: Dispatch<SetStateAction<boolean>>;
  setWarn: Dispatch<SetStateAction<string>>;
  adminPassword: string;
}) {
  const idRef = useRef<HTMLInputElement | null>(null);
  const secretRef = useRef<HTMLInputElement | null>(null);
  const accessAgeRef = useRef<HTMLInputElement | null>(null);
  const refreshAgeRef = useRef<HTMLInputElement | null>(null);
  const redirectRef = useRef<HTMLInputElement | null>(null);

  function click() {
    if (
      idRef.current &&
      secretRef.current &&
      accessAgeRef.current &&
      refreshAgeRef.current &&
      redirectRef.current
    ) {
      if (
        idRef.current.value &&
        secretRef.current.value &&
        accessAgeRef.current.value &&
        refreshAgeRef.current.value &&
        redirectRef.current.value
      ) {
        const client: ClientCreateReq = {
          id: idRef.current.value,
          secret: secretRef.current.value,
          accessTokenAge: +accessAgeRef.current.value,
          refreshTokenAge: +refreshAgeRef.current.value,
          redirects: redirectRef.current.value,
          adminPassword: adminPassword,
        };
        clientCreate(client, setWarn).then((resp) => {
          if (resp) {
            setShowAddForm(false);
          }
        });
      }
    }
  }

  return (
    <div className={"flex-basic-column"}>
      <Input label={"Id"} inputType={"text"} myRef={idRef} />
      <Input label={"Secret"} inputType={"password"} myRef={secretRef} />
      <Input
        label={"Access Token Age (hours)"}
        inputType={"text"}
        myRef={accessAgeRef}
      />
      <Input
        label={"Refresh Token Age (hours)"}
        inputType={"text"}
        myRef={refreshAgeRef}
      />
      <Input
        label={"Redirects (separate by comma)"}
        inputType={"text"}
        myRef={redirectRef}
      />
      <Button text={"Add"} click={click} />
    </div>
  );
}
