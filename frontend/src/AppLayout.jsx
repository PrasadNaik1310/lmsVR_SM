import { useLocation, useNavigate } from "react-router-dom"
import { Outlet } from "react-router-dom";
import Sidebar from "./components/Sidebar.jsx";

export default function AppLayout() {
    const location = useLocation();
    const navigate = useNavigate();
    const noSidebarRoutes = ["/login", "/"]
    const showSidebar = !noSidebarRoutes.includes(location.pathname);

    const handleLogout = () => {
        localStorage.removeItem('auth_token');
        navigate('/');
    };

    return (
        <div className="app flex">
            {showSidebar && <Sidebar onLogout={handleLogout} />}
            <div className={showSidebar ? "flex-1" : "w-full"}>
                <Outlet />
            </div>
        </div>
    );
}