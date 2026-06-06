import {useLocation} from "react-router-dom"
import {Outlet} from "react-router-dom";
import Sidebar from "./components/Sidebar.jsx";
export default function AppLayout(){
    const location = useLocation();
    const noSidebarRoutes = ["/login","/"]
    const showSidebar = !noSidebarRoutes.includes(location.pathname);
 
    return (
    <div className = "app flex">
        {showSidebar && <Sidebar/>}
        <div className = {showSidebar ? "flex-1 ": "w-full"}>
            <Outlet/>
        </div>
    </div>
);
}