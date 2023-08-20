import { Dispatch, SetStateAction, useEffect, useRef, useState } from "react";
import { Input } from "./components/Input.tsx";
import { Button } from "./components/Button.tsx";

export function Admin() {
  const [adminPassword, setAdminPassword] = useState("");
  const [warn, setWarn] = useState("");
  const [showAddForm, setShowAddForm] = useState(false);
  const adminPasswordRef = useRef<HTMLInputElement | null>(null);

  useEffect(() => {
    if (adminPassword !== "" && !showAddForm) {
    }
  }, [adminPassword, showAddForm]);

  function saveAdminPassword() {
    if (
      adminPasswordRef &&
      adminPasswordRef.current &&
      adminPasswordRef.current.value
    ) {
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
          {showAddForm && <AddForm setShowAddForm={setShowAddForm} />}
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
            <tbody></tbody>
          </table>
        </div>
      )}
    </div>
  );
}

function AddForm({
  setShowAddForm,
}: {
  setShowAddForm: Dispatch<SetStateAction<boolean>>;
}) {
  const idRef = useRef<HTMLInputElement | null>(null);
  const secretRef = useRef<HTMLInputElement | null>(null);
  const accessAgeRef = useRef<HTMLInputElement | null>(null);
  const refreshAgeRef = useRef<HTMLInputElement | null>(null);
  const redirectRef = useRef<HTMLInputElement | null>(null);

  function click() {
    setShowAddForm(false);
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
