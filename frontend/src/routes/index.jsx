import { createBrowserRouter } from "react-router-dom";
import AppLayout from "../layouts/AppLayout.jsx";
import Home from "../pages/Home.jsx";
import ManageCompanyBatches from "../pages/ManageCompanyBatches.jsx";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <AppLayout />,
    children: [
      {
        index: true,
        element: <Home />,
      },
      {
        path: "company/batches",
        element: <ManageCompanyBatches />,
      },
    ],
  },
]);