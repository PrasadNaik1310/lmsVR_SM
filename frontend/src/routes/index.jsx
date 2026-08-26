import { createBrowserRouter } from "react-router-dom";
{/*import AppLayout from "../layouts/AppLayout.jsx";*/ }
import AppLayout from "../AppLayout.jsx";
import Home from "../pages/Home.jsx";
import ManageCompanyBatches from "../pages/ManageCompanyBatches.jsx";
import AdmissionsOverview from "../pages/AdmissionsOverview.jsx";
import CourseManagement from "../pages/CourseManagement.jsx";
import CourseDetails from "../pages/CourseDetails";
import ModuleDetails from "../pages/ModuleDetails";
import AcademicSessionManagement from "../pages/AcademicSessionManagement";
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
        path: "company/academic-sessions",
        element: <AcademicSessionManagement />,
      },
      {
        path: "admissions",
        element: <AdmissionsOverview />,
      },
      {
        path: "courses",
        element: <CourseManagement />,
      },
      {
        path: "/courses/:id",
        element: <CourseDetails />,
      },
      {
        path: "/modules/:id",
        element: <ModuleDetails />,
      }
    ],
  },
]);