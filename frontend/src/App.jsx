  import { RouterProvider } from "react-router-dom";
  import {router} from "./routes/index.jsx";
  {/*import {  useLocation } from "react-router-dom";
  import AppLayout from "./AppLayout"
  import Sidebar from "./components/Sidebar";*/}

  export default function App() {
    return <RouterProvider router = {router} />;
  }