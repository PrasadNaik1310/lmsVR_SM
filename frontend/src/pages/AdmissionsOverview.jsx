import React, { useState, useEffect } from "react";
import Sidebar from "../components/Sidebar.jsx"
import * as admissionsApi from "../api/admissions.js";

export default function AdmissionsOverview() {
  const [enquiries, setEnquiries] = useState([]);
  const [applications, setApplications] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const loadData = async () => {
      try {
        setLoading(true);
        const token = localStorage.getItem("auth_token");
        if (!token) {
          throw new Error("No authentication token found. Please login first.");
        }
        const [enquiriesRes, applicationsRes] = await Promise.all([
          admissionsApi.listEnquiries({ page: 1, size: 10 }, token),
          admissionsApi.listApplications({ page: 1, size: 10 }, token),
        ]);
        setEnquiries(enquiriesRes.enquiries|| []);
        console.log(enquiriesRes)
        setApplications(applicationsRes.applications || []);
        console.log(applicationsRes)
        setError(null);
      } catch (err) {
        setError(err.message || 'Failed to load data');
        console.error('Failed to load admissions data:', err);
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, []);

  // Calculate stats from fetched data
  const stats = [
    { title: 'Total Enquiries', value: enquiries.length, color: 'text-emerald-600', bg: 'bg-emerald-50' },
    { title: 'Pending Applications', value: applications.filter(a => a.application_status === 'pending').length, color: 'text-sky-600', bg: 'bg-sky-50' },
    { title: 'Approved Applications', value: applications.filter(a => a.application_status === 'approved').length, color: 'text-amber-600', bg: 'bg-amber-50' },
    { title: 'Rejected Applications', value: applications.filter(a => a.application_status === 'rejected').length, color: 'text-rose-600', bg: 'bg-rose-50' },
  ];

  const formatDate = (dateStr) => {
    if (!dateStr) return 'N/A';
    try {
      return new Date(dateStr).toLocaleDateString('en-IN', { day: 'numeric', month: 'short', year: 'numeric' });
    } catch {
      return dateStr;
    }
  };

  return (
    
    <div className="min-h-screen p-6 bg-slate-50">
      
      
      <div className="max-w-[1200px] mx-auto">
        <div className="mb-6 flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-semibold text-slate-900">Admissions Overview</h1>
        
          </div>
          <div className="flex items-center gap-3">
            <button className="rounded bg-blue-600 px-4 py-2 text-white">+ Create Enquiry</button>
            <div className="h-9 w-9 rounded-full bg-slate-200 text-center leading-9">A</div>
          </div>
        </div>

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
            Error loading admissions data: {error}
          </div>
        )}

        {loading ? (
          <div className="rounded-xl border border-slate-100 bg-white p-6 text-center">
            <p className="text-slate-600">Loading admissions data...</p>
          </div>
        ) : (
          <div className="grid gap-6 lg:grid-cols-3">
            <div className="lg:col-span-2 space-y-4">
              <div className="rounded-xl border border-slate-100 bg-white p-4">
                <div className="flex items-center justify-between">
                <h2 className="text-sm font-semibold text-slate-900">Recent Enquiries</h2>
                <a className="text-sm text-indigo-600">View all</a>
              </div>
              <div className="mt-4 overflow-x-auto">
                <table className="min-w-full text-left text-sm">
                  <thead className="bg-slate-50 text-xs uppercase text-slate-500">
                    <tr>
                      <th className="px-4 py-3">Name</th>
                      <th className="px-4 py-3">Email</th>
                      <th className="px-4 py-3">Phone</th>
                      <th className="px-4 py-3">Course</th>
                      <th className="px-4 py-3">Status</th>
                      <th className="px-4 py-3">Date</th>
                      <th className="px-4 py-3">Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {enquiries.map((e) => (
                      <tr key={e.id} className="border-t border-slate-100">
                        <td className="px-4 py-3">{e.full_name}</td>
                        <td className="px-4 py-3 text-slate-600">{e.email}</td>
                        <td className="px-4 py-3 text-slate-600">{e.phone}</td>
                        <td className="px-4 py-3 text-slate-600">{e.course_id}</td>
                        <td className="px-4 py-3"><span className="rounded-full bg-amber-50 px-2 py-1 text-xs text-amber-700">{e.status}</span></td>
                        <td className="px-4 py-3 text-slate-600">{formatDate(e.created_at)}</td>
                        <td className="px-4 py-3">...</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
              <p className="mt-3 text-xs text-slate-500">Total {enquiries.length} enquiries</p>
            </div>

            <div className="rounded-xl border border-slate-100 bg-white p-4">
              <h2 className="text-sm font-semibold text-slate-900">Recent Applications</h2>
              <div className="mt-4 overflow-x-auto">
                {applications.length > 0 ? (
                  <table className="min-w-full text-left text-sm">
                    <thead className="bg-slate-50 text-xs uppercase text-slate-500">
                      <tr>
                        <th className="px-4 py-3">Enquiry</th>
                        <th className="px-4 py-3">Course</th>
                        <th className="px-4 py-3">Status</th>
                        <th className="px-4 py-3">Submitted</th>
                      </tr>
                    </thead>
                    <tbody>
                      {applications.map((a) => (
                        <tr key={a.id} className="border-t border-slate-100">
                          <td className="px-4 py-3">{a.enquiry_id}</td>
                          <td className="px-4 py-3 text-slate-600">{a.applied_course_id}</td>
                          <td className="px-4 py-3"><span className={`rounded-full px-2 py-1 text-xs ${a.application_status === 'approved' ? 'bg-green-50 text-green-700' : a.application_status === 'rejected' ? 'bg-red-50 text-red-700' : 'bg-blue-50 text-blue-700'}`}>{a.application_status}</span></td>
                          <td className="px-4 py-3 text-slate-600">{formatDate(a.submitted_at)}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                ) : (
                  <p className="text-sm text-slate-600">No applications yet</p>
                )}
              </div>
            </div>
          </div>

          <aside className="space-y-4">
   {         /*<div className="rounded-xl border border-slate-100 bg-white p-4">
              <h3 className="text-sm font-semibold text-slate-900">Application Status</h3>
              <div className="mt-4 flex items-center justify-center">
                <div className="h-36 w-36 rounded-full bg-slate-50 flex items-center justify-center">
                  <div className="text-center">
                    <p className="text-2xl font-semibold">106</p>
                    <p className="text-xs text-slate-500">Total</p>
                  </div>
                </div>
              </div>
            </div>
*/}
            <div className="rounded-xl border border-slate-100 bg-white p-4">
              <h3 className="text-sm font-semibold text-slate-900">Quick Actions</h3>
              <div className="mt-3 grid gap-2">
                <button className="rounded border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700">Create Enquiry</button>
                <button className="rounded border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700">Create Application</button>
              </div>
            </div>
          </aside>
          </div>
        )}
      </div>
      
    </div>
  );
}
