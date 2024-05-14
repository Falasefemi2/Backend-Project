import { Route, Routes } from "react-router-dom";
import CreateUser from "./CreateUser";
import MainChat from "./MainChat";
import Login from "./Login";

export default function App() {
  return (
    <>
      <Routes>
        <Route path="/create-user" element={<CreateUser />} />
        <Route path="/chat" element={<MainChat />} />
        <Route path="/chat/:channelId" element={<MainChat />} />
        <Route path="/" element={<Login />} />
      </Routes>
    </>
  )
}