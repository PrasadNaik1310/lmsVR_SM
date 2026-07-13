import React, { createContext, useContext, useState } from "react";

const TabsContext = createContext();

export function Tabs({ defaultValue, children }) {
  const [value, setValue] = useState(defaultValue || null);

  return (
    <TabsContext.Provider value={{ value, setValue }}>
      <div>{children}</div>
    </TabsContext.Provider>
  );
}

export function TabsList({ children, className = "mb-4 flex space-x-2" }) {
  return <div className={className}>{children}</div>;
}

export function TabsTrigger({ value, children, className = "px-3 py-1 rounded", activeClassName = "bg-gray-900 text-white" }) {
  const ctx = useContext(TabsContext);
  if (!ctx) return null;

  const isActive = ctx.value === value;

  return (
    <button
      type="button"
      onClick={() => ctx.setValue(value)}
      className={`${className} ${isActive ? activeClassName : "bg-white text-gray-700 border"}`}
    >
      {children}
    </button>
  );
}

export function TabsContent({ value, children, className = "" }) {
  const ctx = useContext(TabsContext);
  if (!ctx) return null;

  return ctx.value === value ? <div className={className}>{children}</div> : null;
}

export default Tabs;
