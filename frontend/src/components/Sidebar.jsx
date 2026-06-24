import { NavLink } from "react-router-dom";

const routeLinks = [
  { label: "Dashboard", to: "/" },
  { label: "Admissions", to: "/admissions" },
  { label: "Manage Company", to: "/company/batches" },
  { label: "Courses",to: "/courses"}
];

const sectionLinks = [
  { label: "Academic Sessions", href: "academic-sessions" },
  { label: "Courses / Packages", href: "#courses-packages" },
  { label: "Batches", href: "#batches" },
  { label: "Internal Team", href: "#internal-team" },
];

export default function Sidebar({ onLogout }) {
  return (
    <aside className="hidden w-64 shrink-0 border-r border-slate-100 bg-white lg:block">
      <nav className=" px-3pb-6">
        <ul className="space-y-1 text-sm text-slate-700">
          {routeLinks.map((item) => (
            <li key={item.to}>
              <NavLink
                to={item.to}
                end={item.to === "/"}
                className={({ isActive }) =>
                  [
                    "block rounded-lg px-3 py-2 transition-colors",
                    isActive ? "bg-indigo-50 font-semibold text-indigo-700" : "hover:bg-slate-100",
                  ].join(" ")
                }
              >
                {item.label}
              </NavLink>
            </li>
          ))}

          {sectionLinks.map((item) => (
            <li key={item.href}>
              <a
                href={item.href}
                className="block rounded-lg px-3 py-2 transition-colors hover:bg-slate-100"
              >
                {item.label}
              </a>
            </li>
          ))}
        </ul>

        <div className="mt-8 border-t border-slate-200 pt-6">
          <button
            type="button"
            onClick={onLogout}
            className="w-full rounded-lg px-3 py-2 text-left text-sm font-medium text-red-600 hover:bg-red-50"
          >
            Logout
          </button>
        </div>
      </nav>
    </aside>
  );
}