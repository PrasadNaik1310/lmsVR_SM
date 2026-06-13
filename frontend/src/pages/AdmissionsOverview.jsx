import React, { useState, useEffect, useRef } from "react";
import * as admissionsApi from "../api/admissions.js";
import axios from "axios";
// ── Modal backdrop ────────────────────────────────────────────────────────────
function Modal({ open, onClose, title, children }) {
  useEffect(() => {
    if (open) document.body.style.overflow = "hidden";
    else document.body.style.overflow = "";
    return () => { document.body.style.overflow = ""; };
  }, [open]);

  if (!open) return null;
  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center"
      style={{ background: "rgba(15,23,42,0.45)", backdropFilter: "blur(2px)" }}
      onClick={(e) => { if (e.target === e.currentTarget) onClose(); }}
    >
      <div className="relative w-full max-w-md rounded-2xl bg-white shadow-2xl p-6 mx-4 animate-fadeIn">
        <div className="flex items-center justify-between mb-5">
          <h2 className="text-base font-semibold text-slate-900">{title}</h2>
          <button
            onClick={onClose}
            className="h-7 w-7 rounded-full bg-slate-100 hover:bg-slate-200 flex items-center justify-center text-slate-500 transition-colors"
          >
            ✕
          </button>
        </div>
        {children}
      </div>
    </div>
  );
}

// ── Field ─────────────────────────────────────────────────────────────────────
function Field({ label, error, children }) {
  return (
    <div className="mb-4">
      <label className="block text-xs font-medium text-slate-600 mb-1">{label}</label>
      {children}
      {error && <p className="mt-1 text-xs text-red-500">{error}</p>}
    </div>
  );
}

const inputCls =
  "w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-800 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-400 focus:border-transparent transition";

// ── Row Action Menu ───────────────────────────────────────────────────────────
function ApplicationActionMenu({ application, onActionComplete }) {
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const menuRef = useRef(null);

  // Close on outside click
  useEffect(() => {
    if (!open) return;
    const handler = (e) => {
      if (menuRef.current && !menuRef.current.contains(e.target)) setOpen(false);
    };
    document.addEventListener("mousedown", handler);
    return () => document.removeEventListener("mousedown", handler);
  }, [open]);

  const handleAction = async (action) => {
    setOpen(false);
    setLoading(true);
    const token = localStorage.getItem("auth_token");
    
    try {
      
      if (action === "approve") {
        const res = await fetch(`http://localhost:8080/lms/admissions/applications/${application.id}/approve`, {
          method: "PATCH",
          headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify({application_id:application.id})
          //body: application.id

        });
console.log("Request data:->  %s ",body);
        /*axios.interceptors.request.use((res) => {
  console.log("REQUEST");
  console.log(res.method);
  console.log(res.url);
  console.log(res.data);

  return config;
});*/
        if (!res.ok) throw new Error('Error returned ${res.status}');
      } else {
        const res = await fetch(`http://localhost:8080/lms/admissions/applications/${application.id}/reject`, {
          method: "PATCH",
          headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
          body: JSON.stringify({})
        });
        if (!res.ok) throw new Error("Error returned ${res.status}");
      }
      onActionComplete();
      console.log("Action complete");
    } catch (err) {
      console.error(`Failed to ${action} application:`, err);
    } finally {
      setLoading(false);
    }
  };

  const isPending = application.application_status === "pending";

  return (
    <div className="relative" ref={menuRef}>
      <button
        onClick={() => setOpen((v) => !v)}
        disabled={loading}
        title="Actions"
        className="h-7 w-7 rounded-md flex items-center justify-center text-slate-400 hover:text-slate-700 hover:bg-slate-100 transition "
      >
        {loading ? (
          <span className="block h-3 w-3 rounded-full border-2 border-slate-300 border-t-indigo-500 animate-spin" />
        ) : (
          <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
            <circle cx="8" cy="3" r="1.25" />
            <circle cx="8" cy="8" r="1.25" />
            <circle cx="8" cy="13" r="1.25" />
          </svg>
        )}
      </button>

      {open && (
        <div className="absolute right-0 z-40 mt-1 w-44 rounded-xl border border-slate-100 bg-white shadow-lg py-1 animate-fadeIn">
          <button
            onClick={() => handleAction("approve")}
            /*disabled={!isPending}*/
            className="flex w-full items-center gap-2 px-3 py-2 text-sm text-emerald-700 hover:bg-emerald-50 transition "
          >
            <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <polyline points="2 7 5.5 10.5 12 4" />
            </svg>
            Approve application
          </button>
          <button
            onClick={() => handleAction("reject")}
            /*disabled={!isPending}*/
            className="flex w-full items-center gap-2 px-3 py-2 text-sm text-red-600 hover:bg-red-50 transition "
          >
            <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <line x1="3" y1="3" x2="11" y2="11" />
              <line x1="11" y1="3" x2="3" y2="11" />
            </svg>
            Reject application
          </button>
        </div>
      )}
    </div>
  );
}

// ── Create Enquiry Modal ──────────────────────────────────────────────────────
function CreateEnquiryModal({ open, onClose, onSuccess }) {
  const token = localStorage.getItem("auth_token");
  const [form, setForm] = useState({ full_name: "", email: "", phone: "", interested_course_id: "", notes: "" });
  const [errors, setErrors] = useState({});
  const [submitting, setSubmitting] = useState(false);
  const [apiError, setApiError] = useState("");

  const set = (k, v) => setForm((f) => ({ ...f, [k]: v }));

  const validate = () => {
    const e = {};
    if (!form.full_name.trim()) e.full_name = "Name is required";
    if (!form.email.trim()) e.email = "Email is required";
    else if (!/\S+@\S+\.\S+/.test(form.email)) e.email = "Invalid email";
    if (!form.interested_course_id.trim()) e.interested_course_id = "Course ID is required";
    return e;
  };

  const handleSubmit = async () => {
    const e = validate();
    if (Object.keys(e).length) { setErrors(e); return; }
    setSubmitting(true);
    setApiError("");
    try {
      await admissionsApi.createEnquiry(form, token);
      setForm({ full_name: "", email: "", phone: "", interested_course_id: "", notes: "" });
      setErrors({});
      onSuccess();
      onClose();
    } catch (err) {
      setApiError(err?.response?.data?.error || err.message || "Failed to create enquiry");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Modal open={open} onClose={onClose} title="Create Enquiry">
      <Field label="Full Name *" error={errors.full_name}>
        <input className={inputCls} placeholder="e.g. Arjun Mehta" value={form.full_name} onChange={(e) => set("full_name", e.target.value)} />
      </Field>
      <Field label="Email *" error={errors.email}>
        <input className={inputCls} placeholder="email@example.com" value={form.email} onChange={(e) => set("email", e.target.value)} />
      </Field>
      <Field label="Phone" error={errors.phone}>
        <input className={inputCls} placeholder="+91-9876543210" value={form.phone} onChange={(e) => set("phone", e.target.value)} />
      </Field>
      <Field label="Course ID *" error={errors.interested_course_id}>
        <input className={inputCls} placeholder="UUID of interested course" value={form.interested_course_id} onChange={(e) => set("interested_course_id", e.target.value)} />
      </Field>
      <Field label="Notes">
        <textarea className={inputCls + " resize-none"} rows={3} placeholder="Any additional notes..." value={form.notes} onChange={(e) => set("notes", e.target.value)} />
      </Field>
      {apiError && <p className="mb-3 text-xs text-red-500 rounded bg-red-50 px-3 py-2">{apiError}</p>}
      <div className="flex justify-end gap-2 mt-2">
        <button onClick={onClose} className="rounded-lg px-4 py-2 text-sm text-slate-600 hover:bg-slate-100 transition">Cancel</button>
        <button
          onClick={handleSubmit}
          disabled={submitting}
          className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-60 transition"
        >
          {submitting ? "Creating…" : "Create Enquiry"}
        </button>
      </div>
    </Modal>
  );
}

// ── Create Application Modal ──────────────────────────────────────────────────
function CreateApplicationModal({ open, onClose, onSuccess, enquiries }) {
  const token = localStorage.getItem("auth_token");
  const [form, setForm] = useState({ enquiry_id: "", applied_course_id: "" });
  const [errors, setErrors] = useState({});
  const [submitting, setSubmitting] = useState(false);
  const [apiError, setApiError] = useState("");

  const set = (k, v) => setForm((f) => ({ ...f, [k]: v }));

  const validate = () => {
    const e = {};
    if (!form.enquiry_id.trim()) e.enquiry_id = "Enquiry is required";
    if (!form.applied_course_id.trim()) e.applied_course_id = "Course ID is required";
    return e;
  };

  const handleEnquiryChange = (id) => {
    set("enquiry_id", id);
    const enq = enquiries.find((e) => e.id === id);
    if (enq) set("applied_course_id", enq.interested_course_id);
  };

  const handleSubmit = async () => {
    const e = validate();
    if (Object.keys(e).length) { setErrors(e); return; }
    setSubmitting(true);
    setApiError("");
    try {
      await admissionsApi.createApplication(form, token);
      setForm({ enquiry_id: "", applied_course_id: "" });
      setErrors({});
      onSuccess();
      onClose();
    } catch (err) {
      setApiError(err?.response?.data?.error || err.message || "Failed to create application");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Modal open={open} onClose={onClose} title="Create Application">
      <Field label="Enquiry *" error={errors.enquiry_id}>
        <select
          className={inputCls}
          value={form.enquiry_id}
          onChange={(e) => handleEnquiryChange(e.target.value)}
        >
          <option value="">Select an enquiry…</option>
          {enquiries.map((enq) => (
            <option key={enq.id} value={enq.id}>{enq.full_name} — {enq.email}</option>
          ))}
        </select>
      </Field>
      <Field label="Course ID *" error={errors.applied_course_id}>
        <input
          className={inputCls}
          placeholder="UUID of applied course"
          value={form.applied_course_id}
          onChange={(e) => set("applied_course_id", e.target.value)}
        />
      </Field>
      {apiError && <p className="mb-3 text-xs text-red-500 rounded bg-red-50 px-3 py-2">{apiError}</p>}
      <div className="flex justify-end gap-2 mt-2">
        <button onClick={onClose} className="rounded-lg px-4 py-2 text-sm text-slate-600 hover:bg-slate-100 transition">Cancel</button>
        <button
          onClick={handleSubmit}
          disabled={submitting}
          className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-60 transition"
        >
          {submitting ? "Creating…" : "Create Application"}
        </button>
      </div>
    </Modal>
  );
}

// ── Main Page ─────────────────────────────────────────────────────────────────
export default function AdmissionsOverview() {
  const [enquiries, setEnquiries] = useState([]);
  const [applications, setApplications] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showEnquiryModal, setShowEnquiryModal] = useState(false);
  const [showApplicationModal, setShowApplicationModal] = useState(false);

  const loadData = async () => {
    try {
      setLoading(true);
      const token = localStorage.getItem("auth_token");
      if (!token) throw new Error("No authentication token found. Please login first.");
      const [enquiriesRes, applicationsRes] = await Promise.all([
        admissionsApi.listEnquiries({ page: 1, size: 10 }, token),
        admissionsApi.listApplications({ page: 1, size: 10 }, token),
      ]);
      setEnquiries(enquiriesRes.enquiries || []);
      setApplications(applicationsRes.applications || []);
      setError(null);
    } catch (err) {
      setError(err.message || "Failed to load data");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { loadData(); }, []);

  const stats = [
    { title: "Total Enquiries", value: enquiries.length, color: "text-emerald-600", bg: "bg-emerald-50" },
    { title: "Pending Applications", value: applications.filter((a) => a.application_status === "pending").length, color: "text-sky-600", bg: "bg-sky-50" },
    { title: "Approved Applications", value: applications.filter((a) => a.application_status === "approved").length, color: "text-amber-600", bg: "bg-amber-50" },
    { title: "Rejected Applications", value: applications.filter((a) => a.application_status === "rejected").length, color: "text-rose-600", bg: "bg-rose-50" },
  ];

  const formatDate = (dateStr) => {
    if (!dateStr) return "N/A";
    try { return new Date(dateStr).toLocaleDateString("en-IN", { day: "numeric", month: "short", year: "numeric" }); }
    catch { return dateStr; }
  };

  const statusColor = (status) => {
    const map = {
      new: "bg-amber-50 text-amber-700",
      contacted: "bg-blue-50 text-blue-700",
      follow_up: "bg-purple-50 text-purple-700",
      converted: "bg-green-50 text-green-700",
      closed: "bg-slate-100 text-slate-500",
      pending: "bg-amber-50 text-amber-700",
      approved: "bg-green-50 text-green-700",
      rejected: "bg-red-50 text-red-700",
    };
    return map[status] || "bg-slate-100 text-slate-600";
  };

  return (
    <div className="min-h-screen p-6 bg-slate-50">
      {/* Modals */}
      <CreateEnquiryModal
        open={showEnquiryModal}
        onClose={() => setShowEnquiryModal(false)}
        onSuccess={loadData}
      />
      <CreateApplicationModal
        open={showApplicationModal}
        onClose={() => setShowApplicationModal(false)}
        onSuccess={loadData}
        enquiries={enquiries}
      />

      <div className="max-w-[1200px] mx-auto">
        {/* Header */}
        <div className="mb-6 flex items-center justify-between">
          <h1 className="text-2xl font-semibold text-slate-900">Admissions Overview</h1>
          <div className="flex items-center gap-3">
            <button
              onClick={() => setShowEnquiryModal(true)}
              className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 transition"
            >
              + Create Enquiry
            </button>
            <div className="h-9 w-9 rounded-full bg-slate-200 text-center leading-9 text-sm font-medium text-slate-600">A</div>
          </div>
        </div>

        {/* Stats */}
        <div className="grid gap-4 md:grid-cols-4 mb-6">
          {stats.map((s) => (
            <div key={s.title} className="rounded-xl border border-slate-100 bg-white p-4">
              <div className="flex items-center justify-between">
                <p className="text-xs text-slate-500">{s.title}</p>
                <div className={`h-8 w-8 rounded ${s.bg}`} />
              </div>
              <p className="mt-3 text-2xl font-semibold text-slate-900">{s.value}</p>
            </div>
          ))}
        </div>

        {error && (
          <div className="mb-4 rounded-xl border border-red-100 bg-red-50 p-4 text-sm text-red-700">
            {error}
          </div>
        )}

        {loading ? (
          <div className="rounded-xl border border-slate-100 bg-white p-6 text-center">
            <p className="text-slate-600">Loading admissions data…</p>
          </div>
        ) : (
          <div className="grid gap-6 lg:grid-cols-3">
            <div className="lg:col-span-2 space-y-4">
              {/* Enquiries table */}
              <div className="rounded-xl border border-slate-100 bg-white p-4">
                <div className="flex items-center justify-between mb-4">
                  <h2 className="text-sm font-semibold text-slate-900">Recent Enquiries</h2>
                  <span className="text-xs text-slate-400">Total {enquiries.length}</span>
                </div>
                <div className="overflow-x-auto">
                  <table className="min-w-full text-left text-sm">
                    <thead className="bg-slate-50 text-xs uppercase text-slate-500">
                      <tr>
                        <th className="px-4 py-3">Name</th>
                        <th className="px-4 py-3">Email</th>
                        <th className="px-4 py-3">Phone</th>
                        <th className="px-4 py-3">Status</th>
                        <th className="px-4 py-3">Date</th>
                      </tr>
                    </thead>
                    <tbody>
                      {enquiries.map((e) => (
                        <tr key={e.id} className="border-t border-slate-100 hover:bg-slate-50 transition">
                          <td className="px-4 py-3 font-medium text-slate-800">{e.full_name}</td>
                          <td className="px-4 py-3 text-slate-500">{e.email}</td>
                          <td className="px-4 py-3 text-slate-500">{e.phone}</td>
                          <td className="px-4 py-3">
                            <span className={`rounded-full px-2 py-1 text-xs font-medium ${statusColor(e.status)}`}>{e.status}</span>
                          </td>
                          <td className="px-4 py-3 text-slate-500">{formatDate(e.created_at)}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Applications table */}
              <div className="rounded-xl border border-slate-100 bg-white p-4">
                <div className="flex items-center justify-between mb-4">
                  <h2 className="text-sm font-semibold text-slate-900">Recent Applications</h2>
                  <span className="text-xs text-slate-400">Total {applications.length}</span>
                </div>
                <div className="overflow-x-auto">
                  {applications.length > 0 ? (
                    <table className="min-w-full text-left text-sm">
                      <thead className="bg-slate-50 text-xs uppercase text-slate-500">
                        <tr>
                          <th className="px-4 py-3">Application ID</th>
                          <th className="px-4 py-3">Course ID</th>
                          <th className="px-4 py-3">Status</th>
                          <th className="px-4 py-3">Submitted</th>
                          <th className="px-4 py-3 w-10"></th>
                        </tr>
                      </thead>
                      <tbody>
                        {applications.map((a) => (
                          <tr key={a.id} className="border-t border-slate-100 hover:bg-slate-50 transition">
                            <td className="px-4 py-3 text-slate-500 font-mono text-xs">{a.enquiry_id}</td>
                            <td className="px-4 py-3 text-slate-500 font-mono text-xs">{a.applied_course_id}</td>
                            <td className="px-4 py-3">
                              <span className={`rounded-full px-2 py-1 text-xs font-medium ${statusColor(a.application_status)}`}>
                                {a.application_status}
                              </span>
                            </td>
                            <td className="px-4 py-3 text-slate-500">{formatDate(a.submitted_at)}</td>
                            <td className="px-4 py-3">
                              <ApplicationActionMenu
                                application={a}
                                onActionComplete={loadData}
                              />
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  ) : (
                    <p className="text-sm text-slate-500 py-4 text-center">No applications yet</p>
                  )}
                </div>
              </div>
            </div>

            {/* Sidebar */}
            <aside className="space-y-4">
              <div className="rounded-xl border border-slate-100 bg-white p-4">
                <h3 className="text-sm font-semibold text-slate-900 mb-3">Quick Actions</h3>
                <div className="grid gap-2">
                  <button
                    onClick={() => setShowEnquiryModal(true)}
                    className="rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 hover:bg-indigo-50 hover:border-indigo-200 hover:text-indigo-700 transition text-left"
                  >
                    + Create Enquiry
                  </button>
                  <button
                    onClick={() => setShowApplicationModal(true)}
                    className="rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 hover:bg-indigo-50 hover:border-indigo-200 hover:text-indigo-700 transition text-left"
                  >
                    + Create Application
                  </button>
                </div>
              </div>
            </aside>
          </div>
        )}
      </div>

      <style>{`
        @keyframes fadeIn {
          from { opacity: 0; transform: translateY(-8px) scale(0.98); }
          to   { opacity: 1; transform: translateY(0) scale(1); }
        }
        .animate-fadeIn { animation: fadeIn 0.18s ease; }
        @keyframes spin { to { transform: rotate(360deg); } }
        .animate-spin { animation: spin 0.7s linear infinite; }
      `}</style>
    </div>
  );
}