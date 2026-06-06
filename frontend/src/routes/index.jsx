import { createBrowserRouter } from "react-router-dom";
{/*import AppLayout from "../layouts/AppLayout.jsx";*/}
import AppLayout from "../AppLayout.jsx";
import Home from "../pages/Home.jsx";
import ManageCompanyBatches from "../pages/ManageCompanyBatches.jsx";
import AdmissionsOverview from "../pages/AdmissionsOverview.jsx";

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
      {
        path: "admissions",
        element: <AdmissionsOverview />,
      },
    ],
  },
]);