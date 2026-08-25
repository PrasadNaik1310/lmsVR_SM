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
      <div className="relative w-full max-w-2xl rounded-2xl bg-white shadow-2xl p-6 mx-4 animate-fadeIn" style={{ maxHeight: "90vh", overflowY: "auto" }}>
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
const baseURL= import.meta.env.VITE_API_BASE_URL;
    try {
      if (action === "approve") {
        const res = await fetch(`${baseURL}/admissions/applications/${application.id}/approve`, {
          method: "PATCH",
          headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
          body: JSON.stringify({ application_id: application.id }),
        });
        if (!res.ok) throw new Error(`Error returned ${res.status}`);
      } else {
        const res = await fetch(`${baseURL}/admissions/applications/${application.id}/reject`, {
          method: "PATCH",
          headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
          body: JSON.stringify({}),
        });
        if (!res.ok) throw new Error(`Error returned ${res.status}`);
      }
      onActionComplete();
      console.log("Action complete");
    } catch (err) {
      console.error(`Failed to ${action} application:`, err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="relative" ref={menuRef}>
      <button
        onClick={() => setOpen((v) => !v)}
        disabled={loading}
        title="Actions"
        className="h-7 w-7 rounded-md flex items-center justify-center text-slate-400 hover:text-slate-700 hover:bg-slate-100 transition"
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
            className="flex w-full items-center gap-2 px-3 py-2 text-sm text-emerald-700 hover:bg-emerald-50 transition"
          >
            <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <polyline points="2 7 5.5 10.5 12 4" />
            </svg>
            Approve application
          </button>
          <button
            onClick={() => handleAction("reject")}
            className="flex w-full items-center gap-2 px-3 py-2 text-sm text-red-600 hover:bg-red-50 transition"
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

// ── Import CSV Modal ──────────────────────────────────────────────────────────
function ImportCSVModal({ open, onClose, onSuccess }) {
  const token = localStorage.getItem("auth_token");
  const [rows, setRows] = useState([]);
  const [dragging, setDragging] = useState(false);
  const [importing, setImporting] = useState(false);
  const [progress, setProgress] = useState(0);
  const [done, setDone] = useState(false);
  const fileRef = useRef(null);

  const REQUIRED_COLS = ["full_name", "email", "phone", "interested_course_id", "notes"];

  const reset = () => {
    setRows([]);
    setImporting(false);
    setProgress(0);
    setDone(false);
  };

  const validateRow = (row) => {
    const errors = [];
    if (!row.full_name?.trim()) errors.push("Name required");
    if (!row.email?.trim()) errors.push("Email required");
    else if (!/\S+@\S+\.\S+/.test(row.email)) errors.push("Invalid email");
    if (!row.interested_course_id?.trim()) errors.push("Course ID required");
    return errors;
  };

  const splitCSVLine = (line) => {
    const result = []; let cur = ""; let inQ = false;
    for (const c of line) {
      if (c === '"') { inQ = !inQ; }
      else if (c === "," && !inQ) { result.push(cur); cur = ""; }
      else { cur += c; }
    }
    result.push(cur);
    return result;
  };

  const parseCSV = (text) => {
    const lines = text.trim().split(/\r?\n/);
    if (lines.length < 2) { alert("CSV has no data rows"); return; }
    const headers = lines[0].split(",").map((h) => h.trim().toLowerCase().replace(/"/g, ""));
    const colIdx = {};
    REQUIRED_COLS.forEach((col) => { colIdx[col] = headers.indexOf(col); });
    const parsed = lines.slice(1).map((line, i) => {
      const cols = splitCSVLine(line);
      const row = { _line: i + 2 };
      REQUIRED_COLS.forEach((col) => { row[col] = colIdx[col] >= 0 ? (cols[colIdx[col]] || "").trim() : ""; });
      row._errors = validateRow(row);
      return row;
    }).filter((r) => r.full_name || r.email);
    setRows(parsed);
  };

  const handleFile = (file) => {
    if (!file || !file.name.endsWith(".csv")) { alert("Please select a .csv file"); return; }
    const reader = new FileReader();
    reader.onload = (e) => parseCSV(e.target.result);
    reader.readAsText(file);
  };

  const handleImport = async () => {
    const valid = rows.filter((r) => r._errors.length === 0);
    setImporting(true);
    setProgress(0);
    let completed = 0;
    for (const row of valid) {
      try {
        await admissionsApi.createEnquiry(
          {
            full_name: row.full_name,
            email: row.email,
            phone: row.phone,
            interested_course_id: row.interested_course_id,
            notes: row.notes,
          },
          token
        );
      } catch (err) {
        console.error("Failed to import row:", row._line, err);
      }
      completed++;
      setProgress(Math.round((completed / valid.length) * 100));
    }
    setDone(true);
    setImporting(false);
    onSuccess();
  };

  const downloadTemplate = () => {
    const csv =
      "full_name,email,phone,interested_course_id,notes\n" +
      "Arjun Mehta,arjun@example.com,+91-9876543210,course-uuid-here,Interested in weekend batch\n";
    const a = document.createElement("a");
    a.href = "data:text/csv," + encodeURIComponent(csv);
    a.download = "enquiries_template.csv";
    a.click();
  };

  const valid = rows.filter((r) => r._errors.length === 0);
  const invalid = rows.filter((r) => r._errors.length > 0);

  return (
    <Modal open={open} onClose={() => { reset(); onClose(); }} title="Batch Import Enquiries">
      {/* Drop zone */}
      <div
        onDragOver={(e) => { e.preventDefault(); setDragging(true); }}
        onDragLeave={() => setDragging(false)}
        onDrop={(e) => { e.preventDefault(); setDragging(false); handleFile(e.dataTransfer.files[0]); }}
        onClick={() => fileRef.current?.click()}
        className={`border-2 border-dashed rounded-xl p-8 text-center cursor-pointer transition-colors ${
          dragging
            ? "border-indigo-400 bg-indigo-50"
            : "border-slate-200 bg-slate-50 hover:border-indigo-300 hover:bg-indigo-50"
        }`}
      >
        <p className="text-sm font-medium text-slate-700">
          Drop your CSV file here or <span className="text-indigo-600">browse</span>
        </p>
        <p className="text-xs text-slate-400 mt-1">Supports .csv files up to 5 MB</p>
        <input
          ref={fileRef}
          type="file"
          accept=".csv"
          className="hidden"
          onChange={(e) => handleFile(e.target.files[0])}
        />
      </div>

      {/* Template hint */}
      <div className="mt-3 rounded-lg bg-blue-50 border border-blue-100 px-3 py-2 text-xs text-blue-700">
        Expected columns:{" "}
        <code className="bg-blue-100 px-1 rounded font-mono">full_name</code>,{" "}
        <code className="bg-blue-100 px-1 rounded font-mono">email</code>,{" "}
        <code className="bg-blue-100 px-1 rounded font-mono">phone</code>,{" "}
        <code className="bg-blue-100 px-1 rounded font-mono">interested_course_id</code>,{" "}
        <code className="bg-blue-100 px-1 rounded font-mono">notes</code>
        {" · "}
        <button
          onClick={(e) => { e.stopPropagation(); downloadTemplate(); }}
          className="underline text-blue-600 cursor-pointer bg-transparent border-none p-0"
        >
          Download template
        </button>
      </div>

      {/* Preview table */}
      {rows.length > 0 && (
        <div className="mt-4">
          <div className="flex items-center justify-between mb-2">
            <span className="text-xs font-semibold text-slate-500 uppercase tracking-wide">Preview</span>
            <span className="text-xs text-slate-400">{rows.length} rows parsed</span>
          </div>
          <div className="rounded-lg border border-slate-200" style={{ maxHeight: "320px", overflowY: "auto" }}>
            <table className="min-w-full text-xs text-left">
              <thead className="bg-slate-50 sticky top-0 z-10">
                <tr>
                  <th className="px-3 py-2 text-slate-500 font-semibold">#</th>
                  <th className="px-3 py-2 text-slate-500 font-semibold">Name</th>
                  <th className="px-3 py-2 text-slate-500 font-semibold">Email</th>
                  <th className="px-3 py-2 text-slate-500 font-semibold">Status</th>
                </tr>
              </thead>
              <tbody>
                {rows.map((row, i) => (
                  <tr key={i} className="border-t border-slate-100 hover:bg-slate-50">
                    <td className="px-3 py-2 text-slate-400">{row._line}</td>
                    <td className="px-3 py-2 text-slate-700">{row.full_name || <span className="text-red-400">—</span>}</td>
                    <td className="px-3 py-2 text-slate-500">{row.email || <span className="text-red-400">—</span>}</td>
                    <td className="px-3 py-2">
                      {row._errors.length === 0 ? (
                        <span className="rounded-full bg-green-50 text-green-700 px-2 py-0.5 text-xs font-medium">✓ Valid</span>
                      ) : (
                        <span
                          className="rounded-full bg-red-50 text-red-600 px-2 py-0.5 text-xs font-medium"
                          title={row._errors.join(", ")}
                        >
                          ✕ {row._errors[0]}
                        </span>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Summary cards */}
          <div className="grid grid-cols-3 gap-2 mt-3">
            <div className="rounded-lg bg-slate-50 border border-slate-200 p-2 text-center">
              <p className="text-lg font-semibold text-slate-800">{rows.length}</p>
              <p className="text-xs text-slate-500">Total rows</p>
            </div>
            <div className="rounded-lg bg-green-50 border border-green-100 p-2 text-center">
              <p className="text-lg font-semibold text-green-700">{valid.length}</p>
              <p className="text-xs text-green-600">Ready to import</p>
            </div>
            <div className="rounded-lg bg-red-50 border border-red-100 p-2 text-center">
              <p className="text-lg font-semibold text-red-600">{invalid.length}</p>
              <p className="text-xs text-red-500">Rows with errors</p>
            </div>
          </div>
        </div>
      )}

      {/* Progress bar */}
      {importing && (
        <div className="mt-4">
          <div className="h-1.5 w-full bg-slate-100 rounded-full overflow-hidden">
            <div
              className="h-full bg-indigo-500 rounded-full transition-all duration-300"
              style={{ width: `${progress}%` }}
            />
          </div>
          <p className="text-xs text-slate-500 mt-1 text-right">{progress}% complete</p>
        </div>
      )}

      {done && (
        <p className="mt-3 text-xs text-green-700 bg-green-50 rounded-lg px-3 py-2">
          ✓ {valid.length} enquiries imported successfully
        </p>
      )}

      {/* Footer buttons */}
      <div className="flex justify-end gap-2 mt-4">
        <button
          onClick={() => { reset(); onClose(); }}
          className="rounded-lg px-4 py-2 text-sm text-slate-600 hover:bg-slate-100 transition"
        >
          Cancel
        </button>
        <button
          onClick={handleImport}
          disabled={valid.length === 0 || importing || done}
          className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-50 transition"
        >
          {importing ? `Importing… ${progress}%` : `Import ${valid.length} valid row${valid.length !== 1 ? "s" : ""}`}
        </button>
      </div>
    </Modal>
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
  const [showImportModal, setShowImportModal] = useState(false);

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
      <ImportCSVModal
        open={showImportModal}
        onClose={() => setShowImportModal(false)}
        onSuccess={loadData}
      />

      <div className="max-w-[1200px] mx-auto">
        {/* Header */}
        <div className="mb-6 flex items-center justify-between">
          <h1 className="text-2xl font-semibold text-slate-900">Admissions Overview</h1>
          <div className="flex items-center gap-3">
            <button
              onClick={() => setShowImportModal(true)}
              className="rounded-lg border border-indigo-200 bg-white px-4 py-2 text-sm font-medium text-indigo-600 hover:bg-indigo-50 transition"
            >
              ↑ Import CSV
            </button>
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
                  <button
                    onClick={() => setShowImportModal(true)}
                    className="rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 hover:bg-indigo-50 hover:border-indigo-200 hover:text-indigo-700 transition text-left"
                  >
                    ↑ Import CSV
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