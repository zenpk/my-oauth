import { useState } from "react";

function App() {
  const [adminPassword, setAdminPassword] = useState("");
  console.log(adminPassword);
  setAdminPassword("a");
  return (
    <>
      <h1>Admin Page</h1>
    </>
  );
}

export default App;
